package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/types"
	"golang.org/x/tools/go/loader"
	"log"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	serializationPackageName      = "pbgraphserialization"
	serializationProtoPackageName = "pbgraphserialization_pb"
	tagCommentEnable              = "pbgraphserialization:enable"
)

var invalidIdentifierChar = regexp.MustCompile("[^[:digit:][:alpha:]_]")

type Generator struct {
	targetPackagePath string
	targetPackageName string
	protoPackageName  string
	outputPath        string

	program                   *loader.Program
	serializationPackage      *types.Package
	serializationProtoPackage *types.Package
	targetPackage             *types.Package
	protoPackage              *types.Package

	serializableObjectType   *types.TypeName
	serializerType           *types.TypeName
	deserializerType         *types.TypeName
	protoMessageType         *types.TypeName
	protoSerializationIdType *types.TypeName

	buf                           bytes.Buffer
	targetObjectToProtoObjectInfo map[types.Object]*targetObjectToProtoObjectInfo

	localizationCache map[string]string
	packagePathToName map[string]string
	nameToPackagePath map[string]string

	packageRoots []string
}

type targetObjectToProtoObjectInfo struct {
	protoObject             types.Object
	targetFieldToProtoField map[*types.Var]*types.Var
}

func newTargetObjectToProtoObjectInfo(protoObject types.Object) *targetObjectToProtoObjectInfo {
	return &targetObjectToProtoObjectInfo{
		protoObject:             protoObject,
		targetFieldToProtoField: map[*types.Var]*types.Var{},
	}
}

func NewGenerator(targetPackagePath string, targetPackageName string, protoPackageName string, outputPath string) *Generator {
	var roots []string

	for _, root := range filepath.SplitList(build.Default.GOPATH) {
		roots = append(roots, filepath.Join(root, "src"))
	}

	generator := &Generator{
		targetPackagePath: targetPackagePath,
		targetPackageName: targetPackageName,
		protoPackageName:  protoPackageName,
		outputPath:        outputPath,

		targetObjectToProtoObjectInfo: map[types.Object]*targetObjectToProtoObjectInfo{},

		localizationCache: map[string]string{},
		packagePathToName: map[string]string{},
		nameToPackagePath: map[string]string{},

		packageRoots: roots,
	}

	return generator
}

func (generator *Generator) Generate() error {
	err := generator.loadCode()
	if err != nil {
		return err
	}

	err = generator.populateKnownTypes()
	if err != nil {
		return err
	}

	generator.addPackageImport(generator.protoMessageType.Pkg())
	generator.addPackageImport(generator.serializerType.Pkg())
	generator.generatePrologue()

	err = generator.populateTargetToProtoMaps()
	if err != nil {
		return err
	}

	err = generator.generateCode()
	if err != nil {
		return err
	}

	formatted, err := format.Source(generator.buf.Bytes())
	if err != nil {
		fmt.Println(string(generator.buf.Bytes()))
		log.Fatal(err)
	}
	fmt.Println(string(formatted))

	return nil
}

func (generator *Generator) loadCode() error {
	conf := loader.Config{
		AllowErrors: true,
		ParserMode:  parser.ParseComments,
	}

	conf.ImportWithTests(generator.targetPackagePath)
	program, err := conf.Load()
	if err != nil {
		return err
	}

	generator.program = program

	for pkg := range generator.program.AllPackages {
		if pkg.Name() == generator.targetPackageName {
			generator.targetPackage = pkg
		}

		if pkg.Name() == generator.protoPackageName {
			generator.protoPackage = pkg
		}

		if pkg.Name() == serializationPackageName {
			generator.serializationPackage = pkg
		}

		if pkg.Name() == serializationProtoPackageName {
			generator.serializationProtoPackage = pkg
		}
	}

	if generator.targetPackage == nil {
		return fmt.Errorf("target package '%s' not found", generator.targetPackageName)
	}

	if generator.protoPackage == nil {
		return fmt.Errorf("proto package '%s' not found", generator.protoPackageName)
	}

	if generator.targetPackage == nil {
		return fmt.Errorf("target package '%s' not found", generator.targetPackageName)
	}

	if generator.serializationPackage == nil {
		return fmt.Errorf("serialization package '%s' not found", serializationPackageName)
	}

	if generator.serializationProtoPackage == nil {
		return fmt.Errorf("serialization proto package '%s' not found", serializationProtoPackageName)
	}

	return nil
}

func (generator *Generator) populateKnownTypes() error {
	serializableObjectType, err := generator.getKnownType(serializationPackageName, "SerializableObject", true)
	if err != nil {
		return err
	}
	generator.serializableObjectType = serializableObjectType

	serializerType, err := generator.getKnownType(serializationPackageName, "Serializer", true)
	if err != nil {
		return err
	}
	generator.serializerType = serializerType

	deserializerType, err := generator.getKnownType(serializationPackageName, "Deserializer", true)
	if err != nil {
		return err
	}
	generator.deserializerType = deserializerType

	protoMessageType, err := generator.getKnownType("github.com/gogo/protobuf/proto", "Message", false)
	if err != nil {
		return err
	}
	generator.protoMessageType = protoMessageType

	protoSerializationIdType, err := generator.getKnownType(serializationProtoPackageName, "SerializationId", true)
	if err != nil {
		return err
	}
	generator.protoSerializationIdType = protoSerializationIdType

	return nil
}

func (generator *Generator) populateTargetToProtoMaps() error {
	enabledTypes := generator.getEnabledTypes()

	protoObjects := getObjectsFromScope(generator.protoPackage.Scope())
	for _, targetObject := range enabledTypes {
		protoObjectFound := false
		for _, protoObject := range protoObjects {
			targetTypeName := targetObject.(*types.TypeName)
			protoTypeName, protoObjectIsTypeName := protoObject.(*types.TypeName)
			if !protoObjectIsTypeName {
				continue
			}

			if generator.checkIfTypeNamesMatch(targetTypeName, protoTypeName) {
				generator.targetObjectToProtoObjectInfo[targetObject] = newTargetObjectToProtoObjectInfo(protoObject)
				protoObjectFound = true
				break
			}
		}

		if !protoObjectFound {
			return fmt.Errorf("no matching proto type found for serializable type '%s'", targetObject.Name())
		}
	}

	for targetObject, protoObjectInfo := range generator.targetObjectToProtoObjectInfo {
		targetStruct := targetObject.Type().Underlying().(*types.Struct)
		protoStruct := protoObjectInfo.protoObject.Type().Underlying().(*types.Struct)

		for i := 0; i < targetStruct.NumFields(); i++ {
			targetField := targetStruct.Field(i)
			if !generator.checkIfFieldSerialized(targetField) {
				continue
			}

			protoFieldFound := false
			for j := 0; j < protoStruct.NumFields(); j++ {
				protoField := protoStruct.Field(j)
				if generator.checkIfFieldsMatch(targetField, protoField) {
					protoObjectInfo.targetFieldToProtoField[targetField] = protoField
					protoFieldFound = true
					break
				}
			}

			if !protoFieldFound {
				return fmt.Errorf("no matching proto field found for serializable field '%s'", targetField.Name())
			}
		}
	}

	return nil
}

func (generator *Generator) checkIfTypeNamesMatch(targetTypeName, protoTypeName *types.TypeName) bool {
	return strings.EqualFold(targetTypeName.Name(), protoTypeName.Name())
}

func (generator *Generator) checkIfFieldsMatch(targetField, protoField *types.Var) bool {
	return strings.EqualFold(targetField.Name(), protoField.Name())
}

func (generator *Generator) checkIfFieldSerialized(targetField *types.Var) bool {
	return true
}

func makeValidVariableName(s string) string {
	return strings.Replace(invalidIdentifierChar.ReplaceAllString(s, ""), "_", "", -1)
}

func (generator *Generator) serializeVariable(protoVarName string, targetVarName string, targetVarType types.Type, protoVarType types.Type) error {
	if reflect.TypeOf(protoVarType) != reflect.TypeOf(targetVarType) {
		return fmt.Errorf(
			"variable type mismatch, target has '%v', proto has '%v' [target var name '%s', proto var name '%s']",
			reflect.TypeOf(targetVarType).Elem(),
			reflect.TypeOf(protoVarType).Elem(),
			targetVarName,
			protoVarName,
		)
	}

	switch protoVarType.(type) {
	case *types.Basic:
		protoVarBasic := protoVarType.(*types.Basic)
		generator.printlnf("%s = %s(%s)", protoVarName, generator.renderType(protoVarBasic), targetVarName)
	case *types.Named:
		protoVarNamed := protoVarType.(*types.Named)
		generator.printlnf("%s = %s(%s)", protoVarName, generator.renderType(protoVarNamed), targetVarName)
	case *types.Pointer:
		protoVarPointer := protoVarType.(*types.Pointer)
		protoVarPointerElemNamed := protoVarPointer.Elem().(*types.Named)
		if protoVarPointerElemNamed.Obj() != generator.protoSerializationIdType {
			return fmt.Errorf("unsupported reference type %v", protoVarPointerElemNamed.Obj())
		}

		generator.printlnf("%s = serializer.Serialize(%s).Marshal()", protoVarName, targetVarName)
	case *types.Slice:
		targetSliceElementName := hideName(makeValidVariableName(targetVarName) + "Element")
		targetVarSlice, ok := targetVarType.(*types.Slice)
		if !ok {
			return fmt.Errorf(
				"variable type mismatch, target has '%v', proto has '%v' [target var name '%s', proto var name '%s']",
				reflect.TypeOf(targetVarType).Elem(),
				reflect.TypeOf(protoVarType).Elem(),
				targetVarName,
				protoVarName,
			)
		}

		targetSliceElementType := targetVarSlice.Elem()

		protoSliceElementName := targetSliceElementName + "Serialized"
		protoSliceElementType := protoVarType.(*types.Slice).Elem()

		basicSliceHandled := false
		switch protoSliceElementType.(type) {
		case *types.Basic:
			if protoSliceElementType == targetSliceElementType {
				generator.println()
				generator.printlnf("%s = make([]%s, len(%s))", protoVarName, generator.renderType(targetSliceElementType), targetVarName)
				generator.printlnf("copy(%s, %s)", protoVarName, targetVarName)
				generator.println()

				basicSliceHandled = true
			}
		}

		if !basicSliceHandled {
			generator.println()
			generator.printlnf("for _, %s := range %s {", targetSliceElementName, targetVarName)
			generator.printlnf("var %s %s", protoSliceElementName, generator.renderType(protoSliceElementType))
			err := generator.serializeVariable(protoSliceElementName, targetSliceElementName, targetSliceElementType, protoSliceElementType)
			if err != nil {
				return err
			}
			generator.printlnf("%s = append(%s, %s)", protoVarName, protoVarName, protoSliceElementName)
			generator.printlnf("}")
			generator.println()
		}
	default:
		return fmt.Errorf("unsupported type: %#v (%T)", protoVarType, protoVarType)
	}

	return nil
}

func (generator *Generator) generateCode() error {
	for _, targetObject := range generator.sortedTargetObjectToProtoObjectInfo() {
		protoObjectInfo := generator.targetObjectToProtoObjectInfo[targetObject]
		thisName := hideName(lowerFirst(targetObject.Name()))

		// serializer code
		instanceName := hideName("serialized")
		generator.printlnf(
			"func (%s %s) Serialize(serializer %s) %s {",
			thisName,
			generator.renderType(types.NewPointer(targetObject.Type())),
			generator.renderType(types.NewPointer(generator.serializerType.Type())),
			generator.renderType(generator.protoMessageType.Type()),
		)
		generator.printlnf("%s := &%s{}", instanceName, generator.renderType(protoObjectInfo.protoObject.Type()))

		for _, targetField := range protoObjectInfo.sortedTargetFieldToProtoField() {
			protoField := protoObjectInfo.targetFieldToProtoField[targetField]
			err := generator.serializeVariable(instanceName+"."+protoField.Name(), thisName+"."+targetField.Name(), targetField.Type(), protoField.Type())
			if err != nil {
				return err
			}
		}

		generator.printlnf("return %s", instanceName)
		generator.printlnf("}")
		generator.println()

		// deserializer code
		instanceName = hideName("deserialized")
		generator.printlnf(
			"func (%s %s) Deserialize(deserializer %s, rawMessage %s) (%s, error) {",
			thisName,
			generator.renderType(types.NewPointer(targetObject.Type())),
			generator.renderType(types.NewPointer(generator.deserializerType.Type())),
			generator.renderType(generator.protoMessageType.Type()),
			generator.renderType(generator.serializableObjectType.Type()),
		)

		generator.printlnf("_message := rawMessage.(%s)", generator.renderType(types.NewPointer(protoObjectInfo.protoObject.Type())))
		generator.printlnf("return %s, nil", thisName)
		generator.printlnf("}")
		generator.println()
	}

	return nil
}

func (generator *Generator) generatePrologue() {
	generator.printlnf("// Auto-generated serialization code generated by pbgraphserialization-gen. DO NOT MODIFY!")
	generator.printf("package %s\n\n", generator.targetPackage.Name())

	generator.generateImports()
	generator.printf("\n")
}

func (generator *Generator) generateImports() {
	targetPackagePath := generator.nameToPackagePath[generator.targetPackage.Name()]

	// Sort by import name so that we get a deterministic order
	for _, name := range generator.sortedImportNames() {
		path := generator.nameToPackagePath[name]
		if path == targetPackagePath {
			continue
		}
		generator.printf("import %s \"%s\"\n", name, path)
	}
}

func (generator *Generator) getKnownType(packageIdentifier string, typeName string, searchByPackageName bool) (*types.TypeName, error) {
	packageFound := false
	for pkg := range generator.program.AllPackages {
		var compared string
		if searchByPackageName {
			compared = pkg.Name()
		} else {
			compared = pkg.Path()
		}

		if compared != packageIdentifier {
			continue
		}

		packageFound = true
		objects := getObjectsFromScope(pkg.Scope())
		for _, object := range objects {
			switch object.(type) {
			case *types.TypeName:
				if object.Name() == typeName {
					return object.(*types.TypeName), nil
				}
			}
		}
	}

	if packageFound {
		return nil, fmt.Errorf("known type '%s' not found in package '%s'", typeName, packageIdentifier)
	} else {
		return nil, fmt.Errorf("known package '%s' not found", packageIdentifier)
	}
}

func (generator *Generator) getEnabledTypes() []types.Object {
	targetPackage := generator.targetPackage
	objects := getObjectsFromScope(targetPackage.Scope())

	typeTest := func(object types.Object) bool {
		switch object.(type) {
		case *types.TypeName:
			_, path, _ := generator.program.PathEnclosingInterval(object.Pos(), object.Pos())
			for _, pathNode := range path {
				switch n := pathNode.(type) {
				case *ast.GenDecl:
					comment := n.Doc.Text()

					scanner := bufio.NewScanner(strings.NewReader(comment))
					for scanner.Scan() {
						commentLine := strings.TrimSpace(scanner.Text())
						if commentLine == tagCommentEnable {
							return true
						}
					}
				}
			}

			return false
		default:
			return false
		}
	}

	tempObjects := objects
	objects = objects[:0]
	for _, object := range tempObjects {
		if typeTest(object) {
			objects = append(objects, object)
		}
	}

	return objects
}

func (generator *Generator) renderType(typ types.Type) string {
	switch t := typ.(type) {
	case *types.Named:
		o := t.Obj()
		if o.Pkg() == nil || o.Pkg().Name() == "main" || (o.Pkg() == generator.targetPackage) {
			return o.Name()
		}
		return generator.addPackageImport(o.Pkg()) + "." + o.Name()
	case *types.Basic:
		return t.Name()
	case *types.Pointer:
		return "*" + generator.renderType(t.Elem())
	case *types.Slice:
		return "[]" + generator.renderType(t.Elem())
	case *types.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), generator.renderType(t.Elem()))
	/*case *types.Signature:
		switch t.Results().Len() {
		case 0:
			return fmt.Sprintf(
				"func(%s)",
				generator.renderTypeTuple(t.Params()),
			)
		case 1:
			return fmt.Sprintf(
				"func(%s) %s",
				generator.renderTypeTuple(t.Params()),
				generator.renderType(t.Results().At(0).Type()),
			)
		default:
			return fmt.Sprintf(
				"func(%s)(%s)",
				generator.renderTypeTuple(t.Params()),
				generator.renderTypeTuple(t.Results()),
			)
		}*/
	case *types.Map:
		kt := generator.renderType(t.Key())
		vt := generator.renderType(t.Elem())

		return fmt.Sprintf("map[%s]%s", kt, vt)
	case *types.Chan:
		switch t.Dir() {
		case types.SendRecv:
			return "chan " + generator.renderType(t.Elem())
		case types.RecvOnly:
			return "<-chan " + generator.renderType(t.Elem())
		default:
			return "chan<- " + generator.renderType(t.Elem())
		}
	case *types.Struct:
		var fields []string

		for i := 0; i < t.NumFields(); i++ {
			f := t.Field(i)

			if f.Anonymous() {
				fields = append(fields, generator.renderType(f.Type()))
			} else {
				fields = append(fields, fmt.Sprintf("%s %s", f.Name(), generator.renderType(f.Type())))
			}
		}

		return fmt.Sprintf("struct{%s}", strings.Join(fields, ";"))
	case *types.Interface:
		if t.NumMethods() != 0 {
			panic("Unable to mock inline interfaces with methods")
		}

		return "interface{}"
	default:
		panic(fmt.Sprintf("un-namable type: %#v (%T)", t, t))
	}
}

func (generator *Generator) sortedTargetObjectToProtoObjectInfo() (targetObjects []types.Object) {
	for key := range generator.targetObjectToProtoObjectInfo {
		targetObjects = append(targetObjects, key)
	}
	sort.SliceStable(targetObjects, func(i, j int) bool {
		return targetObjects[i].Name() > targetObjects[j].Name()
	})

	return
}

func (generator *Generator) sortedImportNames() (importNames []string) {
	for name := range generator.nameToPackagePath {
		importNames = append(importNames, name)
	}
	sort.Strings(importNames)
	return
}

func (generator *Generator) addPackageImport(pkg *types.Package) string {
	return generator.addPackageImportWithName(pkg.Path(), pkg.Name())
}

func (generator *Generator) addPackageImportWithName(path, name string) string {
	path = generator.getLocalizedPath(path)
	if existingName, pathExists := generator.packagePathToName[path]; pathExists {
		return existingName
	}

	nonConflictingName := generator.getNonConflictingName(path, name)
	generator.packagePathToName[path] = nonConflictingName
	generator.nameToPackagePath[nonConflictingName] = path
	return nonConflictingName
}

func (generator *Generator) getNonConflictingName(path, name string) string {
	if !generator.importNameExists(name) {
		return name
	}

	// The path will always contain '/' because it is enforced in getLocalizedPath
	// regardless of OS.
	directories := strings.Split(path, "/")

	cleanedDirectories := make([]string, 0, len(directories))
	for _, directory := range directories {
		cleaned := invalidIdentifierChar.ReplaceAllString(directory, "_")
		cleanedDirectories = append(cleanedDirectories, cleaned)
	}
	numDirectories := len(cleanedDirectories)
	var prospectiveName string
	for i := 1; i <= numDirectories; i++ {
		prospectiveName = strings.Join(cleanedDirectories[numDirectories-i:], "")
		if !generator.importNameExists(prospectiveName) {
			return prospectiveName
		}
	}
	// Try adding numbers to the given name
	i := 2
	for {
		prospectiveName = fmt.Sprintf("%v%d", name, i)
		if !generator.importNameExists(prospectiveName) {
			return prospectiveName
		}
		i++
	}
}

func (generator *Generator) importNameExists(name string) bool {
	_, nameExists := generator.nameToPackagePath[name]
	return nameExists
}

func (generator *Generator) getLocalizedPathFromPackage(pkg *types.Package) string {
	return generator.getLocalizedPath(pkg.Path())
}

// TODO(@IvanMalison): Is there not a better way to get the actual
// import path of a package?
func (generator *Generator) getLocalizedPath(path string) string {
	if strings.HasSuffix(path, ".go") {
		path, _ = filepath.Split(path)
	}
	if localized, ok := generator.localizationCache[path]; ok {
		return localized
	}
	directories := strings.Split(path, string(filepath.Separator))
	numDirectories := len(directories)
	vendorIndex := -1
	for i := 1; i <= numDirectories; i++ {
		dir := directories[numDirectories-i]
		if dir == "vendor" {
			vendorIndex = numDirectories - i
			break
		}
	}

	toReturn := path
	if vendorIndex >= 0 {
		toReturn = filepath.Join(directories[vendorIndex+1:]...)
	} else if filepath.IsAbs(path) {
		toReturn = calculateImport(generator.packageRoots, path)
	}

	// Enforce '/' slashes for import paths in every OS.
	toReturn = filepath.ToSlash(toReturn)

	generator.localizationCache[path] = toReturn
	return toReturn
}

func (generator *Generator) printf(format string, vals ...interface{}) {
	fmt.Fprintf(&generator.buf, format, vals...)
}

func (generator *Generator) printlnf(format string, vals ...interface{}) {
	generator.printf(format, vals...)
	fmt.Fprintln(&generator.buf)
}

func (generator *Generator) println() {
	fmt.Fprintln(&generator.buf)
}

func (info *targetObjectToProtoObjectInfo) sortedTargetFieldToProtoField() (targetFields []*types.Var) {
	for key := range info.targetFieldToProtoField {
		targetFields = append(targetFields, key)
	}
	sort.SliceStable(targetFields, func(i, j int) bool {
		score := calculateTypeScore(targetFields[i].Type()) - calculateTypeScore(targetFields[j].Type())
		if score == 0 {
			return targetFields[i].Name() > targetFields[j].Name()
		}
		return score > 0
	})

	return
}

func getObjectsFromScope(scope *types.Scope) []types.Object {
	names := scope.Names()
	objects := make([]types.Object, 0, len(names))
	for _, name := range names {
		objects = append(objects, scope.Lookup(name))
	}

	return objects
}

func calculateImport(set []string, path string) string {
	for _, root := range set {
		if strings.HasPrefix(path, root) {
			packagePath, err := filepath.Rel(root, path)
			if err == nil {
				return packagePath
			} else {
				log.Printf("Unable to localize path %v, %v", path, err)
			}
		}
	}
	return path
}

func calculateTypeScore(typ types.Type) int {
	switch typ.(type) {
	case *types.Basic:
		return 10
	case *types.Named:
		return 9
	case *types.Pointer:
		return 8
	default:
		return 0
	}
}

func hideName(s string) string {
	return "__" + s
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

func isNillable(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Pointer, *types.Array, *types.Map, *types.Interface, *types.Signature, *types.Chan, *types.Slice:
		return true
	case *types.Named:
		return isNillable(t.Underlying())
	}
	return false
}

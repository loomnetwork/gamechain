package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
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
	tagCommentSkip                = "pbgraphserialization:skip"
	tagCommentRoot                = "pbgraphserialization:root"
)

var invalidIdentifierChar = regexp.MustCompile("[^[:digit:][:alpha:]_]")

type Generator struct {
	targetPackagePath string
	targetPackageName string
	protoPackageName  string

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
	protoSerializedGraphType *types.TypeName

	buf                           *bytes.Buffer
	enabledObjectsData            map[types.Object]*TargetObjectMetadata
	targetObjectToProtoObjectInfo map[types.Object]*targetObjectToProtoObjectInfo

	localizationCache map[string]string
	packagePathToName map[string]string
	nameToPackagePath map[string]string

	packageRoots []string
}

type TargetObjectMetadata struct {
	generateRootSerializationMethods bool
	skip                             bool
}

type targetObjectToProtoObjectInfo struct {
	targetObjectMetadata    *TargetObjectMetadata
	protoObject             types.Object
	targetFieldToProtoField map[*types.Var]*types.Var
}

func newTargetObjectToProtoObjectInfo(protoObject types.Object, targetObjectMetadata *TargetObjectMetadata) *targetObjectToProtoObjectInfo {
	return &targetObjectToProtoObjectInfo{
		targetObjectMetadata:    targetObjectMetadata,
		protoObject:             protoObject,
		targetFieldToProtoField: map[*types.Var]*types.Var{},
	}
}

func NewGenerator(targetPackagePath string, targetPackageName string, protoPackageName string) (*Generator, error) {
	var roots []string

	for _, root := range filepath.SplitList(build.Default.GOPATH) {
		roots = append(roots, filepath.Join(root, "src"))
	}

	generator := &Generator{
		targetPackagePath: targetPackagePath,
		targetPackageName: targetPackageName,
		protoPackageName:  protoPackageName,

		enabledObjectsData: map[types.Object]*TargetObjectMetadata{},
		targetObjectToProtoObjectInfo: map[types.Object]*targetObjectToProtoObjectInfo{},

		localizationCache: map[string]string{},
		packagePathToName: map[string]string{},
		nameToPackagePath: map[string]string{},

		packageRoots: roots,
	}

	err := generator.loadCode()
	if err != nil {
		return nil, err
	}

	err = generator.populateKnownTypes()
	if err != nil {
		return nil, err
	}

	generator.addPackageImport(generator.protoMessageType.Pkg())
	generator.addPackageImport(generator.serializerType.Pkg())
	generator.addPackageImport(generator.protoSerializationIdType.Pkg())
	generator.addPackageImport(generator.protoSerializedGraphType.Pkg())

	return generator, nil
}

func (generator *Generator) AddEnabledType(packageIdentifier string, typeName string, searchByPackageName bool, metadata *TargetObjectMetadata) error {
	typeNameInstance, err := generator.getType(packageIdentifier, typeName, searchByPackageName)
	if err != nil {
		return err
	}

	generator.enabledObjectsData[typeNameInstance] = metadata
	return nil
}

func (generator *Generator) AddEnabledTypesFromCode() error {
	objectsData, err := generator.getEnabledObjectsFromCode()
	if err != nil {
		return err
	}

	for k, v := range objectsData {
		generator.enabledObjectsData[k] = v
	}

	return nil
}

func (generator *Generator) Generate() (code string, err error) {
	generator.buf = &bytes.Buffer{}
	generator.targetObjectToProtoObjectInfo = map[types.Object]*targetObjectToProtoObjectInfo{}

	err = generator.populateTargetToProtoMaps()
	if err != nil {
		return "", err
	}

	generator.populateImports()

	generator.generatePrologue()
	err = generator.generateCode()
	if err != nil {
		return "", err
	}

	formatted, err := format.Source(generator.buf.Bytes())
	if err != nil {
		return string(generator.buf.Bytes()), err
	}

	return string(formatted), nil
}

func (generator *Generator) populateImports() {
	for k, v := range generator.targetObjectToProtoObjectInfo {
		generator.addPackageImport(k.Pkg())
		generator.addPackageImport(v.protoObject.Pkg())

		for k2, v2 := range v.targetFieldToProtoField {
			generator.addPackageImport(k2.Pkg())
			generator.addPackageImport(v2.Pkg())
		}
	}
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
	serializableObjectType, err := generator.getType(serializationPackageName, "SerializableObject", true)
	if err != nil {
		return err
	}
	generator.serializableObjectType = serializableObjectType

	serializerType, err := generator.getType(serializationPackageName, "Serializer", true)
	if err != nil {
		return err
	}
	generator.serializerType = serializerType

	deserializerType, err := generator.getType(serializationPackageName, "Deserializer", true)
	if err != nil {
		return err
	}
	generator.deserializerType = deserializerType

	protoMessageType, err := generator.getType("github.com/gogo/protobuf/proto", "Message", false)
	if err != nil {
		return err
	}
	generator.protoMessageType = protoMessageType

	protoSerializationIdType, err := generator.getType(serializationProtoPackageName, "SerializationId", true)
	if err != nil {
		return err
	}
	generator.protoSerializationIdType = protoSerializationIdType

	protoSerializedGraphType, err := generator.getType(serializationProtoPackageName, "SerializedGraph", true)
	if err != nil {
		return err
	}
	generator.protoSerializedGraphType = protoSerializedGraphType

	return nil
}

func (generator *Generator) populateTargetToProtoMaps() error {
	protoObjects := getObjectsFromScope(generator.protoPackage.Scope())
	for targetObject, targetObjectMetadata := range generator.enabledObjectsData {
		if targetObjectMetadata.skip {
			continue
		}

		protoObjectFound := false
		for _, protoObject := range protoObjects {
			targetTypeName := targetObject.(*types.TypeName)
			protoTypeName, protoObjectIsTypeName := protoObject.(*types.TypeName)
			if !protoObjectIsTypeName {
				continue
			}

			if generator.checkIfTypeNamesMatch(targetTypeName, protoTypeName) {
				generator.targetObjectToProtoObjectInfo[targetObject] = newTargetObjectToProtoObjectInfo(protoObject, targetObjectMetadata)
				protoObjectFound = true
				break
			}
		}

		if !protoObjectFound {
			return fmt.Errorf("no matching proto type found for serializable type '%s'", generator.renderType(targetObject.Type()))
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
			var protoField *types.Var
			for j := 0; j < protoStruct.NumFields(); j++ {
				protoField = protoStruct.Field(j)
				if generator.checkIfFieldsMatch(targetField, protoField) {
					protoObjectInfo.targetFieldToProtoField[targetField] = protoField
					protoFieldFound = true
					break
				}
			}

			if !protoFieldFound {
				return fmt.Errorf("no matching proto field found for serializable field '%s.%s'", generator.renderType(targetObject.Type()), targetField.Name())
			}

			err := generator.checkIfFieldMatchIsValid(targetField, protoField)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to handle target object %s", generator.renderType(targetObject.Type())))
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

func (generator *Generator) checkIfFieldMatchIsValid(targetField *types.Var, protoField *types.Var) error {
	if targetField.Type() == generator.protoSerializationIdType.Type() {
		return fmt.Errorf("serializing '%s' (field '%s') doesn't makes sense, check your code", generator.renderType(targetField.Type()), targetField.Name())
	}

	return nil
}

func (generator *Generator) generateVarSerializationCode(serialization bool, destinationVarName string, sourceVarName string, sourceVarType types.Type, destinationVarType types.Type) error {
	var errorTag, elementSuffix string
	if serialization {
		errorTag = "serialization: "
		elementSuffix = "Serialized"
	} else {
		errorTag = "deserialization: "
		elementSuffix = "Deserialized"
	}
	if reflect.TypeOf(destinationVarType) != reflect.TypeOf(sourceVarType) {
		return generator.newTypeMismatchError(errorTag, destinationVarName, sourceVarName, sourceVarType, destinationVarType)
	}

	switch destinationVarType.(type) {
	case *types.Basic, *types.Named:
		switch destinationVarType.(type) {
		case *types.Named:
			destinationVarNamed := destinationVarType.(*types.Named)
			if destinationVarNamed.Obj() == generator.protoSerializationIdType {
				return createSerializationCodeGeneratorError(errorTag, destinationVarName, sourceVarName, "non-pointer SerializationId is not supported")
			}
		}

		assignable := types.ConvertibleTo(sourceVarType, destinationVarType)
		if !assignable {
			return createSerializationCodeGeneratorError(errorTag, destinationVarName, sourceVarName, "type '%s' is not convertible to type '%s'", generator.renderType(sourceVarType), generator.renderType(destinationVarType))
		}

		generator.printlnf("%s = %s(%s)", destinationVarName, generator.renderType(destinationVarType), sourceVarName)
	case *types.Pointer:
		destinationVarPointer := destinationVarType.(*types.Pointer)
		destinationVarPointerElemNamed := destinationVarPointer.Elem().(*types.Named)
		if serialization {
			if destinationVarPointerElemNamed.Obj() != generator.protoSerializationIdType {
				return createSerializationCodeGeneratorError(errorTag, destinationVarName, sourceVarName, "unsupported reference type %s", destinationVarPointerElemNamed.Obj().Type().String())
			}

			generator.printlnf("%s = serializer.Serialize(%s).Serialize()", destinationVarName, sourceVarName)
		} else {
			errName := hideName("err")
			destinationTempVarName := makeValidVariableName(destinationVarName) + elementSuffix
			sourceRealVarType, ok := generator.targetObjectToProtoObjectInfo[destinationVarPointerElemNamed.Obj()]
			if !ok {
				return createSerializationCodeGeneratorError(errorTag, destinationVarName, sourceVarName, "var has non-serialized type %s", generator.renderType(destinationVarPointerElemNamed))
			}

			generator.printlnf("%s, %s := deserializer.Deserialize(", destinationTempVarName, errName)
			generator.printlnf("%s,", sourceVarName)

			generator.printlnf("func() %s { return &%s{} },", generator.renderType(generator.serializableObjectType.Type()), generator.renderType(destinationVarPointerElemNamed))
			generator.printlnf("func() %s { return &%s{} },", generator.renderType(generator.protoMessageType.Type()), generator.renderType(sourceRealVarType.protoObject.Type()))
			generator.printlnf(")")
			generator.println()
			generator.printlnf("if %s != nil {", errName)
			generator.printlnf("return nil, %s", errName)
			generator.printlnf("}")
			generator.println()
			generator.printlnf("%s = %s.(%s)", destinationVarName, destinationTempVarName, generator.renderType(destinationVarType))
		}
	case *types.Slice:
		sourceSliceElementName := hideName(makeValidVariableName(sourceVarName) + "Element")
		sourceVarSlice, ok := sourceVarType.(*types.Slice)
		if !ok {
			return generator.newTypeMismatchError(errorTag, destinationVarName, sourceVarName, sourceVarType, destinationVarType)
		}

		sourceSliceElementType := sourceVarSlice.Elem()

		destinationSliceElementName := sourceSliceElementName + elementSuffix
		destinationSliceElementType := destinationVarType.(*types.Slice).Elem()

		basicSliceHandled := false

		// optimization for the case of primitive type slice with matching types
		switch destinationSliceElementType.(type) {
		case *types.Basic:
			if destinationSliceElementType == sourceSliceElementType {
				generator.println()

				// TODO: add unsafe option to re-use the reference instead of copying?
				generator.printlnf("%s = make([]%s, len(%s))", destinationVarName, generator.renderType(sourceSliceElementType), sourceVarName)
				generator.printlnf("copy(%s, %s)", destinationVarName, sourceVarName)
				generator.println()

				basicSliceHandled = true
			}
		}

		if !basicSliceHandled {
			generator.println()
			generator.printlnf("for _, %s := range %s {", sourceSliceElementName, sourceVarName)
			generator.printlnf("var %s %s", destinationSliceElementName, generator.renderType(destinationSliceElementType))
			err := generator.generateVarSerializationCode(serialization, destinationSliceElementName, sourceSliceElementName, sourceSliceElementType, destinationSliceElementType)
			if err != nil {
				return err
			}
			generator.printlnf("%s = append(%s, %s)", destinationVarName, destinationVarName, destinationSliceElementName)
			generator.printlnf("}")
			generator.println()
		}
	default:
		return createSerializationCodeGeneratorError(errorTag, destinationVarName, sourceVarName, "unsupported type: %#v (%T)", destinationVarType, destinationVarType)
	}

	return nil
}

func (generator *Generator) generateCode() error {
	wrapError := func(err error, prefix string, sourceVarType types.Type, destinationVarType types.Type) error {
		return errors.Wrap(err, fmt.Sprintf("%sfailed to generate code [source type %s, target type %s]", prefix, generator.renderType(sourceVarType), generator.renderType(destinationVarType)))
	}

	for _, targetObject := range generator.sortedTargetObjectToProtoObjectInfo() {
		protoObjectInfo := generator.targetObjectToProtoObjectInfo[targetObject]
		thisName := hideName(lowerFirst(targetObject.Name()))
		serializedName := hideName("serialized")

		sortedTargetFieldToProtoFieldKeys := protoObjectInfo.sortedTargetFieldToProtoField()

		// serializer code
		generator.printlnf(
			"func (%s %s) Serialize(serializer %s) %s {",
			thisName,
			generator.renderType(types.NewPointer(targetObject.Type())),
			generator.renderType(types.NewPointer(generator.serializerType.Type())),
			generator.renderType(generator.protoMessageType.Type()),
		)
		generator.printlnf("%s := &%s{}", serializedName, generator.renderType(protoObjectInfo.protoObject.Type()))

		for _, targetField := range sortedTargetFieldToProtoFieldKeys {
			protoField := protoObjectInfo.targetFieldToProtoField[targetField]
			destinationVarName := serializedName + "." + protoField.Name()
			sourceVarName := thisName + "." + targetField.Name()
			err := generator.generateVarSerializationCode(true, destinationVarName, sourceVarName, targetField.Type(), protoField.Type())
			if err != nil {
				return wrapError(err, "serialization: ", targetField.Type(), protoField.Type())
			}
		}

		generator.printlnf("return %s", serializedName)
		generator.printlnf("}")
		generator.println()

		if protoObjectInfo.targetObjectMetadata.generateRootSerializationMethods {
			serializerName := hideName("serializer")
			generator.printlnf(
				"func (%s %s) SerializeAsRoot() (%s, error) {",
				thisName,
				generator.renderType(types.NewPointer(targetObject.Type())),
				generator.renderType(types.NewPointer(generator.protoSerializedGraphType.Type())),
			)

			generator.printlnf("%s := %sNewSerializerSerialize(%s)", serializerName, generator.renderSerializationPackageAccessor(), thisName)
			generator.printlnf("return %s.SerializeToGraph()", serializerName)
			generator.printlnf("}")
			generator.println()
		}

		// deserializer code
		generator.printlnf(
			"func (%s %s) Deserialize(deserializer %s, rawMessage %s) (%s, error) {",
			thisName,
			generator.renderType(types.NewPointer(targetObject.Type())),
			generator.renderType(types.NewPointer(generator.deserializerType.Type())),
			generator.renderType(generator.protoMessageType.Type()),
			generator.renderType(generator.serializableObjectType.Type()),
		)

		if len(protoObjectInfo.targetFieldToProtoField) > 0 {
			generator.printlnf("%s := rawMessage.(%s)", serializedName, generator.renderType(types.NewPointer(protoObjectInfo.protoObject.Type())))
		}

		for _, targetField := range sortedTargetFieldToProtoFieldKeys {
			protoField := protoObjectInfo.targetFieldToProtoField[targetField]
			destinationVarName := thisName + "." + targetField.Name()
			sourceVarName := serializedName + "." + protoField.Name()
			err := generator.generateVarSerializationCode(false, destinationVarName, sourceVarName, protoField.Type(), targetField.Type())
			if err != nil {
				return wrapError(err, "deserialization: ", targetField.Type(), protoField.Type())
			}
		}

		generator.printlnf("return %s, nil", thisName)
		generator.printlnf("}")
		generator.println()

		if protoObjectInfo.targetObjectMetadata.generateRootSerializationMethods {
			generator.printlnf(
				"func Deserialize%sAsRoot(graph %s) (%s, error) {",
				targetObject.Name(),
				generator.renderType(types.NewPointer(generator.protoSerializedGraphType.Type())),
				generator.renderType(types.NewPointer(targetObject.Type())),
			)

			generator.printlnf("deserializer, err := %sNewDeserializerDeserializeFromGraph(graph)", generator.renderSerializationPackageAccessor())
			generator.printlnf("if err != nil {")
			generator.printlnf("return nil, err")
			generator.printlnf("}")
			generator.println()
			generator.printlnf(
				"deserialized, err := deserializer.DeserializeRoot(&%s{}, &%s{})",
				generator.renderType(targetObject.Type()),
				generator.renderType(generator.targetObjectToProtoObjectInfo[targetObject].protoObject.Type()),
			)

			generator.printlnf("if err != nil {")
			generator.printlnf("return nil, err")
			generator.printlnf("}")
			generator.println()

			generator.printlnf("return deserialized.(%s), nil", generator.renderType(types.NewPointer(targetObject.Type())), )
			generator.printlnf("}")
			generator.println()
		}
	}

	return nil
}

func (generator *Generator) generatePrologue() {
	generator.printlnf("// Generated by pbgraphserialization-gen. DO NOT EDIT.")
	generator.printlnf("// targetPackagePath: %s", generator.targetPackagePath)
	generator.printlnf("// targetPackageName: %s", generator.targetPackageName)
	generator.printlnf("// protoPackageName: %s", generator.protoPackageName)
	generator.println()

	generator.printf("package %s\n\n", generator.targetPackage.Name())

	generator.generateImports()
	generator.printf("\n")
}

func (generator *Generator) generateImports() {
	if len(generator.targetObjectToProtoObjectInfo) == 0 {
		return
	}

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

func (generator *Generator) getType(packageIdentifier string, typeName string, searchByPackageName bool) (*types.TypeName, error) {
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
		return nil, fmt.Errorf("type '%s' not found in package '%s'", typeName, packageIdentifier)
	} else {
		return nil, fmt.Errorf("package '%s' not found", packageIdentifier)
	}
}

func (generator *Generator) getEnabledObjectsFromCode() (map[types.Object]*TargetObjectMetadata, error) {
	targetPackage := generator.targetPackage
	objects := getObjectsFromScope(targetPackage.Scope())

	typeTest := func(object types.Object) (*TargetObjectMetadata, bool, error) {
		switch object.(type) {
		case *types.TypeName:
			_, path, _ := generator.program.PathEnclosingInterval(object.Pos(), object.Pos())
			if path == nil {
				return nil, false, fmt.Errorf("failed to find path for object %s", object.Type())
			}

			for _, pathNode := range path {
				switch n := pathNode.(type) {
				case *ast.GenDecl:
					comment := n.Doc.Text()

					scanner := bufio.NewScanner(strings.NewReader(comment))
					targetObjectMetadata := TargetObjectMetadata{}

					enabled := false
					for scanner.Scan() {
						commentLine := strings.TrimSpace(scanner.Text())
						switch commentLine {
						case tagCommentEnable:
							enabled = true
						case tagCommentRoot:
							targetObjectMetadata.generateRootSerializationMethods = true
						case tagCommentSkip:
							targetObjectMetadata.skip = true
						}
					}

					if enabled {
						return &targetObjectMetadata, true, nil
					}
				}
			}

			return nil, false, nil
		default:
			return nil, false, nil
		}
	}

	enabledObjects := map[types.Object]*TargetObjectMetadata{}
	for _, object := range objects {
		metadata, enabled, err := typeTest(object)
		if err != nil {
			return nil, err
		}
		if enabled {
			enabledObjects[object] = metadata
		}
	}

	return enabledObjects, nil
}

func (generator *Generator) renderSerializationPackageAccessor() string {
	if generator.targetPackage.Name() == serializationPackageName {
		return ""
	} else {
		return serializationPackageName + "."
	}
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
	_, _ = fmt.Fprintf(generator.buf, format, vals...)
}

func (generator *Generator) printlnf(format string, vals ...interface{}) {
	generator.printf(format, vals...)
	_, _ = fmt.Fprintln(generator.buf)
}

func (generator *Generator) println() {
	_, _ = fmt.Fprintln(generator.buf)
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

func (generator *Generator) newTypeMismatchError(prefix string, destinationVarName string, sourceVarName string, sourceVarType types.Type, destinationVarType types.Type) error {
	return fmt.Errorf(
		prefix+"variable type mismatch, source has '%v' (%s), destination has '%v' (%s) %s",
		reflect.TypeOf(sourceVarType).Elem(),
		generator.renderType(sourceVarType),
		reflect.TypeOf(destinationVarType).Elem(),
		generator.renderType(destinationVarType),
		createSerializationCodeGeneratorErrorExplanation(sourceVarName, destinationVarName),
	)
}

func createSerializationCodeGeneratorError(prefix, destinationVarName, sourceVarName, format string, vals ...interface{}) error {
	return fmt.Errorf(
		"%s%s %s",
		prefix,
		fmt.Sprintf(format, vals...),
		createSerializationCodeGeneratorErrorExplanation(sourceVarName, destinationVarName),
	)
}

func createSerializationCodeGeneratorErrorExplanation(destinationVarName string, sourceVarName string) string {
	return fmt.Sprintf(
		"[source var name '%s', destination var name '%s']",
		sourceVarName,
		destinationVarName,
	)
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

func makeValidVariableName(s string) string {
	return strings.Replace(invalidIdentifierChar.ReplaceAllString(s, ""), "_", "", -1)
}

func hideName(s string) string {
	return s
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

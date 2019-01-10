package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/types"
	"golang.org/x/tools/go/loader"
	"strings"
)

const (
	serializationPackageName      = "pbgraphserialization"
	serializationProtoPackageName = "pbgraphserialization_pb"
	tagCommentEnable              = "pbgraphserialization:enable"
)

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
}

func NewGenerator(targetPackagePath string, targetPackageName string, protoPackageName string, outputPath string) *Generator {
	generator := &Generator{
		targetPackagePath: targetPackagePath,
		targetPackageName: targetPackageName,
		protoPackageName:  protoPackageName,
		outputPath:        outputPath,
	}

	return generator
}

func (generator *Generator) Generate() error {
	err := generator.loadCode()
	if err != nil {
		return err
	}

	enabledTypes := generator.getEnabledTypes()

	// find types that have correspondence in proto package
	protoObjects := getObjectsFromScope(generator.protoPackage.Scope())
	targetObjectToProtoObject := make(map[types.Object]types.Object)
	for _, targetType := range enabledTypes {
		protoObjectFound := false
		for _, protoObject := range protoObjects {
			if targetType.Name() == protoObject.Name() {
				targetObjectToProtoObject[targetType] = protoObject
				protoObjectFound = true
				break
			}
		}

		if !protoObjectFound {
			return fmt.Errorf("no matching proto type found for serializable type '%s'", targetType.Name())
		}
	}

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
					comment := strings.TrimSpace(n.Doc.Text())
					if comment == tagCommentEnable {
						return true
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

func getObjectsFromScope(scope *types.Scope) []types.Object {
	names := scope.Names()
	objects := make([]types.Object, 0, len(names))
	for _, name := range names {
		objects = append(objects, scope.Lookup(name))
	}

	return objects
}

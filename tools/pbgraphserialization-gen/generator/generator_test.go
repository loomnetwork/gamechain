package generator

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

const (
	debugPrintCode = true
)

func TestGraphSerializationGenerator_1(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledTypesFromCode()
	assert.NoError(t, err)
	code, err := generator.Generate()
	assert.NoError(t, err)

	printCode(code)
}

func TestGraphSerializationGenerator_Skip(t *testing.T) {
	generator := createTestGenerator(t)

	err := generator.AddEnabledType("pbgraphserialization", "CardAbility", true, &TargetObjectMetadata{skip: false})
	assert.NoError(t, err)
	codeNoSkip, err := generator.Generate()
	assert.NoError(t, err)
	printCode(codeNoSkip)

	err = generator.AddEnabledType("pbgraphserialization", "CardAbility", true, &TargetObjectMetadata{skip: true})
	assert.NoError(t, err)
	codeSkip, err := generator.Generate()
	assert.NoError(t, err)
	printCode(codeSkip)

	codeNoSkipAst := parseCode(codeNoSkip)
	codeSkipAst := parseCode(codeSkip)

	assert.NotNil(t, findFuncDecl(codeNoSkipAst, "CardAbility", "Serialize"))
	assert.NotNil(t, findFuncDecl(codeNoSkipAst, "CardAbility", "Deserialize"))

	assert.Nil(t, findFuncDecl(codeSkipAst, "CardAbility", "Serialize"))
	assert.Nil(t, findFuncDecl(codeSkipAst, "CardAbility", "Deserialize"))
}

func TestGraphSerializationGenerator_GenerateRootMethods(t *testing.T) {
	generator := createTestGenerator(t)

	err := generator.AddEnabledType("pbgraphserialization", "CardAbility", true, &TargetObjectMetadata{generateRootSerializationMethods: true})
	assert.NoError(t, err)
	codeRoot, err := generator.Generate()
	assert.NoError(t, err)
	printCode(codeRoot)

	err = generator.AddEnabledType("pbgraphserialization", "CardAbility", true, &TargetObjectMetadata{generateRootSerializationMethods: false})
	assert.NoError(t, err)
	codeNoRoot, err := generator.Generate()
	assert.NoError(t, err)
	printCode(codeNoRoot)

	codeRootAst := parseCode(codeRoot)
	codeNoRootAst := parseCode(codeNoRoot)

	assert.NotNil(t, findFuncDecl(codeRootAst, "CardAbility", "Serialize"))
	assert.NotNil(t, findFuncDecl(codeRootAst, "CardAbility", "Deserialize"))
	assert.NotNil(t, findFuncDecl(codeRootAst, "CardAbility", "SerializeAsRoot"))
	assert.NotNil(t, findFuncDecl(codeRootAst, "", "DeserializeCardAbilityAsRoot"))

	assert.NotNil(t, findFuncDecl(codeNoRootAst, "CardAbility", "Serialize"))
	assert.NotNil(t, findFuncDecl(codeNoRootAst, "CardAbility", "Deserialize"))
	assert.Nil(t, findFuncDecl(codeNoRootAst, "CardAbility", "SerializeAsRoot"))
	assert.Nil(t, findFuncDecl(codeNoRootAst, "", "DeserializeCardAbilityAsRoot"))
}

func TestGraphSerializationGenerator_EmptyType(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "EmptyType", true, &TargetObjectMetadata{generateRootSerializationMethods:true})
	assert.NoError(t, err)
	code, err := generator.Generate()
	assert.NoError(t, err)

	printCode(code)
	codeAst := parseCode(code)

	assert.NotNil(t, findFuncDecl(codeAst, "EmptyType", "Serialize"))
	assert.NotNil(t, findFuncDecl(codeAst, "EmptyType", "Deserialize"))
	assert.NotNil(t, findFuncDecl(codeAst, "EmptyType", "SerializeAsRoot"))
	assert.NotNil(t, findFuncDecl(codeAst, "", "DeserializeEmptyTypeAsRoot"))
}

func TestGraphSerializationGenerator_AddNonExistentType(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "NonExistent", true, &TargetObjectMetadata{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type 'NonExistent' not found in package 'pbgraphserialization'")
}

func TestGraphSerializationGenerator_AddTypeFromNonExistentPackage(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("NonExistent", "", true, &TargetObjectMetadata{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "package 'NonExistent' not found")
}

func TestGraphSerializationGenerator_NoMatchingProtoType(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "TypeWithNoMatchingProtoType", true, &TargetObjectMetadata{})
	assert.NoError(t, err)
	_, err = generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no matching proto type found for serializable type 'TypeWithNoMatchingProtoType'")
}

func TestGraphSerializationGenerator_FailSerializeSerializationId(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "InvalidType1", true, &TargetObjectMetadata{})
	assert.NoError(t, err)
	_, err = generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "serializing 'pbgraphserialization_pb.SerializationId' (field 'card') doesn't makes sense")
}

func TestGraphSerializationGenerator_VariableTypeMismatch(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "InvalidType2", true, &TargetObjectMetadata{})
	assert.NoError(t, err)
	_, err = generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variable type mismatch, source has 'types.Named' (Card), destination has 'types.Pointer'")
}

func TestGraphSerializationGenerator_NonPointerSerializationId(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "InvalidType3", true, &TargetObjectMetadata{})
	assert.NoError(t, err)
	_, err = generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non-pointer SerializationId is not supported")
}

func TestGraphSerializationGenerator_UnsupportedReferenceType(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "InvalidType4", true, &TargetObjectMetadata{})
	assert.NoError(t, err)
	_, err = generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported reference type")
}

func TestGraphSerializationGenerator_MapUnsupportedType(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "InvalidType5", true, &TargetObjectMetadata{})
	assert.NoError(t, err)
	_, err = generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type: &types.Map")
}

func TestGraphSerializationGenerator_NonConvertibleType(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "InvalidType6", true, &TargetObjectMetadata{})
	assert.NoError(t, err)
	code, err := generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "serialization: type 'string' is not convertible to type 'int32'")

	printCode(code)
}

func TestGraphSerializationGenerator_UnknownReferenceType(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "EntityA", true, &TargetObjectMetadata{})
	assert.NoError(t, err)
	_, err = generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "var has non-serialized type EntityB")
}

func TestGraphSerializationGenerator_UnknownReferenceTypeIfSkip(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "EntityA", true, &TargetObjectMetadata{})
	err = generator.AddEnabledType("pbgraphserialization", "EntityB", true, &TargetObjectMetadata{skip: true})
	assert.NoError(t, err)
	_, err = generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "var has non-serialized type EntityB")
}

func createTestGenerator(t *testing.T) *Generator {
	generator, err := NewGenerator(
		"github.com/loomnetwork/gamechain/library/pbgraphserialization/",
		"pbgraphserialization",
		"pbgraphserialization_pb_test",
	)
	assert.NoError(t, err)

	return generator
}

func findFuncDecl(codeAst *ast.File, receiverName, funcName string,) *ast.FuncDecl {
	funcDecl := linq.From(codeAst.Decls).WhereT(func(d ast.Decl) bool {
		_, ok := d.(*ast.FuncDecl)
		return ok
	}).WhereT(func(funcDecl *ast.FuncDecl) bool {
		return funcDecl.Name.Name == funcName
	}).WhereT(func(funcDecl *ast.FuncDecl) bool {
		if receiverName != "" {
			starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr)
			if ok {
				ident, ok := starExpr.X.(*ast.Ident)
				if ok {
					return ident.Name == receiverName
				}

				return false
			}
			return false
		} else {
			return true
		}
	}).First()

	if funcDecl == nil {
		return nil
	} else {
		return funcDecl.(*ast.FuncDecl)
	}
}

func parseCode(code string) *ast.File {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test", code, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	return file
}

func printCode(code string) {
	if !debugPrintCode {
		return
	}

	fmt.Println("----- Code start")
	fmt.Println(code)
	fmt.Println("----- Code end")
}
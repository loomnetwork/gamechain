package generator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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

	err := generator.AddEnabledType("pbgraphserialization", "CardAbility", true, &targetObjectMetadata{skip: true})
	assert.NoError(t, err)
	codeSkip, err := generator.Generate()
	assert.NoError(t, err)
	printCode(codeSkip)

	err = generator.AddEnabledType("pbgraphserialization", "CardAbility", true, &targetObjectMetadata{skip: false})
	assert.NoError(t, err)
	codeNoSkip, err := generator.Generate()
	assert.NoError(t, err)
	printCode(codeNoSkip)

	assert.True(t, len(codeNoSkip) > len(codeSkip))
}

func TestGraphSerializationGenerator_GenerateRootMethods(t *testing.T) {
	generator := createTestGenerator(t)

	err := generator.AddEnabledType("pbgraphserialization", "CardAbility", true, &targetObjectMetadata{generateRootSerializationMethods: true})
	assert.NoError(t, err)
	codeRoot, err := generator.Generate()
	assert.NoError(t, err)
	printCode(codeRoot)

	err = generator.AddEnabledType("pbgraphserialization", "CardAbility", true, &targetObjectMetadata{generateRootSerializationMethods: false})
	assert.NoError(t, err)
	codeNoRoot, err := generator.Generate()
	assert.NoError(t, err)
	printCode(codeNoRoot)

	assert.True(t, len(codeRoot) > len(codeNoRoot))
}

func TestGraphSerializationGenerator_AddNonExistentType(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "NonExistent", true, &targetObjectMetadata{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type 'NonExistent' not found in package 'pbgraphserialization'")
}

func TestGraphSerializationGenerator_AddTypeFromNonExistentPackage(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("NonExistent", "", true, &targetObjectMetadata{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "package 'NonExistent' not found")
}

func TestGraphSerializationGenerator_NoMatchingProtoType(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "TypeWithNoMatchingProtoType", true, &targetObjectMetadata{})
	assert.NoError(t, err)
	_, err = generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no matching proto type found for serializable type 'TypeWithNoMatchingProtoType'")
}

func TestGraphSerializationGenerator_FailSerializeSerializationId(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "InvalidType1", true, &targetObjectMetadata{})
	assert.NoError(t, err)
	_, err = generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "serializing 'pbgraphserialization_pb.SerializationId' (field 'card') doesn't makes sense")
}

func TestGraphSerializationGenerator_VariableTypeMismatch(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "InvalidType2", true, &targetObjectMetadata{})
	assert.NoError(t, err)
	_, err = generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variable type mismatch, source has 'types.Named' (Card), destination has 'types.Pointer'")
}

func TestGraphSerializationGenerator_NonPointerSerializationId(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "InvalidType3", true, &targetObjectMetadata{})
	assert.NoError(t, err)
	_, err = generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non-pointer SerializationId is not supported")
}

func TestGraphSerializationGenerator_UnsupportedReferenceType(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "InvalidType4", true, &targetObjectMetadata{})
	assert.NoError(t, err)
	_, err = generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported reference type")
}

func TestGraphSerializationGenerator_MapUnsupportedType(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "InvalidType5", true, &targetObjectMetadata{})
	assert.NoError(t, err)
	_, err = generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type: &types.Map")
}

func TestGraphSerializationGenerator_UnknownReferenceType(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "EntityA", true, &targetObjectMetadata{})
	assert.NoError(t, err)
	_, err = generator.Generate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "var has non-serialized type EntityB")
}

func TestGraphSerializationGenerator_UnknownReferenceTypeIfSkip(t *testing.T) {
	generator := createTestGenerator(t)
	err := generator.AddEnabledType("pbgraphserialization", "EntityA", true, &targetObjectMetadata{})
	err = generator.AddEnabledType("pbgraphserialization", "EntityB", true, &targetObjectMetadata{skip:true})
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

func printCode(code string) {
	if !debugPrintCode {
		return
	}

	fmt.Println("----- Code start")
	fmt.Println(code)
	fmt.Println("----- Code end")
}
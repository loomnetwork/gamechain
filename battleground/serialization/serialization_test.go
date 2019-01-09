package serialization

import (
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/loomnetwork/gamechain/types/serialization"
	"github.com/loomnetwork/gamechain/types/test_serialization"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

const (
	printSerializationJsonOutput = true
)

func TestGraphSerialization_Basic(t *testing.T) {
	sourceCardList := getCardList()

	fullCircleSerializationTest(
		t,
		sourceCardList,
		func(t *testing.T, serializer *Serializer) {

		},
		&serializationpb_test.CardList{},
	)
}

func TestGraphSerialization_CrossReference(t *testing.T) {
	sourceEntityA := &EntityA{
		aField: 3,
	}
	sourceEntityB := &EntityB{
		bField: 4,
	}
	sourceEntityA.entityB = sourceEntityB
	sourceEntityB.entityA = sourceEntityA

	fullCircleSerializationTest(
		t,
		sourceEntityA,
		func(t *testing.T, serializer *Serializer) {
			assert.Equal(t, 1, int(serializer.currentId))
			assert.Equal(t, 2, len(serializer.objectToId))
			assert.True(t, proto.Equal(
				serializer.idToSerializedObject[0],
				&serializationpb_test.EntityA{
					EntityB: Id(1).Marshal(),
					AField:  sourceEntityA.aField,
				},
			))
			assert.True(t, proto.Equal(
				serializer.idToSerializedObject[1],
				&serializationpb_test.EntityB{
					EntityA: Id(0).Marshal(),
					BField:  sourceEntityB.bField,
				},
			))
		},
		&serializationpb_test.EntityA{},
	)
}

func TestGraphSerialization_SelfReference(t *testing.T) {
	sourceEntity := &SelfReferenceEntity{
		field: 3,
	}
	sourceEntity.otherEntity = sourceEntity

	fullCircleSerializationTest(
		t,
		sourceEntity,
		func(t *testing.T, serializer *Serializer) {
			assert.Equal(t, 0, int(serializer.currentId))
			assert.Equal(t, 1, len(serializer.objectToId))
			assert.True(t, proto.Equal(
				serializer.idToSerializedObject[0],
				&serializationpb_test.SelfReferenceEntity{
					OtherEntity: Id(0).Marshal(),
					Field:       sourceEntity.field,
				},
			))
		},
		&serializationpb_test.SelfReferenceEntity{},
	)
}

func TestGraphSerialization_DoubleRoot(t *testing.T) {
	entity := &SelfReferenceEntity{
		field: 3,
	}
	entity.otherEntity = entity

	serializer := NewSerializer()
	assert.NotPanics(t, func() { serializer.Serialize(entity) })
	assert.PanicsWithValue(t, ErrOnlyOneRootObject, func() { serializer.Serialize(entity) })
}

func fullCircleSerializationTest(
	t *testing.T,
	sourceRoot SerializableObject,
	checkSerializerFunc func(t *testing.T, serializer *Serializer),
	serializedProtoMessage proto.Message,
) {
	// serialize
	serializer := NewSerializerSerialize(sourceRoot)
	if checkSerializerFunc != nil {
		checkSerializerFunc(t, serializer)
	}

	if printSerializationJsonOutput {
		fmt.Println("=-- Serialized")
		debugOutputGraphAsJson(serializer)
	}

	// deserialize
	marshaled, err := serializer.Marshal()
	assert.NoError(t, err)

	deserializer, err := NewDeserializerUnmarshal(marshaled)
	assert.NoError(t, err)

	// create an empty instance with the same type as sourceRoot
	deserializedRoot := (reflect.New(reflect.ValueOf(sourceRoot).Elem().Type()).Elem()).Addr().Interface().(SerializableObject)

	deserializedRoot, err =
		deserializer.DeserializeRoot(
			deserializedRoot,
			serializedProtoMessage,
		)
	assert.NoError(t, err)

	// validate
	assert.Equal(t, sourceRoot, deserializedRoot)

	if printSerializationJsonOutput {
		fmt.Println("=-- Deserialized")
		reserializedSerializer := NewSerializerSerialize(deserializedRoot)
		debugOutputGraphAsJson(reserializedSerializer)
	}
}

func convertSerializedGraphToDebugGraph(graph *serializationpb.SerializedGraph) *serializationpb.SerializedDebugGraph {
	debugGraph := serializationpb.SerializedDebugGraph{
		Version: SerializerFormatVersion,
	}

	for i := 0; i < len(graph.Objects); i++ {
		objectData := graph.Objects[i]
		objectTypeName := graph.TypeNames[i]

		debugGraph.Objects = append(debugGraph.Objects, &any.Any{
			Value:   objectData,
			TypeUrl: "type.googleapis.com/" + objectTypeName,
		})
	}

	return &debugGraph
}

func debugOutputGraphAsJson(serializer *Serializer) {
	marshaled, err := serializer.DebugMarshal()
	if err != nil {
		panic(err)
	}

	debugGraph := convertSerializedGraphToDebugGraph(marshaled)

	debugOutputProtoMessageAsJson(debugGraph)
}

func debugOutputProtoMessageAsJson(message proto.Message) {
	m := jsonpb.Marshaler{
		OrigName:     true,
		Indent:       "  ",
		EmitDefaults: true,
	}

	if err := m.Marshal(os.Stdout, message); err != nil {
		fmt.Printf("error generating JSON file: %s", err.Error())
	}

	fmt.Println()
}

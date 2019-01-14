package pbgraphserialization

import (
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/loomnetwork/gamechain/library/pbgraphserialization/proto/pbgraphserialization"
	"github.com/loomnetwork/gamechain/library/pbgraphserialization/proto/test_pbgraphserialization"
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
		&pbgraphserialization_pb_test.CardList{},
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
				&pbgraphserialization_pb_test.EntityA{
					EntityB: SerializationId(1).Serialize(),
					AField:  sourceEntityA.aField,
				},
			))
			assert.True(t, proto.Equal(
				serializer.idToSerializedObject[1],
				&pbgraphserialization_pb_test.EntityB{
					EntityA: SerializationId(0).Serialize(),
					BField:  sourceEntityB.bField,
				},
			))
		},
		&pbgraphserialization_pb_test.EntityA{},
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
				&pbgraphserialization_pb_test.SelfReferenceEntity{
					OtherEntity: SerializationId(0).Serialize(),
					Field:       sourceEntity.field,
				},
			))
		},
		&pbgraphserialization_pb_test.SelfReferenceEntity{},
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
	serializedGraph, err := serializer.SerializeToGraph()
	assert.NoError(t, err)

	deserializer, err := NewDeserializerDeserializeFromGraph(serializedGraph)
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

func convertSerializedGraphToDebugGraph(graph *pbgraphserialization_pb.SerializedGraph) *pbgraphserialization_pb.SerializedDebugGraph {
	debugGraph := pbgraphserialization_pb.SerializedDebugGraph{
		Version: serializerFormatVersion,
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
	serializedGraph, err := serializer.SerializeToDebugGraph()
	if err != nil {
		panic(err)
	}

	debugGraph := convertSerializedGraphToDebugGraph(serializedGraph)

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

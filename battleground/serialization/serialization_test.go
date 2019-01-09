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
	"testing"
)

type CardAbility struct {
	targetType string
	effect     string
}

type Card struct {
	name      string
	abilities []*CardAbility
}

type CardList struct {
	abilities []*CardAbility
	cards     []*Card
}

type EntityA struct {
	entityB *EntityB
	aField  int32
}

type EntityB struct {
	entityA *EntityA
	bField  int32
}

type SelfReferenceEntity struct {
	otherEntity *SelfReferenceEntity
	field       int32
}

func (entity *SelfReferenceEntity) Serialize(serializer *Serializer) proto.Message {
	return &serializationpb_test.SelfReferenceEntity{
		OtherEntity: serializer.Serialize(entity).Marshal(),
		Field:       entity.field,
	}
}

func (entity *SelfReferenceEntity) Deserialize(deserializer *Deserializer, rawMessage proto.Message) (SerializableObject, error) {
	return nil, nil
}

func (entityA *EntityA) Serialize(serializer *Serializer) proto.Message {
	return &serializationpb_test.EntityA{
		EntityB: serializer.Serialize(entityA.entityB).Marshal(),
		AField:  entityA.aField,
	}
}

func (entityB *EntityB) Serialize(serializer *Serializer) proto.Message {
	return &serializationpb_test.EntityB{
		EntityA: serializer.Serialize(entityB.entityA).Marshal(),
		BField:  entityB.bField,
	}
}

func (entityA *EntityA) Deserialize(deserializer *Deserializer, rawMessage proto.Message) (SerializableObject, error) {
	message := rawMessage.(*serializationpb_test.EntityA)
	entityBDeserialized, err := deserializer.Deserialize(
		message.EntityB,
		func() SerializableObject {
			return &EntityB{}
		}, func() proto.Message {
			return &serializationpb_test.EntityB{}
		},
	)

	if err != nil {
		return nil, err
	}

	entityA.entityB = entityBDeserialized.(*EntityB)
	entityA.aField = message.AField
	return entityA, nil
}

func (entityB *EntityB) Deserialize(deserializer *Deserializer, rawMessage proto.Message) (SerializableObject, error) {
	message := rawMessage.(*serializationpb_test.EntityB)
	entityADeserialized, err := deserializer.Deserialize(
		message.EntityA,
		func() SerializableObject {
			return &EntityA{}
		}, func() proto.Message {
			return &serializationpb_test.EntityA{}
		},
	)

	if err != nil {
		return nil, err
	}

	entityB.entityA = entityADeserialized.(*EntityA)
	entityB.bField = message.BField
	return entityB, nil
}

func (cardAbility *CardAbility) Serialize(serializer *Serializer) proto.Message {
	return &serializationpb_test.CardAbility{
		Effect:     cardAbility.effect,
		TargetType: cardAbility.targetType,
	}
}

func (cardAbility *CardAbility) Deserialize(deserializer *Deserializer, rawMessage proto.Message) (SerializableObject, error) {
	return nil, nil
}

func (card *Card) Serialize(serializer *Serializer) proto.Message {
	instance := &serializationpb_test.Card{
		Name: card.name,
	}

	for _, ability := range card.abilities {
		instance.Abilities = append(instance.Abilities, serializer.Serialize(ability).Marshal())
	}

	return instance
}

func (card *Card) Deserialize(deserializer *Deserializer, rawMessage proto.Message) (SerializableObject, error) {
	return nil, nil
}

func (cardList *CardList) Serialize(serializer *Serializer) proto.Message {
	instance := &serializationpb_test.CardList{}

	for _, ability := range cardList.abilities {
		instance.Abilities = append(instance.Abilities, serializer.Serialize(ability).Marshal())
	}

	for _, card := range cardList.cards {
		instance.Cards = append(instance.Cards, serializer.Serialize(card).Marshal())
	}

	return instance
}

func (cardList *CardList) Deserialize(deserializer *Deserializer, rawMessage proto.Message) (SerializableObject, error) {
	return nil, nil
}

func getCardList() *CardList {
	var cardList CardList
	cardList.abilities = []*CardAbility{
		{
			effect:     "HEAL",
			targetType: "WALKER",
		},
		{
			effect:     "RAGE",
			targetType: "HEAVY",
		},
	}

	cardList.cards = []*Card{
		{
			name: "Poizom",
			abilities: []*CardAbility{
				cardList.abilities[0],
			},
		},
	}

	return &cardList
}

func TestGraphSerialization_1(t *testing.T) {
	cardList := getCardList()
	serializer := NewSerializer()
	serializer.Serialize(cardList)

	debugOutputGraphAsJson(serializer)
}

func TestGraphSerialization_CrossReference(t *testing.T) {
	// data
	sourceEntityA := &EntityA{
		aField: 3,
	}
	sourceEntityB := &EntityB{
		bField: 4,
	}
	sourceEntityA.entityB = sourceEntityB
	sourceEntityB.entityA = sourceEntityA

	// serialize
	serializer := NewSerializerSerialize(sourceEntityA)

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

	// deserialize
	marshaled, err := serializer.Marshal()
	assert.NoError(t, err)

	deserializer, err := NewDeserializerUnmarshal(marshaled)
	assert.NoError(t, err)

	deserializedEntityA, err := deserializer.DeserializeRoot(&EntityA{}, &serializationpb_test.EntityA{})
	assert.NoError(t, err)

	// validate
	assert.Equal(t, sourceEntityA, deserializedEntityA)
}

func TestGraphSerialization_SelfReference(t *testing.T) {
	entity := &SelfReferenceEntity{
		field: 3,
	}
	entity.otherEntity = entity

	serializer := NewSerializerSerialize(entity)

	assert.Equal(t, 0, int(serializer.currentId))
	assert.Equal(t, 1, len(serializer.objectToId))
	assert.True(t, proto.Equal(
		serializer.idToSerializedObject[0],
		&serializationpb_test.SelfReferenceEntity{
			OtherEntity: Id(0).Marshal(),
			Field:       3,
		},
	))

	debugOutputGraphAsJson(serializer)
}

func TestGraphSerialization_DoubleRoot(t *testing.T) {
	entity := &SelfReferenceEntity{
		field: 3,
	}
	entity.otherEntity = entity

	serializer := NewSerializer()

	assert.NotPanics(t, func() { serializer.Serialize(entity) })
	assert.PanicsWithValue(t, ErrorOnlyOneRootObject, func() { serializer.Serialize(entity) })
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

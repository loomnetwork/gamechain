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
	message := rawMessage.(*serializationpb_test.SelfReferenceEntity)
	otherEntityDeserialized, err := deserializer.Deserialize(
		message.OtherEntity,
		func() SerializableObject { return &SelfReferenceEntity{} },
		func() proto.Message { return &serializationpb_test.SelfReferenceEntity{} },
	)

	if err != nil {
		return nil, err
	}

	entity.otherEntity = otherEntityDeserialized.(*SelfReferenceEntity)
	entity.field = message.Field
	return entity, nil
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
		func() SerializableObject { return &EntityB{} },
		func() proto.Message { return &serializationpb_test.EntityB{} },
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
		func() SerializableObject { return &EntityA{} },
		func() proto.Message { return &serializationpb_test.EntityA{} },
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
	message := rawMessage.(*serializationpb_test.CardAbility)

	cardAbility.targetType = message.TargetType
	cardAbility.effect = message.Effect
	return cardAbility, nil
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
	message := rawMessage.(*serializationpb_test.Card)

	card.name = message.Name
	for i := 0; i < len(message.Abilities); i++ {
		abilityDeserialized, err := deserializer.Deserialize(
			message.Abilities[i],
			func() SerializableObject { return &CardAbility{} },
			func() proto.Message { return &serializationpb_test.CardAbility{} },
		)

		if err != nil {
			return nil, err
		}

		card.abilities = append(card.abilities, abilityDeserialized.(*CardAbility))
	}

	return card, nil
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
	message := rawMessage.(*serializationpb_test.CardList)

	for i := 0; i < len(message.Abilities); i++ {
		abilityDeserialized, err := deserializer.Deserialize(
			message.Abilities[i],
			func() SerializableObject { return &CardAbility{} },
			func() proto.Message { return &serializationpb_test.CardAbility{} },
		)

		if err != nil {
			return nil, err
		}

		cardList.abilities = append(cardList.abilities, abilityDeserialized.(*CardAbility))
	}

	for i := 0; i < len(message.Cards); i++ {
		cardDeserialized, err := deserializer.Deserialize(
			message.Cards[i],
			func() SerializableObject { return &Card{} },
			func() proto.Message { return &serializationpb_test.Card{} },
		)

		if err != nil {
			return nil, err
		}

		cardList.cards = append(cardList.cards, cardDeserialized.(*Card))
	}

	return cardList, nil
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

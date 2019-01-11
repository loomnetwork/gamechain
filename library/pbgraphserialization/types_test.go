package pbgraphserialization

import (
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/library/pbgraphserialization/internal/proto/pbgraphserialization"
	"github.com/loomnetwork/gamechain/library/pbgraphserialization/internal/proto/test_pbgraphserialization"
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

//pbgraphserialization:enable
type SelfReferenceEntity struct {
	otherEntity *SelfReferenceEntity
	field       int32
}

type InvalidType1 struct {
	card pbgraphserialization_pb.SerializationId
}

type InvalidType2 struct {
	card []Card
}

type InvalidType3 struct {
	card []Card
}

type InvalidType4 struct {
	cardAbility *CardAbility
}

type InvalidType5 struct {
	cardAbilities map[int]*CardAbility
}

type TypeWithNoMatchingProtoType struct {
	cardAbility *CardAbility
}

type AwesomeEnum int32

const (
	AwesomeEnum_Foo = 0
	AwesomeEnum_Bar = 1
)

//some cool comments
//pbgraphserialization:enable
//pbgraphserialization:root
// some more comments
type ComplexType struct {
	intArray         []int32
	double           float64
	byteArray        []byte
	awesomeEnum      AwesomeEnum
	awesomeEnumArray []AwesomeEnum
	otherEntity      *SelfReferenceEntity
	otherEntityArray []*SelfReferenceEntity
	stringArray      []string
}

func (entity *SelfReferenceEntity) Serialize(serializer *Serializer) proto.Message {
	return &pbgraphserialization_pb_test.SelfReferenceEntity{
		OtherEntity: serializer.Serialize(entity).Marshal(),
		Field:       entity.field,
	}
}

func (entity *SelfReferenceEntity) Deserialize(deserializer *Deserializer, rawMessage proto.Message) (SerializableObject, error) {
	message := rawMessage.(*pbgraphserialization_pb_test.SelfReferenceEntity)
	otherEntityDeserialized, err := deserializer.Deserialize(
		message.OtherEntity,
		func() SerializableObject { return &SelfReferenceEntity{} },
		func() proto.Message { return &pbgraphserialization_pb_test.SelfReferenceEntity{} },
	)

	if err != nil {
		return nil, err
	}

	entity.otherEntity = otherEntityDeserialized.(*SelfReferenceEntity)
	entity.field = message.Field
	return entity, nil
}

func (entityA *EntityA) Serialize(serializer *Serializer) proto.Message {
	return &pbgraphserialization_pb_test.EntityA{
		EntityB: serializer.Serialize(entityA.entityB).Marshal(),
		AField:  entityA.aField,
	}
}

func (entityA *EntityA) Deserialize(deserializer *Deserializer, rawMessage proto.Message) (SerializableObject, error) {
	message := rawMessage.(*pbgraphserialization_pb_test.EntityA)
	entityBDeserialized, err := deserializer.Deserialize(
		message.EntityB,
		func() SerializableObject { return &EntityB{} },
		func() proto.Message { return &pbgraphserialization_pb_test.EntityB{} },
	)

	if err != nil {
		return nil, err
	}

	entityA.entityB = entityBDeserialized.(*EntityB)
	entityA.aField = message.AField
	return entityA, nil
}

func (entityB *EntityB) Serialize(serializer *Serializer) proto.Message {
	return &pbgraphserialization_pb_test.EntityB{
		EntityA: serializer.Serialize(entityB.entityA).Marshal(),
		BField:  entityB.bField,
	}
}

func (entityB *EntityB) Deserialize(deserializer *Deserializer, rawMessage proto.Message) (SerializableObject, error) {
	message := rawMessage.(*pbgraphserialization_pb_test.EntityB)
	entityADeserialized, err := deserializer.Deserialize(
		message.EntityA,
		func() SerializableObject { return &EntityA{} },
		func() proto.Message { return &pbgraphserialization_pb_test.EntityA{} },
	)

	if err != nil {
		return nil, err
	}

	entityB.entityA = entityADeserialized.(*EntityA)
	entityB.bField = message.BField
	return entityB, nil
}

func (cardAbility *CardAbility) Serialize(serializer *Serializer) proto.Message {
	return &pbgraphserialization_pb_test.CardAbility{
		Effect:     cardAbility.effect,
		TargetType: cardAbility.targetType,
	}
}

func (cardAbility *CardAbility) Deserialize(deserializer *Deserializer, rawMessage proto.Message) (SerializableObject, error) {
	message := rawMessage.(*pbgraphserialization_pb_test.CardAbility)

	cardAbility.targetType = message.TargetType
	cardAbility.effect = message.Effect
	return cardAbility, nil
}

func (card *Card) Serialize(serializer *Serializer) proto.Message {
	instance := &pbgraphserialization_pb_test.Card{
		Name: card.name,
	}

	for _, ability := range card.abilities {
		instance.Abilities = append(instance.Abilities, serializer.Serialize(ability).Marshal())
	}

	return instance
}

func (card *Card) Deserialize(deserializer *Deserializer, rawMessage proto.Message) (SerializableObject, error) {
	message := rawMessage.(*pbgraphserialization_pb_test.Card)

	card.name = message.Name
	for i := 0; i < len(message.Abilities); i++ {
		abilityDeserialized, err := deserializer.Deserialize(
			message.Abilities[i],
			func() SerializableObject { return &CardAbility{} },
			func() proto.Message { return &pbgraphserialization_pb_test.CardAbility{} },
		)

		if err != nil {
			return nil, err
		}

		card.abilities = append(card.abilities, abilityDeserialized.(*CardAbility))
	}

	return card, nil
}

func (cardList *CardList) Serialize(serializer *Serializer) proto.Message {
	instance := &pbgraphserialization_pb_test.CardList{}

	for _, ability := range cardList.abilities {
		instance.Abilities = append(instance.Abilities, serializer.Serialize(ability).Marshal())
	}

	for _, card := range cardList.cards {
		instance.Cards = append(instance.Cards, serializer.Serialize(card).Marshal())
	}

	return instance
}

func (cardList *CardList) Deserialize(deserializer *Deserializer, rawMessage proto.Message) (SerializableObject, error) {
	message := rawMessage.(*pbgraphserialization_pb_test.CardList)

	for i := 0; i < len(message.Abilities); i++ {
		abilityDeserialized, err := deserializer.Deserialize(
			message.Abilities[i],
			func() SerializableObject { return &CardAbility{} },
			func() proto.Message { return &pbgraphserialization_pb_test.CardAbility{} },
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
			func() proto.Message { return &pbgraphserialization_pb_test.Card{} },
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

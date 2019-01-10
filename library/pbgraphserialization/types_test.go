package pbgraphserialization

import (
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/library/pbgraphserialization/proto/pbgraphserialization_test"
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

func (entityB *EntityB) Serialize(serializer *Serializer) proto.Message {
	return &serializationpb_test.EntityB{
		EntityA: serializer.Serialize(entityB.entityA).Marshal(),
		BField:  entityB.bField,
	}
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


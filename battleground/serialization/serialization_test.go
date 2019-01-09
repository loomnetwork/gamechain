package serialization

import (
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/test_serialization"
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

func (entityA *EntityA) Serialize(graph *Graph) proto.Message {
	return &serializationpb_test.EntityA{
		EntityB: graph.SerializeX(entityA.entityB).Marshal(),
		AField:  entityA.aField,
	}
}

func (entityB *EntityB) Serialize(graph *Graph) proto.Message {
	return &serializationpb_test.EntityB{
		EntityA: graph.SerializeX(entityB.entityA).Marshal(),
		BField:  entityB.bField,
	}
}

func (entityA *EntityA) Deserialize(graph *Graph) {

}

func (entityB *EntityB) Deserialize(graph *Graph) {

}

func (cardAbility *CardAbility) Serialize(graph *Graph) proto.Message {
	return &serializationpb_test.CardAbility{
		Effect:     cardAbility.effect,
		TargetType: cardAbility.targetType,
	}
}

func (cardAbility *CardAbility) Deserialize(graph *Graph) {

}

func (card *Card) Serialize(graph *Graph) proto.Message {
	instance := &serializationpb_test.Card{
		Name: card.name,
	}

	for _, ability := range card.abilities {
		instance.Abilities = append(instance.Abilities, graph.SerializeX(ability).Marshal())
	}

	return instance
}

func (card *Card) Deserialize(graph *Graph) {

}

func (cardList *CardList) Serialize(graph *Graph) proto.Message {
	instance := &serializationpb_test.CardList{}

	for _, ability := range cardList.abilities {
		instance.Abilities = append(instance.Abilities, graph.SerializeX(ability).Marshal())
	}

	for _, card := range cardList.cards {
		instance.Cards = append(instance.Cards, graph.SerializeX(card).Marshal())
	}

	return instance
}

func (cardList *CardList) Deserialize(graph *Graph) {

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
	graph := NewGraph()
	graph.SerializeX(cardList)

	m := jsonpb.Marshaler{
		OrigName:     true,
		Indent:       "  ",
		EmitDefaults: true,
	}

	/*	if err := m.Marshal(os.Stdout, serialized); err != nil {
			fmt.Printf("error generating JSON file: %s", err.Error())
		}
	*/
	if err := m.Marshal(os.Stdout, graph.Serialize()); err != nil {
		fmt.Printf("error generating JSON file: %s", err.Error())
	}

	fmt.Println()
}

func TestGraphSerialization_2(t *testing.T) {
	entityA := &EntityA{
		aField: 3,
	}
	entityB := &EntityB{
		bField: 4,
	}
	entityA.entityB = entityB
	entityB.entityA = entityA
	graph := NewGraph()

	graph.SerializeX(entityA)
	//entityA.Serialize(graph)

	//fmt.Println(graph.objectToSerializedInstance)
	m := jsonpb.Marshaler{
		OrigName:     true,
		Indent:       "  ",
		EmitDefaults: true,
	}

	/*	if err := m.Marshal(os.Stdout, serialized); err != nil {
			fmt.Printf("error generating JSON file: %s", err.Error())
		}
	*/
	if err := m.Marshal(os.Stdout, graph.Serialize()); err != nil {
		fmt.Printf("error generating JSON file: %s", err.Error())
	}

	fmt.Println()

	//fmt.Printf("%+v\n", cardList)
}

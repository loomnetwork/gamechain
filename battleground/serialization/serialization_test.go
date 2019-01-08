package serialization

import (
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
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



func (cardAbility *CardAbility) Serialize(graph *Graph) *SerializedInstance {
	instance := &serializationpb_test.CardAbility{
		Effect:     cardAbility.effect,
		TargetType: cardAbility.targetType,
	}

	return graph.AddReference(cardAbility, instance)
}

func (cardAbility *CardAbility) Deserialize(graph *Graph) {

}

func (card *Card) Serialize(graph *Graph) *SerializedInstance {
	instance := &serializationpb_test.Card{
		Name:            card.name,
	}

	for _, ability := range card.abilities {
		serializedAbility := ability.Serialize(graph)
		instance.Abilities = append(instance.Abilities, serializedAbility.id.Marshal())
	}

	return graph.AddReference(card, instance)
}

func (card *Card) Deserialize(graph *Graph) {

}

func (cardList *CardList) Serialize(graph *Graph) *SerializedInstance {
	instance := &serializationpb_test.CardList{}

	for _, ability := range cardList.abilities {
		protoAbility := ability.Serialize(graph)
		instance.Abilities = append(instance.Abilities, protoAbility.id.Marshal())
	}

	for _, card := range cardList.cards {
		protoCard := card.Serialize(graph)
		instance.Cards = append(instance.Cards, protoCard.id.Marshal())
	}

	return graph.AddReference(cardList, instance)
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

	cardList.Serialize(graph)

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

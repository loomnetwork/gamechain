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

func (entity *SelfReferenceEntity) Serialize(graph *Graph) proto.Message {
	return &serializationpb_test.SelfReferenceEntity{
		OtherEntity: graph.Serialize(entity).Marshal(),
		Field:       entity.field,
	}
}

func (entity *SelfReferenceEntity) Deserialize(graph *Graph) {

}

func (entityA *EntityA) Serialize(graph *Graph) proto.Message {
	return &serializationpb_test.EntityA{
		EntityB: graph.Serialize(entityA.entityB).Marshal(),
		AField:  entityA.aField,
	}
}

func (entityB *EntityB) Serialize(graph *Graph) proto.Message {
	return &serializationpb_test.EntityB{
		EntityA: graph.Serialize(entityB.entityA).Marshal(),
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
		instance.Abilities = append(instance.Abilities, graph.Serialize(ability).Marshal())
	}

	return instance
}

func (card *Card) Deserialize(graph *Graph) {

}

func (cardList *CardList) Serialize(graph *Graph) proto.Message {
	instance := &serializationpb_test.CardList{}

	for _, ability := range cardList.abilities {
		instance.Abilities = append(instance.Abilities, graph.Serialize(ability).Marshal())
	}

	for _, card := range cardList.cards {
		instance.Cards = append(instance.Cards, graph.Serialize(card).Marshal())
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
	graph.Serialize(cardList)

	debugOutputGraphAsJson(graph)
}

func TestGraphSerialization_CrossReference(t *testing.T) {
	entityA := &EntityA{
		aField: 3,
	}
	entityB := &EntityB{
		bField: 4,
	}
	entityA.entityB = entityB
	entityB.entityA = entityA

	graph := NewGraphSerialize(entityA)

	assert.Equal(t, int32(1), graph.currentId.Id)
	assert.Equal(t, 2, len(graph.objectToId))
	assert.True(t, proto.Equal(
		graph.idToSerializedObject[0],
		&serializationpb_test.EntityA{
			EntityB: Id{Id: 1}.Marshal(),
			AField:  3,
		},
	))
	assert.True(t, proto.Equal(
		graph.idToSerializedObject[1],
		&serializationpb_test.EntityB{
			EntityA: Id{Id: 0}.Marshal(),
			BField:  4,
		},
	))

	debugOutputGraphAsJson(graph)
}

func TestGraphSerialization_SelfReference(t *testing.T) {
	entity := &SelfReferenceEntity{
		field: 3,
	}
	entity.otherEntity = entity

	graph := NewGraphSerialize(entity)

	assert.Equal(t, int32(0), graph.currentId.Id)
	assert.Equal(t, 1, len(graph.objectToId))
	assert.True(t, proto.Equal(
		graph.idToSerializedObject[0],
		&serializationpb_test.SelfReferenceEntity{
			OtherEntity: Id{Id: 0}.Marshal(),
			Field:       3,
		},
	))

	debugOutputGraphAsJson(graph)
}

func TestGraphSerialization_DoubleRoot(t *testing.T) {
	entity := &SelfReferenceEntity{
		field: 3,
	}
	entity.otherEntity = entity

	graph := NewGraph()

	assert.NotPanics(t, func(){ graph.Serialize(entity) })
	assert.PanicsWithValue(t, ErrorOnlyOneRootObject, func(){ graph.Serialize(entity) })
}

func convertSerializedGraphToDebugGraph(graph *serializationpb.SerializedGraph) *serializationpb.SerializedDebugGraph {
	debugGraph := serializationpb.SerializedDebugGraph{
		Version: graph.Version,
	}

	for i := 0; i < len(graph.Objects); i++ {
		objectData := graph.Objects[i]
		objectTypeName := graph.TypeNames[i]

		debugGraph.Objects = append(debugGraph.Objects, &any.Any{
			Value: objectData,
			TypeUrl: "type.googleapis.com/" + objectTypeName,
		})
	}

	return &debugGraph
}

func debugOutputGraphAsJson(graph *Graph) {
	m := jsonpb.Marshaler{
		OrigName:     true,
		Indent:       "  ",
		EmitDefaults: true,
	}

	marshaled, _ := graph.DebugMarshal()
	debugGraph := convertSerializedGraphToDebugGraph(marshaled)

	if err := m.Marshal(os.Stdout, debugGraph); err != nil {
		fmt.Printf("error generating JSON file: %s", err.Error())
	}

	fmt.Println()
}

package serialization

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/loomnetwork/gamechain/types/serialization"
	"reflect"
)

type Id struct {
	Id int32
}

type Graph struct {
	currentId                  Id
	objectToId map[interface{}]Id
	idToSerializedObject []proto.Message
}

type SerializableObject interface {
	Serialize(graph *Graph) proto.Message
	Deserialize(graph *Graph)
}

func NewGraph() *Graph {
	var graph Graph
	graph.objectToId = make(map[interface{}]Id)
	graph.idToSerializedObject = make([]proto.Message, 1, 1)

	return &graph
}

func (id Id) Marshal() *serializationpb.SerializationId {
	return &serializationpb.SerializationId{SerializationId: id.Id}
}

func (graph *Graph) Serialize() *serializationpb.SerializedGraph {
	instance := serializationpb.SerializedGraph{}

	for i, maxId := 0, int(graph.currentId.Id); i < maxId; i++ {
		serialized := graph.idToSerializedObject[i]
		serializedAny, err := types.MarshalAny(serialized)

		fmt.Println(err)
		instance.Objects = append(instance.Objects, &serializationpb.SerializedInstance{
			SerializationId: &serializationpb.SerializationId{SerializationId: int32(i)},
			Data: &any.Any{
				TypeUrl: serializedAny.TypeUrl,
				Value:   serializedAny.Value,
			},
		})
	}

	return &instance
}

/*func (graph *Graph) GetSerializedInstance(val SerializableObject) *SerializedInstance {
	serializedInstance, ok := graph.objectToSerializedInstance[val]
	fmt.Printf("GetSerializedInstance, ok: %v, %t\n", reflect.TypeOf(val), ok)
	if ok {
		return serializedInstance
	}

	return val.Serialize(graph)
}*/

/*func (graph *Graph) AddReference(val SerializableObject, protoMessage proto.Message) proto.Message {
	id, ok := graph.objectToId[val]
	fmt.Printf("AddReference, ok: %v, %t\n", reflect.TypeOf(val), ok)
	if ok {
		return graph.idToSerializedObject[id.Id]
	}

	id = graph.currentId

	graph.objectToId[val] = id
	graph.currentId.Id = graph.currentId.Id + 1
	graph.idToSerializedObject[int32(graph.currentId.Id)] = protoMessage

	return id
}*/

func (graph *Graph) SerializeX(val SerializableObject) Id {
	id, alreadyAdded := graph.AddReference(val)
	if alreadyAdded {
		return id
	}

	message := val.Serialize(graph)
	graph.AddSerialized(id, message)
	return id
}

func (graph *Graph) AddSerialized(id Id, protoMessage proto.Message) {
	fmt.Printf("1 - len %d cap %d, id: %d\n", len(graph.idToSerializedObject), cap(graph.idToSerializedObject), id.Id)

	if cap(graph.idToSerializedObject) - 1 < int(id.Id) {
		graph.idToSerializedObject = append(graph.idToSerializedObject, make([]proto.Message, cap(graph.idToSerializedObject), cap(graph.idToSerializedObject))...)
		//graph.idToSerializedObject = graph.idToSerializedObject[:int(id.Id + 1)]
	}

	fmt.Printf("2 - len %d cap %d, id: %d\n", len(graph.idToSerializedObject), cap(graph.idToSerializedObject), id.Id)
	graph.idToSerializedObject[int32(id.Id)] = protoMessage
}

func (graph *Graph) AddReference(val SerializableObject) (id Id, alreadyAdded bool) {
	id, ok := graph.objectToId[val]
	fmt.Printf("AddReference, ok: %v, %t\n", reflect.TypeOf(val), ok)
	if ok {
		return id, true
	}

	id = graph.currentId
	graph.objectToId[val] = id
	graph.currentId.Id = graph.currentId.Id + 1

	return id, false
}
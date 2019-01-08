package serialization

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/loomnetwork/gamechain/types/serialization"
)

type Id struct {
	Id int32
}

type SerializedInstance struct {
	id Id
	protoMessage proto.Message
}

type Graph struct {
	currentId                  Id
	objectToSerializedInstance map[interface{}]*SerializedInstance
}

type SerializableObject interface {
	Serialize(graph *Graph) proto.Message
	Deserialize(graph *Graph)
}

func NewGraph() *Graph {
	var graph Graph
	graph.objectToSerializedInstance = make(map[interface{}]*SerializedInstance)

	return &graph
}

func (id *Id) Marshal() *serializationpb.SerializationId {
	return &serializationpb.SerializationId{SerializationId: id.Id}
}

func (graph *Graph) Serialize() *serializationpb.SerializedGraph {
	instance := serializationpb.SerializedGraph{}

	for _, serializedInstance := range graph.objectToSerializedInstance {
		anyx, _ := types.MarshalAny(serializedInstance.protoMessage)
		instance.Objects = append(instance.Objects, &serializationpb.SerializedInstance{
			SerializationId: &serializationpb.SerializationId{SerializationId: serializedInstance.id.Id},
			Data: &any.Any{
				TypeUrl:anyx.TypeUrl,
				Value:anyx.Value,
			},
		})
	}

	return &instance
}

func (graph *Graph) AddReference(val interface{}, protoMessage proto.Message) *SerializedInstance {
	serializedInstance, ok := graph.objectToSerializedInstance[val]
	if ok {
		return serializedInstance
	}

	serializedInstance = &SerializedInstance{
		id: graph.currentId,
		protoMessage: protoMessage,
	}

	graph.objectToSerializedInstance[val] = serializedInstance
	graph.currentId.Id = graph.currentId.Id + 1

	return serializedInstance
}
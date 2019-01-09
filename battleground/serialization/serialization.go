package serialization

import (
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/serialization"
)

type Id struct {
	Id int32
}

type SerializableObject interface {
	Serialize(graph *Graph) proto.Message
	Deserialize(graph *Graph)
}

func (id Id) Marshal() *serializationpb.SerializationId {
	return &serializationpb.SerializationId{SerializationId: id.Id}
}

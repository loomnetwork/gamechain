package serialization

import (
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/serialization"
	"math"
)

const (
	SerializerFormatVersion = 1
)

var (
	NilSerializationId = Id(math.MaxUint32)
)

type Id uint32

type SerializableObject interface {
	Serialize(serializer *Serializer) proto.Message
	Deserialize(deserializer *Deserializer, rawMessage proto.Message) (SerializableObject, error)
}

type UnmarshaledProtoMessageCreator func() proto.Message

type SerializableObjectCreator func() SerializableObject

func (id Id) Marshal() *serializationpb.SerializationId {
	return &serializationpb.SerializationId{SerializationId: uint32(id)}
}

func Unmarshal(id *serializationpb.SerializationId) Id {
	return Id(id.SerializationId)
}

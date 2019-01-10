package pbgraphserialization

import (
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/library/pbgraphserialization/internal/proto/pbgraphserialization"
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

func (id Id) Marshal() *pbgraphserialization_pb.SerializationId {
	return &pbgraphserialization_pb.SerializationId{SerializationId: uint32(id)}
}

func Unmarshal(id *pbgraphserialization_pb.SerializationId) Id {
	return Id(id.SerializationId)
}

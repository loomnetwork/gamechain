package pbgraphserialization

import (
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/library/pbgraphserialization/internal/proto/pbgraphserialization"
	"math"
)

const (
	serializerFormatVersion = 1
)

var (
	nilSerializationId = SerializationId(math.MaxUint32)
)

type SerializableObject interface {
	Serialize(serializer *Serializer) proto.Message
	Deserialize(deserializer *Deserializer, rawMessage proto.Message) (SerializableObject, error)
}

type ProtoMessageCreator func() proto.Message

type SerializableObjectCreator func() SerializableObject

type SerializationId uint32

func (id SerializationId) Serialize() *pbgraphserialization_pb.SerializationId {
	return &pbgraphserialization_pb.SerializationId{SerializationId: uint32(id)}
}

func DeserializeSerializationId(id *pbgraphserialization_pb.SerializationId) SerializationId {
	return SerializationId(id.SerializationId)
}

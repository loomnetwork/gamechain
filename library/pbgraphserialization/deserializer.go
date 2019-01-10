package pbgraphserialization

import (
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/library/pbgraphserialization/internal/proto/pbgraphserialization"
)

type Deserializer struct {
	idToMarshaledObject [][]byte
	idToObject          map[Id]SerializableObject
}

func NewDeserializer() *Deserializer {
	var deserializer Deserializer
	deserializer.idToObject = make(map[Id]SerializableObject)
	return &deserializer
}

func NewDeserializerUnmarshal(marshaledGraph *pbgraphserialization_pb.SerializedGraph) (*Deserializer, error) {
	deserializer := NewDeserializer()
	err := deserializer.Unmarshal(marshaledGraph)
	if err != nil {
		return nil, err
	}
	return deserializer, nil
}

func (deserializer *Deserializer) Unmarshal(marshaledGraph *pbgraphserialization_pb.SerializedGraph) error {
	count := len(marshaledGraph.Objects)
	deserializer.idToMarshaledObject = make([][]byte, count, count)
	for i := 0; i < count; i++ {
		deserializer.idToMarshaledObject[i] = marshaledGraph.Objects[i]
	}

	return nil
}

func (deserializer *Deserializer) Deserialize(id *pbgraphserialization_pb.SerializationId, targetCreator SerializableObjectCreator, unmarshaledProtoMessageCreator UnmarshaledProtoMessageCreator) (SerializableObject, error) {
	unmarshaledId := Unmarshal(id)

	if unmarshaledId == NilSerializationId {
		return nil, nil
	}

	object, ok := deserializer.idToObject[unmarshaledId]
	if ok {
		return object, nil
	}

	marshaled := deserializer.getMarshaledObject(id)
	unmarshaled := unmarshaledProtoMessageCreator()

	err := proto.Unmarshal(marshaled, unmarshaled)
	if err != nil {
		return nil, err
	}

	object = targetCreator()
	deserializer.idToObject[unmarshaledId] = object
	object, err = object.Deserialize(deserializer, unmarshaled)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (deserializer *Deserializer) DeserializeNoError(id *pbgraphserialization_pb.SerializationId, targetCreator SerializableObjectCreator, unmarshaledProtoMessageCreator UnmarshaledProtoMessageCreator) SerializableObject {
	deserialized, _ := deserializer.Deserialize(id, targetCreator, unmarshaledProtoMessageCreator)
	return deserialized
}

func (deserializer *Deserializer) DeserializeRoot(root SerializableObject, message proto.Message) (SerializableObject, error) {
	id := Id(0)
	err := proto.Unmarshal(deserializer.idToMarshaledObject[id], message)
	if err != nil {
		return nil, err
	}
	deserializer.idToObject[id] = root
	root, err = root.Deserialize(deserializer, message)
	if err != nil {
		return nil, err
	}

	return root, nil
}

func (deserializer *Deserializer) getMarshaledObject(id *pbgraphserialization_pb.SerializationId) []byte {
	return deserializer.idToMarshaledObject[id.SerializationId]
}

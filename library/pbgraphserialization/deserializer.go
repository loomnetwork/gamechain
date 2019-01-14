package pbgraphserialization

import (
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/library/pbgraphserialization/proto/pbgraphserialization"
)

type Deserializer struct {
	idToMarshaledObject [][]byte
	idToObject          map[SerializationId]SerializableObject
}

func NewDeserializer() *Deserializer {
	var deserializer Deserializer
	deserializer.idToObject = make(map[SerializationId]SerializableObject)
	return &deserializer
}

func NewDeserializerDeserializeFromGraph(serializedGraph *pbgraphserialization_pb.SerializedGraph) (*Deserializer, error) {
	deserializer := NewDeserializer()
	err := deserializer.DeserializeFromGraph(serializedGraph)
	if err != nil {
		return nil, err
	}
	return deserializer, nil
}

func (deserializer *Deserializer) DeserializeFromGraph(serializedGraph *pbgraphserialization_pb.SerializedGraph) error {
	count := len(serializedGraph.Objects)
	deserializer.idToMarshaledObject = make([][]byte, count, count)
	for i := 0; i < count; i++ {
		deserializer.idToMarshaledObject[i] = serializedGraph.Objects[i]
	}

	return nil
}

func (deserializer *Deserializer) Deserialize(id *pbgraphserialization_pb.SerializationId, targetCreator SerializableObjectCreator, protoMessageCreator ProtoMessageCreator) (SerializableObject, error) {
	deserializedId := DeserializeSerializationId(id)

	if deserializedId == nilSerializationId {
		return nil, nil
	}

	object, ok := deserializer.idToObject[deserializedId]
	if ok {
		return object, nil
	}

	marshaled := deserializer.getMarshaledObject(id)
	unmarshaled := protoMessageCreator()

	err := proto.Unmarshal(marshaled, unmarshaled)
	if err != nil {
		return nil, err
	}

	object = targetCreator()
	deserializer.idToObject[deserializedId] = object
	object, err = object.Deserialize(deserializer, unmarshaled)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (deserializer *Deserializer) DeserializeRoot(root SerializableObject, message proto.Message) (SerializableObject, error) {
	id := SerializationId(0)
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

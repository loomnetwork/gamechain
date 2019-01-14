package pbgraphserialization

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/library/pbgraphserialization/proto/pbgraphserialization"
	"math"
	"reflect"
)

var (
	ErrOnlyOneRootObject   = errors.New("only one root object is allowed")
	initialSerializationId = SerializationId(math.MaxUint32 - 1)
)

type Serializer struct {
	serializingRoot      bool
	currentId            SerializationId
	objectToId           map[interface{}]SerializationId
	idToSerializedObject []proto.Message
	idToMarshaledObject  [][]byte
}

func NewSerializer() *Serializer {
	var serializer Serializer
	serializer.currentId = initialSerializationId
	serializer.objectToId = make(map[interface{}]SerializationId)
	serializer.idToSerializedObject = make([]proto.Message, 8, 8)

	return &serializer
}

func NewSerializerSerialize(object SerializableObject) *Serializer {
	serializer := NewSerializer()
	serializer.Serialize(object)
	return serializer
}

func (serializer *Serializer) Serialize(object SerializableObject) SerializationId {
	if serializer.currentId != initialSerializationId && !serializer.serializingRoot {
		panic(ErrOnlyOneRootObject)
	}

	serializingRoot := false
	if !serializer.serializingRoot {
		serializer.serializingRoot = true
		serializingRoot = true
	}
	id := serializer.serializeInternal(object)
	if serializingRoot {
		serializer.serializingRoot = false
	}

	return id
}

func (serializer *Serializer) SerializeToGraph() (*pbgraphserialization_pb.SerializedGraph, error) {
	return serializer.serializeToGraphInternal(false)
}

func (serializer *Serializer) SerializeToDebugGraph() (*pbgraphserialization_pb.SerializedGraph, error) {
	return serializer.serializeToGraphInternal(true)
}

func (serializer *Serializer) serializeToGraphInternal(saveTypeInfo bool) (*pbgraphserialization_pb.SerializedGraph, error) {
	instance := pbgraphserialization_pb.SerializedGraph{
		Version: serializerFormatVersion,
	}

	maxId := int(serializer.currentId) + 1
	instance.Objects = make([][]byte, maxId, maxId)
	if saveTypeInfo {
		instance.TypeNames = make([]string, maxId, maxId)
	}
	for i := 0; i < maxId; i++ {
		serialized := serializer.idToSerializedObject[i]

		marshaled, err := proto.Marshal(serialized)
		if err != nil {
			return nil, err
		}

		instance.Objects[i] = marshaled
		if saveTypeInfo {
			marshaledTypeName := proto.MessageName(serialized)
			instance.TypeNames[i] = marshaledTypeName
		}
	}

	return &instance, nil
}

func (serializer *Serializer) serializeInternal(object SerializableObject) SerializationId {
	if reflect.ValueOf(object).IsNil() {
		return nilSerializationId
	}

	id, alreadyAdded := serializer.addReference(object)
	if alreadyAdded {
		return id
	}

	message := object.Serialize(serializer)
	serializer.addSerialized(id, message)
	return id
}

func (serializer *Serializer) addSerialized(id SerializationId, protoMessage proto.Message) {
	messageCap := cap(serializer.idToSerializedObject)
	if messageCap - 1 < int(id) {
		serializer.idToSerializedObject = append(serializer.idToSerializedObject, make([]proto.Message, messageCap, messageCap)...)
	}

	serializer.idToSerializedObject[id] = protoMessage
}

func (serializer *Serializer) addReference(object SerializableObject) (id SerializationId, alreadyAdded bool) {
	id, ok := serializer.objectToId[object]
	if ok {
		return id, true
	}

	if serializer.currentId == initialSerializationId {
		serializer.currentId = 0
		id = 0
	} else {
		id = serializer.currentId + 1
	}

	serializer.objectToId[object] = id
	serializer.currentId = id

	return id, false
}

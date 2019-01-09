package serialization

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/serialization"
	"math"
	"reflect"
)

var (
	ErrorOnlyOneRootObject = errors.New("only one root object is allowed")
	initialSerializationId = Id(math.MaxUint32 - 1)
)

type Serializer struct {
	serializingRoot      bool
	currentId            Id
	objectToId           map[interface{}]Id
	idToSerializedObject []proto.Message
	idToMarshaledObject  [][]byte
}

func NewSerializer() *Serializer {
	var serializer Serializer
	serializer.currentId = initialSerializationId
	serializer.objectToId = make(map[interface{}]Id)
	serializer.idToSerializedObject = make([]proto.Message, 8, 8)

	return &serializer
}

func NewSerializerSerialize(object SerializableObject) *Serializer {
	serializer := NewSerializer()
	serializer.Serialize(object)
	return serializer
}

func (serializer *Serializer) Serialize(object SerializableObject) Id {
	if serializer.currentId != initialSerializationId && !serializer.serializingRoot {
		panic(ErrorOnlyOneRootObject)
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

func (serializer *Serializer) Marshal() (*serializationpb.SerializedGraph, error) {
	return serializer.marshalInternal(false)
}

func (serializer *Serializer) DebugMarshal() (*serializationpb.SerializedGraph, error) {
	return serializer.marshalInternal(true)
}

func (serializer *Serializer) marshalInternal(saveTypeInfo bool) (*serializationpb.SerializedGraph, error) {
	instance := serializationpb.SerializedGraph{
		Version: SerializerFormatVersion,
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

func (serializer *Serializer) serializeInternal(object SerializableObject) Id {
	if reflect.ValueOf(object).IsNil() {
		return NilSerializationId
	}

	id, alreadyAdded := serializer.addReference(object)
	if alreadyAdded {
		return id
	}

	message := object.Serialize(serializer)
	serializer.addSerialized(id, message)
	return id
}

func (serializer *Serializer) addSerialized(id Id, protoMessage proto.Message) {
	messageCap := cap(serializer.idToSerializedObject)
	if messageCap-1 < int(id) {
		serializer.idToSerializedObject = append(serializer.idToSerializedObject, make([]proto.Message, messageCap, messageCap)...)
	}

	serializer.idToSerializedObject[id] = protoMessage
}

func (serializer *Serializer) addReference(object SerializableObject) (id Id, alreadyAdded bool) {
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

package serialization

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/serialization"
)

var (
	ErrorOnlyOneRootObject         = errors.New("only one root object is allowed")
)

type Graph struct {
	Version int32

	serializingRoot      bool
	currentId            Id
	objectToId           map[interface{}]Id
	idToSerializedObject []proto.Message
}

func NewGraph() *Graph {
	var graph Graph
	graph.currentId.Id = -1
	graph.objectToId = make(map[interface{}]Id)
	graph.idToSerializedObject = make([]proto.Message, 16, 16)

	return &graph
}

func NewGraphSerialize(object SerializableObject) *Graph {
	graph := NewGraph()
	graph.Serialize(object)
	return graph
}

func (graph *Graph) Serialize(object SerializableObject) Id {
	if graph.currentId.Id != -1 && !graph.serializingRoot {
		panic(ErrorOnlyOneRootObject)
	}

	serializingRoot := false
	if !graph.serializingRoot {
		graph.serializingRoot = true
		serializingRoot = true
	}
	id := graph.serializeInternal(object)
	if serializingRoot {
		graph.serializingRoot = false
	}

	return id
}

func (graph *Graph) Marshal() (*serializationpb.SerializedGraph, error) {
	return graph.marshalInternal(false)
}

func (graph *Graph) DebugMarshal() (*serializationpb.SerializedGraph, error) {
	return graph.marshalInternal(true)
}

func (graph *Graph) marshalInternal(saveTypeInfo bool) (*serializationpb.SerializedGraph, error) {
	instance := serializationpb.SerializedGraph{
		Version: graph.Version,
	}

	maxId := int(graph.currentId.Id) + 1
	instance.Objects = make([][]byte, maxId, maxId)
	if saveTypeInfo {
		instance.TypeNames = make([]string, maxId, maxId)
	}
	for i := 0; i < maxId; i++ {
		serialized := graph.idToSerializedObject[i]

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

func (graph *Graph) serializeInternal(object SerializableObject) Id {
	id, alreadyAdded := graph.addReference(object)
	if alreadyAdded {
		return id
	}

	message := object.Serialize(graph)
	graph.addSerialized(id, message)
	return id
}

func (graph *Graph) addSerialized(id Id, protoMessage proto.Message) {
	messageCap := cap(graph.idToSerializedObject)
	if messageCap-1 < int(id.Id) {
		graph.idToSerializedObject = append(graph.idToSerializedObject, make([]proto.Message, messageCap, messageCap)...)
	}

	graph.idToSerializedObject[int32(id.Id)] = protoMessage
}

func (graph *Graph) addReference(object SerializableObject) (id Id, alreadyAdded bool) {
	id, ok := graph.objectToId[object]
	if ok {
		return id, true
	}

	id.Id = graph.currentId.Id + 1
	graph.objectToId[object] = id
	graph.currentId = id

	return id, false
}

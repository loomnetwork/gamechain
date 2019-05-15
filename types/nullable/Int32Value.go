package nullable

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/nullable/nullable_pb"
)

type Int32Value struct {
	Value int32
}

func (value *Int32Value) Size() int {
	return proto.Size(value.protoType())
}

func (value Int32Value) Marshal() ([]byte, error) {
	return proto.Marshal(value.protoType())
}

func (value *Int32Value) Unmarshal(data []byte) error {
	protoValue := &nullable_pb.Int32Value{}
	err := proto.Unmarshal(data, protoValue)
	if err != nil {
		return err
	}

	value.Value = protoValue.Value
	return nil
}

func (value Int32Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(value.Value)
}

func (value *Int32Value) UnmarshalJSON(data []byte) error {
	var raw int32
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	*value = Int32Value{Value: raw}
	return nil
}

func (value *Int32Value) protoType() *nullable_pb.Int32Value {
	return &nullable_pb.Int32Value{
		Value: value.Value,
	}
}
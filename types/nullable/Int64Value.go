package nullable

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/nullable/nullable_pb"
)

type Int64Value struct {
	Value int64
}

func (value *Int64Value) Size() int {
	return proto.Size(value.protoType())
}

func (value Int64Value) Marshal() ([]byte, error) {
	return proto.Marshal(value.protoType())
}

func (value *Int64Value) Unmarshal(data []byte) error {
	protoValue := &nullable_pb.Int64Value{}
	err := proto.Unmarshal(data, protoValue)
	if err != nil {
		return err
	}

	value.Value = protoValue.Value
	return nil
}

func (value Int64Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(value.Value)
}

func (value *Int64Value) UnmarshalJSON(data []byte) error {
	var raw int64
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	*value = Int64Value{Value: raw}
	return nil
}

func (value *Int64Value) protoType() *nullable_pb.Int64Value {
	return &nullable_pb.Int64Value{
		Value: value.Value,
	}
}

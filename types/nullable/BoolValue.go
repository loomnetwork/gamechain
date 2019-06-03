package nullable

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/nullable/nullable_pb"
)

type BoolValue struct {
	Value bool
}

func (value *BoolValue) Size() int {
	return proto.Size(value.protoType())
}

func (value BoolValue) Marshal() ([]byte, error) {
	return proto.Marshal(value.protoType())
}

func (value *BoolValue) Unmarshal(data []byte) error {
	protoValue := &nullable_pb.BoolValue{}
	err := proto.Unmarshal(data, protoValue)
	if err != nil {
		return err
	}

	value.Value = protoValue.Value
	return nil
}

func (value BoolValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(value.Value)
}

func (value *BoolValue) UnmarshalJSON(data []byte) error {
	var raw bool
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	*value = BoolValue{Value: raw}
	return nil
}

func (value *BoolValue) protoType() *nullable_pb.BoolValue {
	return &nullable_pb.BoolValue{
		Value: value.Value,
	}
}
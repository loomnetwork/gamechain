package nullable

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/nullable/nullable_pb"
)

type StringValue struct {
	Value string
}

func (value *StringValue) Size() int {
	return proto.Size(value.protoType())
}

func (value StringValue) Marshal() ([]byte, error) {
	return proto.Marshal(value.protoType())
}

func (value *StringValue) Unmarshal(data []byte) error {
	protoValue := &nullable_pb.StringValue{}
	err := proto.Unmarshal(data, protoValue)
	if err != nil {
		return err
	}

	value.Value = protoValue.Value
	return nil
}

func (value StringValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(value.Value)
}

func (value *StringValue) UnmarshalJSON(data []byte) error {
	var raw string
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	*value = StringValue{Value: raw}
	return nil
}

func (value *StringValue) protoType() *nullable_pb.StringValue {
	return &nullable_pb.StringValue{
		Value: value.Value,
	}
}
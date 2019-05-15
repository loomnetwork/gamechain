package battleground_nullable

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
)

type CardTypeEnumValue struct {
	Value zb.CardType_Enum
}

func (value *CardTypeEnumValue) Size() int {
	return proto.Size(value.protoType())
}

func (value CardTypeEnumValue) Marshal() ([]byte, error) {
	return proto.Marshal(value.protoType())
}

func (value *CardTypeEnumValue) Unmarshal(data []byte) error {
	protoValue := &zb.CardTypeEnumValue{}
	err := proto.Unmarshal(data, protoValue)
	if err != nil {
		return err
	}

	value.Value = protoValue.Value
	return nil
}

func (value CardTypeEnumValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(value.Value)
}

func (value *CardTypeEnumValue) UnmarshalJSON(data []byte) error {
	var raw zb.CardType_Enum
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	*value = CardTypeEnumValue{Value: raw}
	return nil
}

func (value *CardTypeEnumValue) protoType() *zb.CardTypeEnumValue {
	return &zb.CardTypeEnumValue{
		Value: value.Value,
	}
}
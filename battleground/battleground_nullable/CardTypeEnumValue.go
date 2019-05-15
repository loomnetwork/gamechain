package battleground_nullable

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
)

type CardTypeEnumValue struct {
	Value zb_enums.CardType_Enum
}

func (value *CardTypeEnumValue) Size() int {
	return proto.Size(value.protoType())
}

func (value CardTypeEnumValue) Marshal() ([]byte, error) {
	return proto.Marshal(value.protoType())
}

func (value *CardTypeEnumValue) Unmarshal(data []byte) error {
	protoValue := &zb_enums.CardTypeEnumValue{}
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
	var raw zb_enums.CardType_Enum
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	*value = CardTypeEnumValue{Value: raw}
	return nil
}

func (value *CardTypeEnumValue) protoType() *zb_enums.CardTypeEnumValue {
	return &zb_enums.CardTypeEnumValue{
		Value: value.Value,
	}
}
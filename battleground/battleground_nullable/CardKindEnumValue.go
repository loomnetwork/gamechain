package battleground_nullable

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb_enums"
)

type CardKindEnumValue struct {
	Value zb_enums.CardKind_Enum
}

func (value *CardKindEnumValue) Size() int {
	return proto.Size(value.protoType())
}

func (value CardKindEnumValue) Marshal() ([]byte, error) {
	return proto.Marshal(value.protoType())
}

func (value *CardKindEnumValue) Unmarshal(data []byte) error {
	protoValue := &zb_enums.CardKindEnumValue{}
	err := proto.Unmarshal(data, protoValue)
	if err != nil {
		return err
	}

	value.Value = protoValue.Value
	return nil
}

func (value CardKindEnumValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(value.Value)
}

func (value *CardKindEnumValue) UnmarshalJSON(data []byte) error {
	var raw zb_enums.CardKind_Enum
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	*value = CardKindEnumValue{Value: raw}
	return nil
}

func (value *CardKindEnumValue) protoType() *zb_enums.CardKindEnumValue {
	return &zb_enums.CardKindEnumValue{
		Value: value.Value,
	}
}
package battleground_nullable

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
)

type FactionEnumValue struct {
	Value zb.Faction_Enum
}

func (value *FactionEnumValue) Size() int {
	return proto.Size(value.protoType())
}

func (value FactionEnumValue) Marshal() ([]byte, error) {
	return proto.Marshal(value.protoType())
}

func (value *FactionEnumValue) Unmarshal(data []byte) error {
	protoValue := &zb.FactionEnumValue{}
	err := proto.Unmarshal(data, protoValue)
	if err != nil {
		return err
	}

	value.Value = protoValue.Value
	return nil
}

func (value FactionEnumValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(value.Value)
}

func (value *FactionEnumValue) UnmarshalJSON(data []byte) error {
	var raw zb.Faction_Enum
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	*value = FactionEnumValue{Value: raw}
	return nil
}

func (value *FactionEnumValue) protoType() *zb.FactionEnumValue {
	return &zb.FactionEnumValue{
		Value: value.Value,
	}
}
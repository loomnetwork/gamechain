package battleground_nullable

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
)

type CreatureRankEnumValue struct {
	Value zb_enums.CreatureRank_Enum
}

func (value *CreatureRankEnumValue) Size() int {
	return proto.Size(value.protoType())
}

func (value CreatureRankEnumValue) Marshal() ([]byte, error) {
	return proto.Marshal(value.protoType())
}

func (value *CreatureRankEnumValue) Unmarshal(data []byte) error {
	protoValue := &zb_enums.CreatureRankEnumValue{}
	err := proto.Unmarshal(data, protoValue)
	if err != nil {
		return err
	}

	value.Value = protoValue.Value
	return nil
}

func (value CreatureRankEnumValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(value.Value)
}

func (value *CreatureRankEnumValue) UnmarshalJSON(data []byte) error {
	var raw zb_enums.CreatureRank_Enum
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	*value = CreatureRankEnumValue{Value: raw}
	return nil
}

func (value *CreatureRankEnumValue) protoType() *zb_enums.CreatureRankEnumValue {
	return &zb_enums.CreatureRankEnumValue{
		Value: value.Value,
	}
}
package battleground_nullable

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
)

type CreatureRankEnumValue struct {
	Value zb.CreatureRank_Enum
}

func (value *CreatureRankEnumValue) Size() int {
	return proto.Size(value.protoType())
}

func (value CreatureRankEnumValue) Marshal() ([]byte, error) {
	return proto.Marshal(value.protoType())
}

func (value *CreatureRankEnumValue) Unmarshal(data []byte) error {
	protoValue := &zb.CreatureRankEnumValue{}
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
	var raw zb.CreatureRank_Enum
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	*value = CreatureRankEnumValue{Value: raw}
	return nil
}

func (value *CreatureRankEnumValue) protoType() *zb.CreatureRankEnumValue {
	return &zb.CreatureRankEnumValue{
		Value: value.Value,
	}
}
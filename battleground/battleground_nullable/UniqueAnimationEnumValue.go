package battleground_nullable

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
)

type UniqueAnimationEnumValue struct {
	Value zb_enums.UniqueAnimation_Enum
}

func (value *UniqueAnimationEnumValue) Size() int {
	return proto.Size(value.protoType())
}

func (value UniqueAnimationEnumValue) Marshal() ([]byte, error) {
	return proto.Marshal(value.protoType())
}

func (value *UniqueAnimationEnumValue) Unmarshal(data []byte) error {
	protoValue := &zb_enums.UniqueAnimationEnumValue{}
	err := proto.Unmarshal(data, protoValue)
	if err != nil {
		return err
	}

	value.Value = protoValue.Value
	return nil
}

func (value UniqueAnimationEnumValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(value.Value)
}

func (value *UniqueAnimationEnumValue) UnmarshalJSON(data []byte) error {
	var raw zb_enums.UniqueAnimation_Enum
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	*value = UniqueAnimationEnumValue{Value: raw}
	return nil
}

func (value *UniqueAnimationEnumValue) protoType() *zb_enums.UniqueAnimationEnumValue {
	return &zb_enums.UniqueAnimationEnumValue{
		Value: value.Value,
	}
}
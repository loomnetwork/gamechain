package battleground_proto

import (
	"bytes"
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb/zb_custombase"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
)

type CardKey struct {
	MouldId int64
	Variant zb_enums.CardVariant_Enum
}

func (value *CardKey) Size() int {
	return proto.Size(value.protoType())
}

func (value CardKey) Marshal() ([]byte, error) {
	return proto.Marshal(value.protoType())
}

func (value *CardKey) Unmarshal(data []byte) error {
	protoValue := &zb_custombase.CardKey{}
	err := proto.Unmarshal(data, protoValue)
	if err != nil {
		return err
	}

	value.MouldId = protoValue.MouldId
	value.Variant = protoValue.Variant
	return nil
}

func (value CardKey) MarshalJSON() ([]byte, error) {
	m := jsonpb.Marshaler{
		OrigName:     false,
		Indent:       "",
		EmitDefaults: true,
	}

	jsonString, err := m.MarshalToString(value.protoType())
	if err != nil {
		return nil, err
	}

	return []byte(jsonString), nil
}

func (value *CardKey) UnmarshalJSON(data []byte) error {
	var raw zb_custombase.CardKey
	err := jsonpb.Unmarshal(bytes.NewReader(data), &raw)
	if err != nil {
		return err
	}
	*value = CardKey{
		MouldId: raw.MouldId,
		Variant: raw.Variant,
	}
	return nil
}

func (value *CardKey) String() string {
	out := fmt.Sprintf("MouldId: %d", value.MouldId)
	if value.Variant != zb_enums.CardVariant_Standard {
		out += fmt.Sprintf(", Variant: %s", zb_enums.CardVariant_Enum_name[int32(value.Variant)])
	}

	return out
}

func (value *CardKey) protoType() *zb_custombase.CardKey {
	return &zb_custombase.CardKey {
		MouldId: value.MouldId,
		Variant: value.Variant,
	}
}

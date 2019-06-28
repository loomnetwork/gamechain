package battleground_utility

import (
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/go-loom/common"
	"github.com/loomnetwork/go-loom/types"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"math/big"
	"os"
)

func ProtoMessageToJsonString(pb proto.Message) (string, error) {
	m := jsonpb.Marshaler{
		OrigName:     false,
		Indent:       "  ",
		EmitDefaults: true,
	}

	json, err := m.MarshalToString(pb)
	if err != nil {
		return "", fmt.Errorf("error marshaling Proto to JSON: %s", err.Error())
	}

	return json, nil
}

func ProtoMessageToJsonStringNoIndent(pb proto.Message) (string, error) {
	m := jsonpb.Marshaler{
		OrigName:     false,
		Indent:       "",
		EmitDefaults: true,
	}

	json, err := m.MarshalToString(pb)
	if err != nil {
		return "", fmt.Errorf("error marshaling Proto to JSON: %s", err.Error())
	}

	return json, nil
}

func PrintProtoMessageAsJson(out io.Writer, pb proto.Message) error {
	m := jsonpb.Marshaler{
		OrigName:     false,
		Indent:       "  ",
		EmitDefaults: true,
	}

	if err := m.Marshal(out, pb); err != nil {
		return fmt.Errorf("error marshaling Proto to JSON: %s", err.Error())
	}

	return nil
}

func PrintProtoMessageAsJsonToStdout(pb proto.Message) error {
	if err := PrintProtoMessageAsJson(os.Stdout, pb); err != nil {
		return err
	}

	return nil
}

func ReadJsonFileToProtoMessage(filename string, message proto.Message) error {
	json, err := ReadFileToString(filename)
	if err != nil {
		return errors.Wrap(err, "error reading "+filename)
	}

	if err := jsonpb.UnmarshalString(json, message); err != nil {
		return errors.Wrap(err, "error parsing JSON file "+filename)
	}

	return nil
}

func ReadJsonStringToProtoMessage(json string, message proto.Message) error {
	if err := jsonpb.UnmarshalString(json, message); err != nil {
		return errors.Wrap(err, "error parsing JSON ")
	}

	return nil
}

func ReadFileToString(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func MarshalBigIntProto(v *big.Int) *types.BigUInt {
	return &types.BigUInt{Value: common.BigUInt{Int: v}}
}

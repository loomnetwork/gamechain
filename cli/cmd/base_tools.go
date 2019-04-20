package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

func readJsonFileToProtobuf(filename string, message proto.Message) error {
	json, err := readFileToString(filename)
	if err != nil {
		return errors.Wrap(err, "error reading " + filename)
	}

	if err := jsonpb.UnmarshalString(json, message); err != nil {
		return errors.Wrap(err, "error parsing JSON file " + filename)
	}

	return nil
}

func readFileToString(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func printProtoMessageAsJSON(out io.Writer, pb proto.Message) error {
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

func printProtoMessageAsJSONToStdout(pb proto.Message) error {
	if err := printProtoMessageAsJSON(os.Stdout, pb); err != nil {
		return err
	}

	return nil
}

package battleground

import (
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"io"
	"os"
)

func protoMessageToJSON(pb proto.Message) (string, error) {
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
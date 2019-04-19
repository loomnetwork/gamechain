package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

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

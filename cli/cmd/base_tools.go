package cmd

import (
	"fmt"
	"os"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

func printProtoMessageAsJSONToStdout(pb proto.Message) error {
	m := jsonpb.Marshaler{
		OrigName:     true,
		Indent:       "  ",
		EmitDefaults: true,
	}

	if err := m.Marshal(os.Stdout, pb); err != nil {
		return fmt.Errorf("error parsing JSON file: %s", err.Error())
	}

	return nil
}

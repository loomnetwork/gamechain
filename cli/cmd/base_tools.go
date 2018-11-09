package cmd

import (
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"os"
)

func printProtoMessageAsJsonToStdout(pb proto.Message) error {
	m := jsonpb.Marshaler{
		OrigName: true,
		Indent: "  ",
	}

	if err := m.Marshal(os.Stdout, pb); err != nil {
		return fmt.Errorf("error parsing JSON file: %s", err.Error())
	}

	return nil
}

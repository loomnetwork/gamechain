package battleground

import (
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
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


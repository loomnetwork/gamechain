package nullable_test

import (
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/nullable"
	"github.com/loomnetwork/gamechain/types/nullable/nullable_test_pb"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestNullableWithValue(t *testing.T) {
	message := &nullable_test_pb.TestMessage{
		Int32Value: &nullable.Int32Value{
			Value: 373,
		},
	}

	json, err := protoMessageToJSON(message)
	assert.Nil(t, err)
	assert.Equal(t, "{\"int32Value\":373}", json)

	messageUnmarshaled := &nullable_test_pb.TestMessage{
	}

	err = readJsonStringToProtobuf(json, messageUnmarshaled)
	assert.Nil(t, err)
	assert.True(t, proto.Equal(message, messageUnmarshaled))
}

func TestNullableNull(t *testing.T) {
	message := &nullable_test_pb.TestMessage{
	}

	json, err := protoMessageToJSON(message)
	assert.Nil(t, err)
	assert.Equal(t, "{\"int32Value\":null}", json)

	messageUnmarshaled := &nullable_test_pb.TestMessage{
	}

	err = readJsonStringToProtobuf(json, messageUnmarshaled)
	assert.Nil(t, err)
	assert.True(t, proto.Equal(message, messageUnmarshaled))
}

func readJsonStringToProtobuf(json string, message proto.Message) error {
	if err := jsonpb.UnmarshalString(json, message); err != nil {
		return errors.Wrap(err, "error parsing JSON")
	}

	return nil
}

func protoMessageToJSON(pb proto.Message) (string, error) {
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

func printProtoMessageAsJSON(out io.Writer, pb proto.Message) error {
	m := jsonpb.Marshaler{
		OrigName:     false,
		Indent:       "",
		EmitDefaults: true,
	}

	if err := m.Marshal(out, pb); err != nil {
		return fmt.Errorf("error marshaling Proto to JSON: %s", err.Error())
	}

	return nil
}
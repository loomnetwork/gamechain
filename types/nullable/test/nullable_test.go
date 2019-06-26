package nullable_test

import (
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/tools/battleground_utility"
	"github.com/loomnetwork/gamechain/types/nullable"
	"github.com/loomnetwork/gamechain/types/nullable/nullable_test_pb"
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestNullableWithValue(t *testing.T) {
	message := &nullable_test_pb.TestMessage{
		Int32Value: &nullable.Int32Value{
			Value: 373,
		},
	}

	json, err := battleground_utility.ProtoMessageToJsonString(message)
	assert.Nil(t, err)
	assert.Equal(t, "{\"int32Value\":373}", json)

	messageUnmarshaled := &nullable_test_pb.TestMessage{
	}

	err = battleground_utility.ReadJsonStringToProtoMessage(json, messageUnmarshaled)
	assert.Nil(t, err)
	assert.True(t, proto.Equal(message, messageUnmarshaled))
}

func TestNullableNull(t *testing.T) {
	message := &nullable_test_pb.TestMessage{
	}

	json, err := battleground_utility.ProtoMessageToJsonString(message)
	assert.Nil(t, err)
	assert.Equal(t, "{\"int32Value\":null}", json)

	messageUnmarshaled := &nullable_test_pb.TestMessage{
	}

	err = battleground_utility.ReadJsonStringToProtoMessage(json, messageUnmarshaled)
	assert.Nil(t, err)
	assert.True(t, proto.Equal(message, messageUnmarshaled))
}
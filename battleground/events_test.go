package battleground

import (
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestUserEventsMarshal(t *testing.T) {
	emitMsg := createUserEventBase("SomeUser")
	emitMsg.Event = &zb_calls.UserEvent_FullCardCollectionSync{
		FullCardCollectionSync: &zb_calls.UserEvent_FullCardCollectionSyncEvent{},
	}

	_, err := proto.Marshal(emitMsg)
	assert.Nil(t, err)
}



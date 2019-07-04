package battleground

import "github.com/loomnetwork/gamechain/types/zb/zb_calls"

func createUserEventTopic(userId string) string {
	return TopicUserEventPrefix + userId
}

func createUserEventBase(userId string) *zb_calls.UserEvent {
	return &zb_calls.UserEvent{
		UserId: userId,
	}
}

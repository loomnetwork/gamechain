package main

import (
	"testing"

	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/stretchr/testify/assert"

	"github.com/loomnetwork/zombie_battleground/battleground"
)

func TestReplayAndValidate(t *testing.T) {
	fakeCtx := setupFakeContext()
	zbContract := &battleground.ZombieBattleground{}
	gprJSON := zb.GameReplay{
		Events: []*zb.PlayerActionEvent{
			&zb.PlayerActionEvent{},
		},
	}
	actionList := gprJSON.Events
	err := replayAndValidate(*fakeCtx, zbContract, actionList)
	assert.Nil(t, err)
}

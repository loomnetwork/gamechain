package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loomnetwork/gamechain/types/zb"
)

func TestReplayAndValidate(t *testing.T) {
	/*
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
	*/
}

func TestComparePlayerTestsFailing(t *testing.T) {
	newPlayerStates := []*zb.PlayerState{
		&zb.PlayerState{
			Id:            "test1",
			Defense:       20,
			CurrentGoo:    0,
			GooVials:      0,
			CurrentAction: zb.PlayerActionType_DrawCard,
		},
	}

	logPlayerStates := []*zb.PlayerState{
		&zb.PlayerState{
			Id:            "test1",
			Defense:       21,
			CurrentGoo:    0,
			GooVials:      0,
			CurrentAction: zb.PlayerActionType_DrawCard,
		},
	}

	err := comparePlayerStates(newPlayerStates, logPlayerStates)
	assert.NotNil(t, err)
}

func TestComparePlayerTests(t *testing.T) {
	newPlayerStates := []*zb.PlayerState{
		&zb.PlayerState{
			Id:            "test1",
			Defense:       20,
			CurrentGoo:    0,
			GooVials:      0,
			CurrentAction: zb.PlayerActionType_DrawCard,
		},
	}

	logPlayerStates := []*zb.PlayerState{
		&zb.PlayerState{
			Id:            "test1",
			Defense:       20,
			CurrentGoo:    0,
			GooVials:      0,
			CurrentAction: zb.PlayerActionType_DrawCard,
		},
	}

	err := comparePlayerStates(newPlayerStates, logPlayerStates)
	assert.Nil(t, err)
}

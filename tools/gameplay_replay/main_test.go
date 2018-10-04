package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loomnetwork/zombie_battleground/types/zb"
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
			Hp:            20,
			Mana:          1,
			CurrentAction: zb.PlayerActionType_DrawCard,
		},
	}

	logPlayerStates := []*zb.PlayerState{
		&zb.PlayerState{
			Id:            "test1",
			Hp:            21, // different
			Mana:          1,
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
			Hp:            20,
			Mana:          1,
			CurrentAction: zb.PlayerActionType_DrawCard,
		},
	}

	logPlayerStates := []*zb.PlayerState{
		&zb.PlayerState{
			Id:            "test1",
			Hp:            20,
			Mana:          1,
			CurrentAction: zb.PlayerActionType_DrawCard,
		},
	}

	err := comparePlayerStates(newPlayerStates, logPlayerStates)
	assert.Nil(t, err)
}

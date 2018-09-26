package battleground

import (
	"testing"

	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/stretchr/testify/assert"
)

var defaultDeck1 = zb.Deck{
	Id:     1,
	HeroId: 1,
	Name:   "Default1",
	Cards: []*zb.CardCollection{
		{CardName: "Pyromaz", Amount: 4},
		{CardName: "Quazi", Amount: 4},
		{CardName: "Burrrnn", Amount: 4},
		{CardName: "Cynderman", Amount: 4},
		{CardName: "Werezomb", Amount: 4},
		{CardName: "Modo", Amount: 4},
		{CardName: "Fire-Maw", Amount: 4},
		{CardName: "Zhampion", Amount: 2},
	},
}

var defaultDeck2 = zb.Deck{
	Id:     2,
	HeroId: 2,
	Name:   "Default2",
	Cards: []*zb.CardCollection{
		{CardName: "Gargantua", Amount: 4},
		{CardName: "Cerberus", Amount: 4},
		{CardName: "Izze", Amount: 4},
		{CardName: "Znowman", Amount: 4},
		{CardName: "Ozmoziz", Amount: 4},
		{CardName: "Jetter", Amount: 4},
		{CardName: "Freezzee", Amount: 4},
		{CardName: "Geyzer", Amount: 2},
	},
}

func TestGameStateFunc(t *testing.T) {
	var uid1 = "id1"
	var uid2 = "id2"
	state := &zb.GameState{
		Id: 1,
		PlayerStates: []*zb.PlayerState{
			&zb.PlayerState{
				Id:   uid1,
				Hp:   10,
				Mana: 0,
				Deck: &defaultDeck1,
			},
			&zb.PlayerState{
				Id:   uid2,
				Hp:   10,
				Mana: 0,
				Deck: &defaultDeck2,
			},
		},
		PlayerActions: []*zb.PlayerAction{
			&zb.PlayerAction{ActionType: zb.PlayerActionType_DrawCardPlayer, PlayerId: uid1},
			&zb.PlayerAction{ActionType: zb.PlayerActionType_EndTurn, PlayerId: uid1},
			&zb.PlayerAction{ActionType: zb.PlayerActionType_DrawCardPlayer, PlayerId: uid2},
			&zb.PlayerAction{ActionType: zb.PlayerActionType_EndTurn, PlayerId: uid2},
			&zb.PlayerAction{
				ActionType: zb.PlayerActionType_CardAttack,
				PlayerId:   uid1,
				Action:     &zb.PlayerAction_CardAttack{},
			},
		},
		CurrentActionIndex: -1, // must start with -1
	}
	gp := &Gameplay{
		State: state,
	}
	err := RunStateMachine(gp)
	assert.Nil(t, err)
	// 5 player actions should be added
	assert.EqualValues(t, 4, gp.State.CurrentActionIndex)
	// add more action
	err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_EndTurn, PlayerId: uid1})
	assert.Nil(t, err)
	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_CardAttack,
		PlayerId:   uid2,
		Action:     &zb.PlayerAction_CardAttack{},
	})
	assert.Nil(t, err)
	// 2 more player actions should be added
	assert.EqualValues(t, 6, gp.State.CurrentActionIndex)
}

func TestInvalidUserTurn(t *testing.T) {
	var uid1 = "id1"
	var uid2 = "id2"
	state := &zb.GameState{
		Id: 2,
		PlayerStates: []*zb.PlayerState{
			&zb.PlayerState{Id: uid1, Deck: &defaultDeck1},
			&zb.PlayerState{Id: uid2, Deck: &defaultDeck1},
		},
		PlayerActions:      []*zb.PlayerAction{},
		CurrentActionIndex: -1, // must start with -1
	}
	gp := &Gameplay{
		State: state,
	}
	// add more action
	err := gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_EndTurn, PlayerId: uid2})
	assert.Equal(t, err, errInvalidPlayer)
	err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_DrawCardPlayer, PlayerId: uid1})
	assert.Nil(t, err)
	err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_EndTurn, PlayerId: uid1})
	assert.Nil(t, err)
	gp.PrintState()
}

func TestInvalidAction(t *testing.T) {}

func TestGameAddAction(t *testing.T) {}

func TestGameResumeAtAction(t *testing.T) {}

package battleground

import (
	"testing"

	"github.com/loomnetwork/zombie_battleground/types/zb"
)

var defaultDeck = zb.Deck{
	Id:     0,
	HeroId: 2,
	Name:   "Default",
	Cards: []*zb.CardCollection{
		{CardName: "Banshee", Amount: 2},
		{CardName: "Breezee", Amount: 2},
		{CardName: "Buffer", Amount: 2},
		{CardName: "Soothsayer", Amount: 2},
		{CardName: "Wheezy", Amount: 2},
		{CardName: "Whiffer", Amount: 2},
		{CardName: "Whizpar", Amount: 1},
		{CardName: "Zhocker", Amount: 1},
		{CardName: "Bouncer", Amount: 1},
		{CardName: "Dragger", Amount: 1},
		{CardName: "Guzt", Amount: 1},
		{CardName: "Pushhh", Amount: 1},
	},
}

func TestGameStateFunc(t *testing.T) {
	var uid1 = "id1"
	var uid2 = "id2"
	state := zb.GameState{
		Id: 1,
		PlayerStates: []*zb.PlayerState{
			&zb.PlayerState{
				Id:   uid1,
				Hp:   10,
				Mana: 0,
				Deck: &defaultDeck,
			},
			&zb.PlayerState{
				Id:   uid2,
				Hp:   10,
				Mana: 0,
				Deck: &defaultDeck,
			},
		},
		PlayerActions: []*zb.PlayerAction{
			&zb.PlayerAction{ActionType: zb.ActionType_DRAW_CARD, PlayerId: uid1},
			&zb.PlayerAction{ActionType: zb.ActionType_END_TURN, PlayerId: uid1},
			&zb.PlayerAction{ActionType: zb.ActionType_DRAW_CARD, PlayerId: uid2},
			&zb.PlayerAction{ActionType: zb.ActionType_END_TURN, PlayerId: uid2},
			&zb.PlayerAction{
				ActionType: zb.ActionType_CARD_ATTACK,
				PlayerId:   uid1,
				Action:     &zb.PlayerAction_CardAttack{},
			},
		},
		CurrentActionIndex: -1, // must start with -1
	}
	gp := InitGameplay(state)
	gp.play()
}

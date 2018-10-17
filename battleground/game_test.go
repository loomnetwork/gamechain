package battleground

import (
	"testing"

	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/stretchr/testify/assert"
)

func TestGameStateFunc(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	var deckList zb.DeckList
	err := ctx.Get(MakeVersionedKey("v1", defaultDeckKey), &deckList)
	assert.Nil(t, err)
	player1 := "player-1"
	player2 := "player-2"
	players := []*zb.PlayerState{
		{Id: player1, Deck: deckList.Decks[0]},
		{Id: player2, Deck: deckList.Decks[0]},
	}
	seed := int64(0)
	gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil)
	assert.Nil(t, err)

	// // add more action
	err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_EndTurn, PlayerId: player1})
	assert.Nil(t, err)
	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_CardPlay,
		PlayerId:   player2,
		Action: &zb.PlayerAction_CardPlay{
			CardPlay: &zb.PlayerActionCardPlay{},
		},
	})
	assert.Nil(t, err)
	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_CardPlay,
		PlayerId:   player2,
		Action: &zb.PlayerAction_CardPlay{
			CardPlay: &zb.PlayerActionCardPlay{},
		},
	})
	assert.Nil(t, err)
	err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_EndTurn, PlayerId: player2})
	assert.Nil(t, err)

	// card attack
	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_CardAttack,
		PlayerId:   player1,
		Action: &zb.PlayerAction_CardAttack{
			CardAttack: &zb.PlayerActionCardAttack{
				Attacker: &zb.CardInstance{
					InstanceId: 1,
				},
				AffectObjectType: zb.AffectObjectType_CARD,
				Target: &zb.Unit{
					InstanceId: 2,
				},
			},
		},
	})
	assert.Nil(t, err)
	// card ability used
	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_CardAbilityUsed,
		PlayerId:   player1,
		Action: &zb.PlayerAction_CardAbilityUsed{
			CardAbilityUsed: &zb.PlayerActionCardAbilityUsed{
				Card: &zb.CardInstance{
					InstanceId: 1,
				},
				Targets: []*zb.Unit{
					&zb.Unit{
						InstanceId:       2,
						AffectObjectType: zb.AffectObjectType_CARD,
					},
				},
			},
		},
	})
	assert.Nil(t, err)
	// overload skill used
	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_OverlordSkillUsed,
		PlayerId:   player1,
		Action: &zb.PlayerAction_OverlordSkillUsed{
			OverlordSkillUsed: &zb.PlayerActionOverlordSkillUsed{
				SkillId:          1,
				AffectObjectType: zb.AffectObjectType_CARD,
				Target: &zb.Unit{
					InstanceId: 2,
				},
			},
		},
	})
	assert.Nil(t, err)

	// rankbuff
	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_RankBuff,
		PlayerId:   player1,
		Action: &zb.PlayerAction_RankBuff{
			RankBuff: &zb.PlayerActionRankBuff{
				Card: &zb.CardInstance{
					InstanceId: 1,
				},
				Targets: []*zb.Unit{
					&zb.Unit{
						InstanceId:       2,
						AffectObjectType: zb.AffectObjectType_CARD,
					},
				},
			},
		},
	})
	assert.Nil(t, err)

	// leave match
	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_LeaveMatch,
		PlayerId:   player1,
		Action: &zb.PlayerAction_LeaveMatch{
			LeaveMatch: &zb.PlayerActionLeaveMatch{},
		},
	})
	assert.Nil(t, err)
	assert.True(t, gp.State.IsEnded)
	assert.Equal(t, gp.State.Winner, player2)

	gp.PrintState()
}

func TestInvalidUserTurn(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	var deckList zb.DeckList
	err := ctx.Get(MakeVersionedKey("v1", defaultDeckKey), &deckList)
	assert.Nil(t, err)
	player1 := "player-1"
	player2 := "player-2"
	players := []*zb.PlayerState{
		{Id: player1, Deck: deckList.Decks[0]},
		{Id: player2, Deck: deckList.Decks[0]},
	}
	seed := int64(0)
	gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil)
	assert.Nil(t, err)
	// add more action
	err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_EndTurn, PlayerId: player2})
	assert.Equal(t, err, errInvalidPlayer)
	err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_DrawCard, PlayerId: player1})
	assert.Nil(t, err)
	err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_EndTurn, PlayerId: player1})
	assert.Nil(t, err)
	gp.PrintState()
}

func TestInitialGameplayWithMulligan(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	var deckList zb.DeckList
	err := ctx.Get(MakeVersionedKey("v1", defaultDeckKey), &deckList)
	assert.Nil(t, err)
	player1 := "player-1"
	player2 := "player-2"
	players := []*zb.PlayerState{
		{Id: player1, Deck: deckList.Decks[0]},
		{Id: player2, Deck: deckList.Decks[0]},
	}
	seed := int64(0)
	gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil)
	assert.Nil(t, err)

	// mulligan keep all the cards
	player1Mulligan := gp.State.PlayerStates[0].CardsInHand
	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_Mulligan,
		PlayerId:   player1,
		Action: &zb.PlayerAction_Mulligan{
			Mulligan: &zb.PlayerActionMulligan{
				MulliganedCards: player1Mulligan,
			},
		},
	})
	assert.Nil(t, err)
	for _, card := range player1Mulligan {
		_, found := containCardInCardList(card, gp.State.PlayerStates[0].CardsInHand)
		assert.True(t, found, "mulliganed card should be player hand")
	}

	// mulligan keep only 2 of the card
	player2Mulligan := gp.State.PlayerStates[1].CardsInHand[:2]
	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_Mulligan,
		PlayerId:   player2,
		Action: &zb.PlayerAction_Mulligan{
			Mulligan: &zb.PlayerActionMulligan{
				MulliganedCards: player2Mulligan,
			},
		},
	})
	assert.Nil(t, err)
	for _, card := range player2Mulligan {
		_, found := containCardInCardList(card, gp.State.PlayerStates[1].CardsInHand)
		assert.True(t, found, "mulliganed card should be player hand")
	}
	gp.PrintState()
}

func TestInitialGameplayWithInvalidMulligan(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	var deckList zb.DeckList
	err := ctx.Get(MakeVersionedKey("v1", defaultDeckKey), &deckList)
	assert.Nil(t, err)
	player1 := "player-1"
	player2 := "player-2"
	players := []*zb.PlayerState{
		{Id: player1, Deck: deckList.Decks[0]},
		{Id: player2, Deck: deckList.Decks[0]},
	}
	seed := int64(0)
	gp, err := NewGamePlay(ctx, 5, "v1", players, seed, nil)
	assert.Nil(t, err)

	// mulligan keep only 2 of the card
	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_Mulligan,
		PlayerId:   player2,
		Action: &zb.PlayerAction_Mulligan{
			Mulligan: &zb.PlayerActionMulligan{
				MulliganedCards: []*zb.CardInstance{
					&zb.CardInstance{
						Prototype: &zb.CardPrototype{Name: "test1"},
					},
					&zb.CardInstance{
						Prototype: &zb.CardPrototype{Name: "test2"},
					},
					&zb.CardInstance{
						Prototype: &zb.CardPrototype{Name: "test3"},
					},
				},
			},
		},
	})
	assert.NotNil(t, err)
	gp.PrintState()
}

func TestPopulateDeckCards(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context
	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)
	getDeckResp1, _ := c.GetDeck(ctx, &zb.GetDeckRequest{
		UserId: "player-1",
		DeckId: 1,
	})
	getDeckResp2, _ := c.GetDeck(ctx, &zb.GetDeckRequest{
		UserId: "player-2",
		DeckId: 1,
	})
	playerStates := []*zb.PlayerState{
		&zb.PlayerState{
			Id:   "player-1",
			Deck: getDeckResp1.Deck,
		},
		&zb.PlayerState{
			Id:   "player-2",
			Deck: getDeckResp2.Deck,
		},
	}
	err := populateDeckCards(ctx, playerStates, "v1")
	assert.Nil(t, err)
	assert.NotNil(t, playerStates[0].CardsInDeck)
	assert.NotNil(t, playerStates[1].CardsInDeck)
	assert.Equal(t, len(playerStates[0].Deck.Cards), len(playerStates[0].CardsInDeck))
	assert.Equal(t, len(playerStates[1].Deck.Cards), len(playerStates[1].CardsInDeck))
}

func TestDrawCard(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	var deckList zb.DeckList
	err := ctx.Get(MakeVersionedKey("v1", defaultDeckKey), &deckList)
	assert.Nil(t, err)
	player1 := "player-1"
	player2 := "player-2"

	// DrawCard cannot be called twice for the same turn
	t.Run("Call DrawCard twice (Invalid)", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deckList.Decks[0]},
			{Id: player2, Deck: deckList.Decks[0]},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil)
		assert.Nil(t, err)
		// add more action
		err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_DrawCard, PlayerId: player1})
		assert.Nil(t, err)
		err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_DrawCard, PlayerId: player1})
		assert.Equal(t, errInvalidAction, err)

	})

	t.Run("DrawCard after Endturn", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deckList.Decks[0]},
			{Id: player2, Deck: deckList.Decks[0]},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 4, "v1", players, seed, nil)
		assert.Nil(t, err)
		// add more action
		err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_DrawCard, PlayerId: player1})
		assert.Nil(t, err)
		err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_EndTurn, PlayerId: player1})
		assert.Nil(t, err)
		err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_DrawCard, PlayerId: player2})
		assert.Nil(t, err)
		err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_EndTurn, PlayerId: player2})
		assert.Nil(t, err)
		err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_DrawCard, PlayerId: player1})
		assert.Nil(t, err)
	})
}
func TestCardAttack(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	var deckList zb.DeckList
	err := ctx.Get(MakeVersionedKey("v1", defaultDeckKey), &deckList)
	assert.Nil(t, err)
	player1 := "player-1"
	player2 := "player-2"

	t.Run("Both cards are damaged and survive", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deckList.Decks[0]},
			{Id: player2, Deck: deckList.Decks[0]},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil)
		assert.Nil(t, err)

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, &zb.CardInstance{
			InstanceId: 1,
			Defense:    3,
			Attack:     2,
		})
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, &zb.CardInstance{
			InstanceId: 2,
			Defense:    5,
			Attack:     1,
		})

		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardAttack,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardAttack{
				CardAttack: &zb.PlayerActionCardAttack{
					Attacker: &zb.CardInstance{
						InstanceId: 1,
					},
					AffectObjectType: zb.AffectObjectType_CARD,
					Target: &zb.Unit{
						InstanceId: 2,
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(2), gp.State.PlayerStates[0].CardsInPlay[0].Defense)
		assert.Equal(t, int32(3), gp.State.PlayerStates[1].CardsInPlay[0].Defense)
	})

	t.Run("Target is killed", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deckList.Decks[0]},
			{Id: player2, Deck: deckList.Decks[0]},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil)
		assert.Nil(t, err)

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, &zb.CardInstance{
			InstanceId: 1,
			Defense:    3,
			Attack:     2,
		})
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, &zb.CardInstance{
			InstanceId: 2,
			Defense:    1,
			Attack:     1,
		})

		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardAttack,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardAttack{
				CardAttack: &zb.PlayerActionCardAttack{
					Attacker: &zb.CardInstance{
						InstanceId: 1,
					},
					AffectObjectType: zb.AffectObjectType_CARD,
					Target: &zb.Unit{
						InstanceId: 2,
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(2), gp.State.PlayerStates[0].CardsInPlay[0].Defense)
		assert.Zero(t, len(gp.State.PlayerStates[1].CardsInPlay))
	})

	t.Run("Attacker and target are killed", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deckList.Decks[0]},
			{Id: player2, Deck: deckList.Decks[0]},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil)
		assert.Nil(t, err)

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, &zb.CardInstance{
			InstanceId: 1,
			Defense:    1,
			Attack:     1,
		})
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, &zb.CardInstance{
			InstanceId: 2,
			Defense:    1,
			Attack:     1,
		})

		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardAttack,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardAttack{
				CardAttack: &zb.PlayerActionCardAttack{
					Attacker: &zb.CardInstance{
						InstanceId: 1,
					},
					AffectObjectType: zb.AffectObjectType_CARD,
					Target: &zb.Unit{
						InstanceId: 2,
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.Zero(t, len(gp.State.PlayerStates[0].CardsInPlay))
		assert.Zero(t, len(gp.State.PlayerStates[1].CardsInPlay))
	})
}

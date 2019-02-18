package battleground

import (
	"os"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/stretchr/testify/assert"
)

var (
	firstPlayerHasFirstTurnCheats = []*zb.DebugCheatsConfiguration{{Enabled: true, ForceFirstTurnUserId: "player-1"}, {Enabled: true}}
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
	gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(gp.State.PlayerStates[0].CardsInHand))
	assert.Equal(t, 0, len(gp.State.PlayerStates[0].CardsInPlay))
	assert.Equal(t, 7, len(gp.State.PlayerStates[0].CardsInDeck))
	assert.Equal(t, 0, len(gp.State.PlayerStates[0].CardsInGraveyard))

	assert.Equal(t, 3, len(gp.State.PlayerStates[1].CardsInHand))
	assert.Equal(t, 0, len(gp.State.PlayerStates[1].CardsInPlay))
	assert.Equal(t, 8, len(gp.State.PlayerStates[1].CardsInDeck))
	assert.Equal(t, 0, len(gp.State.PlayerStates[1].CardsInGraveyard))

	// add more action
	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_CardPlay,
		PlayerId:   player1,
		Action: &zb.PlayerAction_CardPlay{
			CardPlay: &zb.PlayerActionCardPlay{
				Card: &zb.InstanceId{Id: 2},
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 3, len(gp.State.PlayerStates[0].CardsInHand))
	assert.Equal(t, 1, len(gp.State.PlayerStates[0].CardsInPlay))
	assert.Equal(t, 7, len(gp.State.PlayerStates[0].CardsInDeck))
	assert.Equal(t, 0, len(gp.State.PlayerStates[0].CardsInGraveyard))

	assert.Equal(t, 3, len(gp.State.PlayerStates[1].CardsInHand))
	assert.Equal(t, 0, len(gp.State.PlayerStates[1].CardsInPlay))
	assert.Equal(t, 8, len(gp.State.PlayerStates[1].CardsInDeck))
	assert.Equal(t, 0, len(gp.State.PlayerStates[1].CardsInGraveyard))

	err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_EndTurn, PlayerId: player1})
	assert.Nil(t, err)
	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_CardPlay,
		PlayerId:   player2,
		Action: &zb.PlayerAction_CardPlay{
			CardPlay: &zb.PlayerActionCardPlay{
				Card: &zb.InstanceId{Id: 13},
			},
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
				Attacker: &zb.InstanceId{Id: 2},
				Target: &zb.Unit{
					InstanceId: &zb.InstanceId{Id: 13},
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
				Card: &zb.InstanceId{Id: 1},
				Targets: []*zb.Unit{
					&zb.Unit{
						InstanceId: &zb.InstanceId{Id: 2},
					},
				},
			},
		},
	})
	assert.Nil(t, err)
	// overlord skill used
	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_OverlordSkillUsed,
		PlayerId:   player1,
		Action: &zb.PlayerAction_OverlordSkillUsed{
			OverlordSkillUsed: &zb.PlayerActionOverlordSkillUsed{
				SkillId: 1,
				Target: &zb.Unit{
					InstanceId: &zb.InstanceId{Id: 2},
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
				Card: &zb.InstanceId{Id: 1},
				Targets: []*zb.Unit{
					&zb.Unit{
						InstanceId: &zb.InstanceId{Id: 2},
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
	gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
	assert.Nil(t, err)
	// add more action
	err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_EndTurn, PlayerId: player2})
	assert.Equal(t, err, errInvalidPlayer)
	cardID := gp.State.PlayerStates[0].CardsInHand[0].InstanceId
	err = gp.AddAction(&zb.PlayerAction{ActionType: zb.PlayerActionType_CardPlay, PlayerId: player1, Action: &zb.PlayerAction_CardPlay{CardPlay: &zb.PlayerActionCardPlay{Card: cardID}}})
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
	gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
	assert.Nil(t, err)

	// mulligan all the cards
	player1Mulligan := []*zb.CardInstance{}
	for _, mulliganCard := range gp.State.PlayerStates[0].CardsInHand[:3] {
		player1Mulligan = append(player1Mulligan, mulliganCard)
	}

	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_Mulligan,
		PlayerId:   player1,
		Action: &zb.PlayerAction_Mulligan{
			Mulligan: &zb.PlayerActionMulligan{
				MulliganedCards: getInstanceIdsFromCardInstances(player1Mulligan),
			},
		},
	})
	assert.Nil(t, err)
	for _, card := range player1Mulligan {
		_, _, found := findCardInCardListByName(card, gp.State.PlayerStates[0].CardsInHand)
		assert.False(t, found, "mulliganed card should not be in player hand")
	}
	assert.True(t, len(gp.State.PlayerStates[0].CardsInHand) >= 3, "cards in hand should still be >= 3")

	// mulligan 2 of the card
	player2Mulligan := []*zb.CardInstance{}
	for _, mulliganCard := range gp.State.PlayerStates[1].CardsInHand[:2] {
		player2Mulligan = append(player2Mulligan, mulliganCard)
	}

	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_Mulligan,
		PlayerId:   player2,
		Action: &zb.PlayerAction_Mulligan{
			Mulligan: &zb.PlayerActionMulligan{
				MulliganedCards: getInstanceIdsFromCardInstances(player2Mulligan),
			},
		},
	})
	assert.Nil(t, err)
	for _, card := range player2Mulligan {
		_, _, found := findCardInCardListByName(card, gp.State.PlayerStates[1].CardsInHand)
		assert.False(t, found, "mulliganed card should not be in player hand")
	}
	assert.True(t, len(gp.State.PlayerStates[1].CardsInHand) >= 3, "cards in hand should still be >= 3")
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
	gp, err := NewGamePlay(ctx, 5, "v1", players, seed, nil, true, nil)
	assert.Nil(t, err)

	// mulligan keep only 2 of the card
	err = gp.AddAction(&zb.PlayerAction{
		ActionType: zb.PlayerActionType_Mulligan,
		PlayerId:   player2,
		Action: &zb.PlayerAction_Mulligan{
			Mulligan: &zb.PlayerActionMulligan{
				MulliganedCards: []*zb.InstanceId{
					{Id: -1},
					{Id: -2},
					{Id: -3},
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

	cardLibrary, err := getCardLibrary(ctx, "v1")
	assert.Nil(t, err)

	err = populateDeckCards(cardLibrary, playerStates, true)
	assert.Nil(t, err)
	assert.NotNil(t, playerStates[0].CardsInDeck)
	assert.NotNil(t, playerStates[1].CardsInDeck)
	s0 := int64(0) // sum of number of cards
	s1 := int64(0)
	for _, cardCollection := range playerStates[0].Deck.Cards {
		s0 += cardCollection.Amount
	}
	assert.Equal(t, s0, int64(len(playerStates[0].CardsInDeck)))

	for _, cardCollection := range playerStates[1].Deck.Cards {
		s1 += cardCollection.Amount
	}
	assert.Equal(t, s1, int64(len(playerStates[1].CardsInDeck)))
}

func TestCardAttack(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setupInitFromFile(c, pubKeyHexString, &addr, &ctx, t)

	player1 := "player-1"
	player2 := "player-2"

	deck0 := &zb.Deck{
		Id:     0,
		HeroId: 2,
		Name:   "Default",
		Cards: []*zb.DeckCard{
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
			{CardName: "Pushhh", Amount: 1},
		},
	}

	t.Run("Both cards are damaged and survive",
		func(t *testing.T) {
			players := []*zb.PlayerState{
				{Id: player1, Deck: deck0},
				{Id: player2, Deck: deck0},
			}
			seed := int64(0)
			gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
			assert.Nil(t, err)

			gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, &zb.CardInstance{
				InstanceId: &zb.InstanceId{Id: 1},
				Prototype:  &zb.Card{},
				Instance: &zb.CardInstanceSpecificData{
					Defense: 3,
					Attack:  2,
				},
				OwnerIndex: 0,
			})
			gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, &zb.CardInstance{
				InstanceId: &zb.InstanceId{Id: 2},
				Prototype:  &zb.Card{},
				Instance: &zb.CardInstanceSpecificData{
					Defense: 5,
					Attack:  1,
				},
				OwnerIndex: 1,
			})

			err = gp.AddAction(&zb.PlayerAction{
				ActionType: zb.PlayerActionType_CardAttack,
				PlayerId:   player1,
				Action: &zb.PlayerAction_CardAttack{
					CardAttack: &zb.PlayerActionCardAttack{
						Attacker: &zb.InstanceId{Id: 1},
						Target: &zb.Unit{
							InstanceId: &zb.InstanceId{Id: 2},
						},
					},
				},
			})
			assert.Nil(t, err)
			assert.Equal(t, int32(2), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
			assert.Equal(t, int32(3), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Defense)
		})

	t.Run("Target is killed", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 1},
			Prototype:  &zb.Card{},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 3,
				Attack:  2,
			},
			OwnerIndex: 0,
		})
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 2},
			Prototype:  &zb.Card{},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 1,
				Attack:  1,
			},
			OwnerIndex: 1,
		})

		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardAttack,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardAttack{
				CardAttack: &zb.PlayerActionCardAttack{
					Attacker: &zb.InstanceId{Id: 1},
					Target: &zb.Unit{
						InstanceId: &zb.InstanceId{Id: 2},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(2), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Zero(t, len(gp.State.PlayerStates[1].CardsInPlay))
		assert.Equal(t, 1, len(gp.State.PlayerStates[1].CardsInGraveyard))
		assert.Equal(t, int32(2), gp.State.PlayerStates[1].CardsInGraveyard[0].InstanceId.Id)
	})

	t.Run("Attacker and target are killed", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 1},
			Prototype:  &zb.Card{},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 1,
				Attack:  1,
			},
			OwnerIndex: 0,
		})
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 2},
			Prototype:  &zb.Card{},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 1,
				Attack:  1,
			},
			OwnerIndex: 1,
		})

		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardAttack,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardAttack{
				CardAttack: &zb.PlayerActionCardAttack{
					Attacker: &zb.InstanceId{Id: 1},
					Target: &zb.Unit{
						InstanceId: &zb.InstanceId{Id: 2},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.Zero(t, len(gp.State.PlayerStates[0].CardsInPlay))
		assert.Zero(t, len(gp.State.PlayerStates[1].CardsInPlay))
		assert.Equal(t, 1, len(gp.State.PlayerStates[0].CardsInGraveyard))
		assert.Equal(t, 1, len(gp.State.PlayerStates[1].CardsInGraveyard))
		assert.Equal(t, int32(1), gp.State.PlayerStates[0].CardsInGraveyard[0].InstanceId.Id)
		assert.Equal(t, int32(2), gp.State.PlayerStates[1].CardsInGraveyard[0].InstanceId.Id)
	})

	t.Run("Opponent overlord is attacked", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 2},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 3,
				Attack:  2,
			},
			OwnerIndex: 0,
		})
		gp.State.PlayerStates[1].Defense = 3

		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardAttack,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardAttack{
				CardAttack: &zb.PlayerActionCardAttack{
					Attacker: &zb.InstanceId{Id: 2},
					Target: &zb.Unit{
						InstanceId: &zb.InstanceId{Id: 1},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(1), gp.State.PlayerStates[1].Defense)
	})

	t.Run("Opponent overlord is attacked and defeated", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, firstPlayerHasFirstTurnCheats)
		assert.Nil(t, err)

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 2},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 3,
				Attack:  2,
			},
			OwnerIndex: 0,
		})
		gp.State.PlayerStates[1].Defense = 1

		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardAttack,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardAttack{
				CardAttack: &zb.PlayerActionCardAttack{
					Attacker: &zb.InstanceId{Id: 2},
					Target: &zb.Unit{
						InstanceId: &zb.InstanceId{Id: 1},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(-1), gp.State.PlayerStates[1].Defense)
		assert.Equal(t, "player-1", gp.State.Winner)
		assert.True(t, gp.isEnded())
	})
}

func TestCardPlay(t *testing.T) {
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
	t.Run("Normal Card Play", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deckList.Decks[0]},
			{Id: player2, Deck: deckList.Decks[0]},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 4, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)
		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardPlay{
				CardPlay: &zb.PlayerActionCardPlay{
					Card: &zb.InstanceId{Id: 3},
				},
			},
		})
		assert.Nil(t, err)
	})
	t.Run("Card not found in hand", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deckList.Decks[0]},
			{Id: player2, Deck: deckList.Decks[0]},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 4, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)
		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardPlay{
				CardPlay: &zb.PlayerActionCardPlay{
					Card: &zb.InstanceId{Id: -1},
				},
			},
		})

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "not found in hand")
	})
	t.Run("CardPlay from empty hand", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deckList.Decks[0]},
			{Id: player2, Deck: deckList.Decks[0]},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 5, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)
		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardPlay{
				CardPlay: &zb.PlayerActionCardPlay{
					Card: &zb.InstanceId{Id: 2},
				},
			},
		})
		assert.Nil(t, err)
		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardPlay{
				CardPlay: &zb.PlayerActionCardPlay{
					Card: &zb.InstanceId{Id: 3},
				},
			},
		})
		assert.Nil(t, err)
		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardPlay{
				CardPlay: &zb.PlayerActionCardPlay{
					Card: &zb.InstanceId{Id: 4},
				},
			},
		})
		assert.Nil(t, err)
		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardPlay{
				CardPlay: &zb.PlayerActionCardPlay{
					Card: &zb.InstanceId{Id: 5},
				},
			},
		})
		assert.Nil(t, err)
		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardPlay{
				CardPlay: &zb.PlayerActionCardPlay{
					Card: &zb.InstanceId{Id: 6},
				},
			},
		})
		assert.Equal(t, errNoCardsInHand, err)
	})
}

func TestCheats(t *testing.T) {
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
	t.Run("CheatDestroyCardsOnBoard", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deckList.Decks[0]},
			{Id: player2, Deck: deckList.Decks[0]},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 4, "v1", players, seed, nil, true, []*zb.DebugCheatsConfiguration{{Enabled: true}, {Enabled: true}})
		assert.Nil(t, err)
		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardPlay{
				CardPlay: &zb.PlayerActionCardPlay{
					Card: &zb.InstanceId{Id: 3},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(gp.activePlayer().CardsInPlay))
		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CheatDestroyCardsOnBoard,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CheatDestroyCardsOnBoard{
				CheatDestroyCardsOnBoard: &zb.PlayerActionCheatDestroyCardsOnBoard{
					DestroyedCards: []*zb.InstanceId{{Id: 3}},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, 0, len(gp.activePlayer().CardsInPlay))
		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CheatDestroyCardsOnBoard,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CheatDestroyCardsOnBoard{
				CheatDestroyCardsOnBoard: &zb.PlayerActionCheatDestroyCardsOnBoard{
					DestroyedCards: []*zb.InstanceId{{Id: 500}},
				},
			},
		})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "card with instance id 500 not found")
	})
}

/*func TestGameReplyState(t *testing.T) {
c := &ZombieBattleground{}
var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
var addr loom.Address
var ctx contract.Context

setup(c, pubKeyHexString, &addr, &ctx, t)

setupAccount(c, ctx, &zb.UpsertAccountRequest{
	UserId:  "ZombieSlayer_885304049522535028281909288283089888794321162558469122615045277120216755610",
	Version: "v1",
}, t)
setupAccount(c, ctx, &zb.UpsertAccountRequest{
	UserId:  "ZombieSlayer_133859560841827127472479602243651375194640870164584945309062317679873410439",
	Version: "v1",
}, t)

setupGameStateFromFile(c, &ctx)
setupMatchFromFile(c, &ctx)

gameState := getGameStateFromFile(c, &ctx)
for i := 0; i < len(gameState.PlayerActions); i++ {
	_, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
		MatchId:      gameState.Id,
		PlayerAction: gameState.PlayerActions[i],
	})
	if err != nil {
		panic(err)
	}
}

clientGameState := getClientGameStateFromFile(c, &ctx)

if gameState.Id != clientGameState.Id {
	assert.Error(t, nil, "Id are not equal")
}

if gameState.CurrentActionIndex != clientGameState.CurrentActionIndex {
	assert.Error(t, nil, "ActionIndexes are not equal")
}

stateCompare(gameState, clientGameState, t)

/*response, err := c.GetGameState(ctx, &zb.GetGameStateRequest{
	MatchId: gameState.Id,
})
if err != nil {
	panic(err)
}
fmt.Println(response.GameState)*/
//fmt.Println(gameState)

/*fmt.Println(response.GameState.PlayerStates[0].Defense)
	fmt.Println(gameState.PlayerStates[0].Defense)
	fmt.Println(clientGameState.PlayerStates[1].Defense)

}*/

func setupGameStateFromFile(c *ZombieBattleground, ctx *contract.Context) {
	// read from game-state file
	f, err := os.Open("./init_game_state.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var gameStateData zb.GameState

	if err := new(jsonpb.Unmarshaler).Unmarshal(f, &gameStateData); err != nil {
		panic(err)
	}

	if err := saveGameState(*ctx, &gameStateData); err != nil {
		panic(err)
	}

}

func setupMatchFromFile(c *ZombieBattleground, ctx *contract.Context) {
	// read from game-state file
	f, err := os.Open("./match.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var match zb.Match

	if err := new(jsonpb.Unmarshaler).Unmarshal(f, &match); err != nil {
		panic(err)
	}

	if err := saveMatch(*ctx, &match); err != nil {
		panic(err)
	}
}

func getGameStateFromFile(c *ZombieBattleground, ctx *contract.Context) zb.GameState {
	f, err := os.Open("./game_state.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var gameState zb.GameState

	if err := new(jsonpb.Unmarshaler).Unmarshal(f, &gameState); err != nil {
		panic(err)
	}

	return gameState

}

func getClientGameStateFromFile(c *ZombieBattleground, ctx *contract.Context) zb.GameState {
	f, err := os.Open("./client_state.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var gameState zb.GameState

	if err := new(jsonpb.Unmarshaler).Unmarshal(f, &gameState); err != nil {
		panic(err)
	}

	return gameState

}

func stateCompare(serverState zb.GameState, clientState zb.GameState, t *testing.T) bool {
	serverPlayerStates := map[string]*zb.PlayerState{}
	clientPlayerStates := map[string]*zb.PlayerState{}

	for i := 0; i < len(clientState.PlayerStates); i++ {
		clientPlayerStates[clientState.PlayerStates[i].Id] = clientState.PlayerStates[i]
	}

	for i := 0; i < len(serverState.PlayerStates); i++ {
		serverPlayerStates[serverState.PlayerStates[i].Id] = serverState.PlayerStates[i]
	}

	//compare player defense, goocost, deck
	for k, v := range serverPlayerStates {
		if clientPlayerStates[k].Defense != v.Defense {
			t.Errorf("Overlord defenses do not match %d, %d", clientPlayerStates[k].Defense, v.Defense)
			return false
		}
		if clientPlayerStates[k].CurrentGoo != v.CurrentGoo {
			t.Errorf("CurrentGoo do not match %d, %d", clientPlayerStates[k].CurrentGoo, v.CurrentGoo)
			return false
		}
		if compareDecks(clientPlayerStates[k].CardsInDeck, v.CardsInDeck, t) {
			t.Errorf("Cards in deck do not match")
			return false
		}
		if compareDecks(clientPlayerStates[k].CardsInHand, v.CardsInHand, t) {
			t.Errorf("Cards in hand do not match")
			return false
		}
		if compareDecks(clientPlayerStates[k].CardsInPlay, v.CardsInPlay, t) {
			t.Errorf("Cards in play do not match")
			return false
		}
		if compareDecks(clientPlayerStates[k].CardsInGraveyard, v.CardsInGraveyard, t) {
			t.Errorf("Cards in graveyard do not match")
			return false
		}
	}
	return true
}

func compareDecks(d1 []*zb.CardInstance, d2 []*zb.CardInstance, t *testing.T) bool {
	if len(d1) != len(d2) {
		t.Errorf("Number of cards are not equal %d, %d", len(d1), len(d2))
		return false
	}

	for i := 0; i < len(d1); i++ {
		result := compareCards(d1[i], d2[i])
		if !result {
			t.Errorf("Card instances are not equal %v, %v", d1[i], d2[i])
			return false
		}
	}

	return true
}

func compareCards(c1 *zb.CardInstance, c2 *zb.CardInstance) bool {
	if c1.Instance.Defense != c2.Instance.Defense {
		return false
	}

	if c1.Instance.Attack != c2.Instance.Attack {
		return false
	}

	if c1.Instance.GooCost != c2.Instance.GooCost {
		return false
	}

	return true
}

package battleground

import (
	"testing"

	"github.com/gogo/protobuf/proto"
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

	// // add more action
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
	// overload skill used
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

	err = populateDeckCards(ctx, cardLibrary, playerStates, true)
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
			})
			gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, &zb.CardInstance{
				InstanceId: &zb.InstanceId{Id: 2},
				Prototype:  &zb.Card{},
				Instance: &zb.CardInstanceSpecificData{
					Defense: 5,
					Attack:  1,
				},
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
		})
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 2},
			Prototype:  &zb.Card{},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 1,
				Attack:  1,
			},
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
		})
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 2},
			Prototype:  &zb.Card{},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 1,
				Attack:  1,
			},
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

	t.Run("Rage ability works", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb.Card{
			Defense: 5,
			Attack:  2,
			Abilities: []*zb.CardAbility{
				{
					Type:    zb.CardAbilityType_Rage,
					Trigger: zb.CardAbilityTrigger_GotDamage,
					Value:   2,
				},
			},
		}
		instance0 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 1},
			Instance:   newCardInstanceSpecificDataFromCardDetails(card0),
			Prototype:  proto.Clone(card0).(*zb.Card),
			AbilitiesInstances: []*zb.CardAbilityInstance{
				&zb.CardAbilityInstance{
					IsActive: true,
					Trigger:  card0.Abilities[0].Trigger,
					AbilityType: &zb.CardAbilityInstance_Rage{
						Rage: &zb.CardAbilityRage{
							AddedAttack: card0.Abilities[0].Value,
						},
					},
				},
			},
		}

		instance1 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 2},
			Prototype:  &zb.Card{},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 5,
				Attack:  1,
			},
		}

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance0)
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, instance1)

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
		assert.Equal(t, int32(4), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Attack)
		assert.Equal(t, int32(4), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, int32(1), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Defense)
		gp.DebugState()
		assert.Equal(t, int32(4), gp.actionOutcomes[0].GetRage().NewAttack)
		assert.Equal(t, int32(1), gp.actionOutcomes[0].GetRage().InstanceId.Id)
	})

	t.Run("PriorityAttack ability when target is killed", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb.Card{
			Defense: 5,
			Attack:  2,
			Abilities: []*zb.CardAbility{
				{
					Type:    zb.CardAbilityType_PriorityAttack,
					Trigger: zb.CardAbilityTrigger_Attack,
				},
			},
		}
		instance0 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 1},
			Instance:   newCardInstanceSpecificDataFromCardDetails(card0),
			Prototype:  proto.Clone(card0).(*zb.Card),
			AbilitiesInstances: []*zb.CardAbilityInstance{
				&zb.CardAbilityInstance{
					IsActive: true,
					Trigger:  card0.Abilities[0].Trigger,
					AbilityType: &zb.CardAbilityInstance_PriorityAttack{
						PriorityAttack: &zb.CardAbilityPriorityAttack{
							AttackerOldDefense: card0.Defense,
						},
					},
				},
			},
		}
		instance1 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 2},
			Prototype:  &zb.Card{},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 1,
				Attack:  1,
			},
		}

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance0)
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, instance1)

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
		assert.Equal(t, int32(5), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Zero(t, len(gp.State.PlayerStates[1].CardsInPlay))
	})

	t.Run("PriorityAttack ability does not trigger when target is not killed", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb.Card{
			Defense: 5,
			Attack:  2,
			Abilities: []*zb.CardAbility{
				{
					Type: zb.CardAbilityType_PriorityAttack,
				},
			},
		}
		instance0 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 1},
			Instance:   newCardInstanceSpecificDataFromCardDetails(card0),
			Prototype:  proto.Clone(card0).(*zb.Card),
		}
		instance1 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 2},
			Prototype:  &zb.Card{},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 5,
				Attack:  1,
			},
		}

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance0)
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, instance1)

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
		assert.Equal(t, int32(4), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, int32(3), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Defense)
	})

	t.Run("Reanimate ability get activated when attacker death", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb.Card{
			Defense: 3,
			Attack:  2,
			Abilities: []*zb.CardAbility{
				{
					Type:    zb.CardAbilityType_ReanimateUnit,
					Trigger: zb.CardAbilityTrigger_Death,
				},
			},
		}
		instance0 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 1},
			Instance:   newCardInstanceSpecificDataFromCardDetails(card0),
			Prototype:  proto.Clone(card0).(*zb.Card),
			AbilitiesInstances: []*zb.CardAbilityInstance{
				&zb.CardAbilityInstance{
					IsActive: true,
					Trigger:  card0.Abilities[0].Trigger,
					AbilityType: &zb.CardAbilityInstance_Reanimate{
						Reanimate: &zb.CardAbilityReanimate{
							Attack:  card0.Attack,
							Defence: card0.Defense,
						},
					},
				},
			},
		}
		instance1 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 2},
			Prototype:  &zb.Card{},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 5,
				Attack:  4,
			},
		}

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance0)
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, instance1)

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
		assert.Equal(t, int32(2), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Attack)
		assert.Equal(t, int32(3), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, false, gp.State.PlayerStates[0].CardsInPlay[0].AbilitiesInstances[0].IsActive)
	})
}

func TestAdditionalDamageToHeavyInAttack(t *testing.T) {
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

	t.Run("AdditionalDamageToHeavyInAttack ability does trigger when target is heavy", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb.Card{
			Defense: 5,
			Attack:  2,
			Abilities: []*zb.CardAbility{
				{
					Type:    zb.CardAbilityType_AdditionalDamageToHeavyInAttack,
					Trigger: zb.CardAbilityTrigger_Attack,
				},
			},
		}
		instance0 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 1},
			Instance:   newCardInstanceSpecificDataFromCardDetails(card0),
			Prototype:  proto.Clone(card0).(*zb.Card),
			AbilitiesInstances: []*zb.CardAbilityInstance{
				&zb.CardAbilityInstance{
					IsActive: true,
					Trigger:  card0.Abilities[0].Trigger,
					AbilityType: &zb.CardAbilityInstance_AdditionalDamageToHeavyInAttack{
						AdditionalDamageToHeavyInAttack: &zb.CardAbilityAdditionalDamageToHeavyInAttack{
							AddedAttack: 2,
						},
					},
				},
			},
		}
		instance1 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 2},
			Prototype:  &zb.Card{},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 5,
				Attack:  1,
				Type:    zb.CreatureType_Heavy,
			},
		}

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance0)
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, instance1)

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
		assert.Equal(t, int32(1), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Defense)
	})

	t.Run("AdditionalDamageToHeavyInAttack ability does NOT trigger when target is NOT heavy", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb.Card{
			Defense: 5,
			Attack:  2,
			Abilities: []*zb.CardAbility{
				{
					Type:    zb.CardAbilityType_AdditionalDamageToHeavyInAttack,
					Trigger: zb.CardAbilityTrigger_Attack,
				},
			},
		}
		instance0 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 1},
			Instance:   newCardInstanceSpecificDataFromCardDetails(card0),
			Prototype:  proto.Clone(card0).(*zb.Card),
			AbilitiesInstances: []*zb.CardAbilityInstance{
				&zb.CardAbilityInstance{
					IsActive: true,
					Trigger:  card0.Abilities[0].Trigger,
					AbilityType: &zb.CardAbilityInstance_AdditionalDamageToHeavyInAttack{
						AdditionalDamageToHeavyInAttack: &zb.CardAbilityAdditionalDamageToHeavyInAttack{
							AddedAttack: 2,
						},
					},
				},
			},
		}
		instance1 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 2},
			Prototype:  &zb.Card{},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 5,
				Attack:  1,
				Type:    zb.CreatureType_Feral,
			},
		}

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance0)
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, instance1)

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
		assert.Equal(t, int32(3), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Defense)
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

package battleground

import (
	"fmt"
	battleground_proto "github.com/loomnetwork/gamechain/battleground/proto"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	"os"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	loom "github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	assert "github.com/stretchr/testify/require"
)

var (
	firstPlayerHasFirstTurnCheats = []*zb_data.DebugCheatsConfiguration{{Enabled: true, ForceFirstTurnUserId: "player-1"}, {Enabled: true}}
)

func TestGameStateFunc(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	defaultDecks, err := loadDefaultDecks(ctx, "v1")
	assert.Nil(t, err)
	player1 := "player-1"
	player2 := "player-2"
	players := []*zb_data.PlayerState{
		{Id: player1, Deck: defaultDecks.Decks[0]},
		{Id: player2, Deck: defaultDecks.Decks[0]},
	}
	seed := int64(0)
	gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(gp.State.PlayerStates[0].CardsInHand))
	assert.Equal(t, 0, len(gp.State.PlayerStates[0].CardsInPlay))
	assert.Equal(t, 26, len(gp.State.PlayerStates[0].CardsInDeck))
	assert.Equal(t, 0, len(gp.State.PlayerStates[0].CardsInGraveyard))

	assert.Equal(t, 3, len(gp.State.PlayerStates[1].CardsInHand))
	assert.Equal(t, 0, len(gp.State.PlayerStates[1].CardsInPlay))
	assert.Equal(t, 27, len(gp.State.PlayerStates[1].CardsInDeck))
	assert.Equal(t, 0, len(gp.State.PlayerStates[1].CardsInGraveyard))

	// add more action
	err = gp.AddAction(&zb_data.PlayerAction{
		ActionType: zb_enums.PlayerActionType_CardPlay,
		PlayerId:   player1,
		Action: &zb_data.PlayerAction_CardPlay{
			CardPlay: &zb_data.PlayerActionCardPlay{
				Card: &zb_data.InstanceId{Id: 2},
			},
		},
	})
	err = gp.AddAction(&zb_data.PlayerAction{
		ActionType: zb_enums.PlayerActionType_CardPlay,
		PlayerId:   player1,
		Action: &zb_data.PlayerAction_CardPlay{
			CardPlay: &zb_data.PlayerActionCardPlay{
				Card: &zb_data.InstanceId{Id: 3},
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(gp.State.PlayerStates[0].CardsInHand))
	assert.Equal(t, 2, len(gp.State.PlayerStates[0].CardsInPlay))
	assert.Equal(t, 26, len(gp.State.PlayerStates[0].CardsInDeck))
	assert.Equal(t, 0, len(gp.State.PlayerStates[0].CardsInGraveyard))

	assert.Equal(t, 3, len(gp.State.PlayerStates[1].CardsInHand))
	assert.Equal(t, 0, len(gp.State.PlayerStates[1].CardsInPlay))
	assert.Equal(t, 27, len(gp.State.PlayerStates[1].CardsInDeck))
	assert.Equal(t, 0, len(gp.State.PlayerStates[1].CardsInGraveyard))

	err = gp.AddAction(&zb_data.PlayerAction{ActionType: zb_enums.PlayerActionType_EndTurn, PlayerId: player1})
	assert.Nil(t, err)
	err = gp.AddAction(&zb_data.PlayerAction{
		ActionType: zb_enums.PlayerActionType_CardPlay,
		PlayerId:   player2,
		Action: &zb_data.PlayerAction_CardPlay{
			CardPlay: &zb_data.PlayerActionCardPlay{
				Card: &zb_data.InstanceId{Id: 32},
			},
		},
	})
	assert.Nil(t, err)
	err = gp.AddAction(&zb_data.PlayerAction{ActionType: zb_enums.PlayerActionType_EndTurn, PlayerId: player2})
	assert.Nil(t, err)

	// card ability used
	err = gp.AddAction(&zb_data.PlayerAction{
		ActionType: zb_enums.PlayerActionType_CardAbilityUsed,
		PlayerId:   player1,
		Action: &zb_data.PlayerAction_CardAbilityUsed{
			CardAbilityUsed: &zb_data.PlayerActionCardAbilityUsed{
				Card: &zb_data.InstanceId{Id: 3},
				Targets: []*zb_data.Unit{
					&zb_data.Unit{
						InstanceId: &zb_data.InstanceId{Id: 32},
					},
				},
			},
		},
	})
	assert.Nil(t, err)
	// card attack
	err = gp.AddAction(&zb_data.PlayerAction{
		ActionType: zb_enums.PlayerActionType_CardAttack,
		PlayerId:   player1,
		Action: &zb_data.PlayerAction_CardAttack{
			CardAttack: &zb_data.PlayerActionCardAttack{
				Attacker: &zb_data.InstanceId{Id: 2},
				Target: &zb_data.Unit{
					InstanceId: &zb_data.InstanceId{Id: 32},
				},
			},
		},
	})
	assert.Nil(t, err)
	// overlord skill used
	err = gp.AddAction(&zb_data.PlayerAction{
		ActionType: zb_enums.PlayerActionType_OverlordSkillUsed,
		PlayerId:   player1,
		Action: &zb_data.PlayerAction_OverlordSkillUsed{
			OverlordSkillUsed: &zb_data.PlayerActionOverlordSkillUsed{
				SkillId: 1,
				Target: &zb_data.Unit{
					InstanceId: &zb_data.InstanceId{Id: 2},
				},
			},
		},
	})
	assert.Nil(t, err)

	// rankbuff
	err = gp.AddAction(&zb_data.PlayerAction{
		ActionType: zb_enums.PlayerActionType_RankBuff,
		PlayerId:   player1,
		Action: &zb_data.PlayerAction_RankBuff{
			RankBuff: &zb_data.PlayerActionRankBuff{
				Card: &zb_data.InstanceId{Id: 1},
				Targets: []*zb_data.Unit{
					&zb_data.Unit{
						InstanceId: &zb_data.InstanceId{Id: 2},
					},
				},
			},
		},
	})
	assert.Nil(t, err)

	// leave match
	err = gp.AddAction(&zb_data.PlayerAction{
		ActionType: zb_enums.PlayerActionType_LeaveMatch,
		PlayerId:   player1,
		Action: &zb_data.PlayerAction_LeaveMatch{
			LeaveMatch: &zb_data.PlayerActionLeaveMatch{},
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

	defaultDecks, err := loadDefaultDecks(ctx, "v1")
	assert.Nil(t, err)
	player1 := "player-1"
	player2 := "player-2"
	players := []*zb_data.PlayerState{
		{Id: player1, Deck: defaultDecks.Decks[0]},
		{Id: player2, Deck: defaultDecks.Decks[0]},
	}
	seed := int64(0)
	gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
	assert.Nil(t, err)
	// add more action
	err = gp.AddAction(&zb_data.PlayerAction{ActionType: zb_enums.PlayerActionType_EndTurn, PlayerId: player2})
	assert.Equal(t, err, errInvalidPlayer)
	cardID := gp.State.PlayerStates[0].CardsInHand[0].InstanceId
	err = gp.AddAction(&zb_data.PlayerAction{ActionType: zb_enums.PlayerActionType_CardPlay, PlayerId: player1, Action: &zb_data.PlayerAction_CardPlay{CardPlay: &zb_data.PlayerActionCardPlay{Card: cardID}}})
	assert.Nil(t, err)
	err = gp.AddAction(&zb_data.PlayerAction{ActionType: zb_enums.PlayerActionType_EndTurn, PlayerId: player1})
	assert.Nil(t, err)
	gp.PrintState()
}

func TestInitialGameplayWithMulligan(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	defaultDecks, err := loadDefaultDecks(ctx, "v1")
	assert.Nil(t, err)
	player1 := "player-1"
	player2 := "player-2"
	players := []*zb_data.PlayerState{
		{Id: player1, Deck: defaultDecks.Decks[0]},
		{Id: player2, Deck: defaultDecks.Decks[0]},
	}
	seed := int64(0)
	gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
	assert.Nil(t, err)

	// mulligan all the cards
	player1Mulligan := []*zb_data.CardInstance{}
	for _, mulliganCard := range gp.State.PlayerStates[0].CardsInHand[:3] {
		player1Mulligan = append(player1Mulligan, mulliganCard)
	}

	err = gp.AddAction(&zb_data.PlayerAction{
		ActionType: zb_enums.PlayerActionType_Mulligan,
		PlayerId:   player1,
		Action: &zb_data.PlayerAction_Mulligan{
			Mulligan: &zb_data.PlayerActionMulligan{
				MulliganedCards: getInstanceIdsFromCardInstances(player1Mulligan),
			},
		},
	})

	assert.Nil(t, err)
	for _, card := range player1Mulligan {
		_, _, found := findCardInCardListByInstanceId(card.InstanceId, gp.State.PlayerStates[0].CardsInHand)
		assert.False(t, found, "mulliganed card should not be in player hand")
	}
	assert.True(t, len(gp.State.PlayerStates[0].CardsInHand) >= 3, "cards in hand should still be >= 3")

	// mulligan 2 of the card
	player2Mulligan := []*zb_data.CardInstance{}
	for _, mulliganCard := range gp.State.PlayerStates[1].CardsInHand[:2] {
		player2Mulligan = append(player2Mulligan, mulliganCard)
	}

	err = gp.AddAction(&zb_data.PlayerAction{
		ActionType: zb_enums.PlayerActionType_Mulligan,
		PlayerId:   player2,
		Action: &zb_data.PlayerAction_Mulligan{
			Mulligan: &zb_data.PlayerActionMulligan{
				MulliganedCards: getInstanceIdsFromCardInstances(player2Mulligan),
			},
		},
	})
	assert.Nil(t, err)
	for _, card := range player2Mulligan {
		_, _, found := findCardInCardListByInstanceId(card.InstanceId, gp.State.PlayerStates[1].CardsInHand)
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

	defaultDecks, err := loadDefaultDecks(ctx, "v1")
	assert.Nil(t, err)
	player1 := "player-1"
	player2 := "player-2"
	players := []*zb_data.PlayerState{
		{Id: player1, Deck: defaultDecks.Decks[0]},
		{Id: player2, Deck: defaultDecks.Decks[0]},
	}
	seed := int64(0)
	gp, err := NewGamePlay(ctx, 5, "v1", players, seed, nil, true, nil)
	assert.Nil(t, err)

	// mulligan keep only 2 of the card
	err = gp.AddAction(&zb_data.PlayerAction{
		ActionType: zb_enums.PlayerActionType_Mulligan,
		PlayerId:   player2,
		Action: &zb_data.PlayerAction_Mulligan{
			Mulligan: &zb_data.PlayerActionMulligan{
				MulliganedCards: []*zb_data.InstanceId{
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
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)
	getDeckResp1, err := c.GetDeck(ctx, &zb_calls.GetDeckRequest{
		UserId: "player-1",
		DeckId: 1,
		Version: "v1",
	})
	assert.Nil(t, err)
	getDeckResp2, err := c.GetDeck(ctx, &zb_calls.GetDeckRequest{
		UserId: "player-2",
		DeckId: 1,
		Version: "v1",
	})
	assert.Nil(t, err)
	playerStates := []*zb_data.PlayerState{
		&zb_data.PlayerState{
			Id:   "player-1",
			Deck: getDeckResp1.Deck,
		},
		&zb_data.PlayerState{
			Id:   "player-2",
			Deck: getDeckResp2.Deck,
		},
	}

	cardLibrary, err := loadCardLibrary(ctx, "v1")
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

	deck0 := &zb_data.Deck{
		Id:     0,
		OverlordId: 1,
		Name:   "Default",
		Cards: []*zb_data.DeckCard{
			{CardKey: battleground_proto.CardKey{MouldId: 90}, Amount: 2},
			{CardKey: battleground_proto.CardKey{MouldId: 91}, Amount: 2},
			{CardKey: battleground_proto.CardKey{MouldId: 96}, Amount: 2},
			{CardKey: battleground_proto.CardKey{MouldId: 3}, Amount: 2},
			{CardKey: battleground_proto.CardKey{MouldId: 2}, Amount: 2},
			{CardKey: battleground_proto.CardKey{MouldId: 92}, Amount: 2},
			{CardKey: battleground_proto.CardKey{MouldId: 1}, Amount: 1},
			{CardKey: battleground_proto.CardKey{MouldId: 93}, Amount: 1},
			{CardKey: battleground_proto.CardKey{MouldId: 7}, Amount: 1},
			{CardKey: battleground_proto.CardKey{MouldId: 94}, Amount: 1},
			{CardKey: battleground_proto.CardKey{MouldId: 5}, Amount: 1},
		},
	}

	t.Run("Both cards are damaged and survive",
		func(t *testing.T) {
			players := []*zb_data.PlayerState{
				{Id: player1, Deck: deck0},
				{Id: player2, Deck: deck0},
			}
			seed := int64(0)
			gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
			assert.Nil(t, err)

			gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, &zb_data.CardInstance{
				InstanceId: &zb_data.InstanceId{Id: 1},
				Prototype:  &zb_data.Card{},
				Instance: &zb_data.CardInstanceSpecificData{
					Defense: 3,
					Damage:  2,
				},
				OwnerIndex: 0,
			})
			gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, &zb_data.CardInstance{
				InstanceId: &zb_data.InstanceId{Id: 2},
				Prototype:  &zb_data.Card{},
				Instance: &zb_data.CardInstanceSpecificData{
					Defense: 5,
					Damage:  1,
				},
				OwnerIndex: 1,
			})

			err = gp.AddAction(&zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_CardAttack,
				PlayerId:   player1,
				Action: &zb_data.PlayerAction_CardAttack{
					CardAttack: &zb_data.PlayerActionCardAttack{
						Attacker: &zb_data.InstanceId{Id: 1},
						Target: &zb_data.Unit{
							InstanceId: &zb_data.InstanceId{Id: 2},
						},
					},
				},
			})
			assert.Nil(t, err)
			assert.Equal(t, int32(2), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
			assert.Equal(t, int32(3), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Defense)
		})

	t.Run("Target is killed", func(t *testing.T) {
		players := []*zb_data.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 1},
			Prototype:  &zb_data.Card{},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 3,
				Damage:  2,
			},
			OwnerIndex: 0,
		})
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 2},
			Prototype:  &zb_data.Card{},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 1,
				Damage:  1,
			},
			OwnerIndex: 1,
		})

		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardAttack,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardAttack{
				CardAttack: &zb_data.PlayerActionCardAttack{
					Attacker: &zb_data.InstanceId{Id: 1},
					Target: &zb_data.Unit{
						InstanceId: &zb_data.InstanceId{Id: 2},
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
		players := []*zb_data.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 1},
			Prototype:  &zb_data.Card{},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 1,
				Damage:  1,
			},
			OwnerIndex: 0,
		})
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 2},
			Prototype:  &zb_data.Card{},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 1,
				Damage:  1,
			},
			OwnerIndex: 1,
		})

		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardAttack,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardAttack{
				CardAttack: &zb_data.PlayerActionCardAttack{
					Attacker: &zb_data.InstanceId{Id: 1},
					Target: &zb_data.Unit{
						InstanceId: &zb_data.InstanceId{Id: 2},
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
		players := []*zb_data.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 2},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 3,
				Damage:  2,
			},
			OwnerIndex: 0,
		})
		gp.State.PlayerStates[1].Defense = 3

		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardAttack,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardAttack{
				CardAttack: &zb_data.PlayerActionCardAttack{
					Attacker: &zb_data.InstanceId{Id: 2},
					Target: &zb_data.Unit{
						InstanceId: &zb_data.InstanceId{Id: 1},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(1), gp.State.PlayerStates[1].Defense)
	})

	t.Run("Opponent overlord is attacked and defeated", func(t *testing.T) {
		players := []*zb_data.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, firstPlayerHasFirstTurnCheats)
		assert.Nil(t, err)

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 2},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 3,
				Damage:  2,
			},
			OwnerIndex: 0,
		})
		gp.State.PlayerStates[1].Defense = 1

		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardAttack,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardAttack{
				CardAttack: &zb_data.PlayerActionCardAttack{
					Attacker: &zb_data.InstanceId{Id: 2},
					Target: &zb_data.Unit{
						InstanceId: &zb_data.InstanceId{Id: 1},
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

	defaultDecks, err := loadDefaultDecks(ctx, "v1")
	assert.Nil(t, err)
	player1 := "player-1"
	player2 := "player-2"
	t.Run("Normal Card Play", func(t *testing.T) {
		players := []*zb_data.PlayerState{
			{Id: player1, Deck: defaultDecks.Decks[0]},
			{Id: player2, Deck: defaultDecks.Decks[0]},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 4, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)
		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardPlay{
				CardPlay: &zb_data.PlayerActionCardPlay{
					Card: &zb_data.InstanceId{Id: 3},
				},
			},
		})
		assert.Nil(t, err)
	})
	t.Run("Card not found in hand", func(t *testing.T) {
		players := []*zb_data.PlayerState{
			{Id: player1, Deck: defaultDecks.Decks[0]},
			{Id: player2, Deck: defaultDecks.Decks[0]},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 4, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)
		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardPlay{
				CardPlay: &zb_data.PlayerActionCardPlay{
					Card: &zb_data.InstanceId{Id: -1},
				},
			},
		})

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "not found in hand")
	})
	t.Run("CardPlay from empty hand", func(t *testing.T) {
		players := []*zb_data.PlayerState{
			{Id: player1, Deck: defaultDecks.Decks[0]},
			{Id: player2, Deck: defaultDecks.Decks[0]},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 5, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)
		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardPlay{
				CardPlay: &zb_data.PlayerActionCardPlay{
					Card: &zb_data.InstanceId{Id: 2},
				},
			},
		})
		assert.Nil(t, err)
		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardPlay{
				CardPlay: &zb_data.PlayerActionCardPlay{
					Card: &zb_data.InstanceId{Id: 3},
				},
			},
		})
		assert.Nil(t, err)
		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardPlay{
				CardPlay: &zb_data.PlayerActionCardPlay{
					Card: &zb_data.InstanceId{Id: 4},
				},
			},
		})
		assert.Nil(t, err)
		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardPlay{
				CardPlay: &zb_data.PlayerActionCardPlay{
					Card: &zb_data.InstanceId{Id: 5},
				},
			},
		})
		assert.Nil(t, err)
		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardPlay{
				CardPlay: &zb_data.PlayerActionCardPlay{
					Card: &zb_data.InstanceId{Id: 6},
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

	defaultDecks, err := loadDefaultDecks(ctx, "v1")
	assert.Nil(t, err)
	player1 := "player-1"
	player2 := "player-2"
	t.Run("CheatDestroyCardsOnBoard", func(t *testing.T) {
		players := []*zb_data.PlayerState{
			{Id: player1, Deck: defaultDecks.Decks[0]},
			{Id: player2, Deck: defaultDecks.Decks[0]},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 4, "v1", players, seed, nil, true, []*zb_data.DebugCheatsConfiguration{{Enabled: true}, {Enabled: true}})
		assert.Nil(t, err)
		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardPlay{
				CardPlay: &zb_data.PlayerActionCardPlay{
					Card: &zb_data.InstanceId{Id: 3},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(gp.activePlayer().CardsInPlay))
		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CheatDestroyCardsOnBoard,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CheatDestroyCardsOnBoard{
				CheatDestroyCardsOnBoard: &zb_data.PlayerActionCheatDestroyCardsOnBoard{
					DestroyedCards: []*zb_data.InstanceId{{Id: 3}},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, 0, len(gp.activePlayer().CardsInPlay))
		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CheatDestroyCardsOnBoard,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CheatDestroyCardsOnBoard{
				CheatDestroyCardsOnBoard: &zb_data.PlayerActionCheatDestroyCardsOnBoard{
					DestroyedCards: []*zb_data.InstanceId{{Id: 500}},
				},
			},
		})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "card with instance id 500 not found")
	})
}

func TestGameReplayState(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "ZombieSlayer_17765869228194024927116692302141924240301573798750730796590384853844410321577",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "ZombieSlayer_8551218729826748508527518552469681437189485015566002759474151402174095726156",
		Version: "v1",
	}, t)

	setupGameStateFromFile(c, &ctx)
	setupMatchFromFile(c, &ctx)

	gameState := getGameStateFromFile(c, &ctx)
	for i := 0; i < len(gameState.PlayerActions); i++ {
		_, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
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

	result := stateCompare(gameState, clientGameState, t)
	assert.Nil(t, result, "States Comparision Failed")

}

func setupGameStateFromFile(c *ZombieBattleground, ctx *contract.Context) {
	// read from game-state file
	f, err := os.Open("./test_data/init_game_state.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var gameStateData zb_data.GameState

	if err := new(jsonpb.Unmarshaler).Unmarshal(f, &gameStateData); err != nil {
		panic(err)
	}

	if err := saveGameState(*ctx, &gameStateData); err != nil {
		panic(err)
	}

}

func setupMatchFromFile(c *ZombieBattleground, ctx *contract.Context) {
	// read from game-state file
	f, err := os.Open("./test_data/match.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var match zb_data.Match

	if err := new(jsonpb.Unmarshaler).Unmarshal(f, &match); err != nil {
		panic(err)
	}

	if err := saveMatch(*ctx, &match); err != nil {
		panic(err)
	}
}

func getGameStateFromFile(c *ZombieBattleground, ctx *contract.Context) zb_data.GameState {
	f, err := os.Open("./test_data/game_state.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var gameState zb_data.GameState

	if err := new(jsonpb.Unmarshaler).Unmarshal(f, &gameState); err != nil {
		panic(err)
	}

	return gameState

}

func getClientGameStateFromFile(c *ZombieBattleground, ctx *contract.Context) zb_data.GameState {
	f, err := os.Open("./test_data/client_state.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var gameState zb_data.GameState

	if err := new(jsonpb.Unmarshaler).Unmarshal(f, &gameState); err != nil {
		panic(err)
	}

	return gameState

}

func stateCompare(serverState zb_data.GameState, clientState zb_data.GameState, t *testing.T) error {
	serverPlayerStates := map[string]*zb_data.PlayerState{}
	clientPlayerStates := map[string]*zb_data.PlayerState{}

	for i := 0; i < len(clientState.PlayerStates); i++ {
		clientPlayerStates[clientState.PlayerStates[i].Id] = clientState.PlayerStates[i]
	}

	for i := 0; i < len(serverState.PlayerStates); i++ {
		serverPlayerStates[serverState.PlayerStates[i].Id] = serverState.PlayerStates[i]
	}

	//compare player defense, goocost, deck
	for k, v := range serverPlayerStates {
		if clientPlayerStates[k].Defense != v.Defense {
			return fmt.Errorf("Overlord defenses do not match %d, %d", clientPlayerStates[k].Defense, v.Defense)
		}
		if clientPlayerStates[k].CurrentGoo != v.CurrentGoo {
			return fmt.Errorf("CurrentGoo do not match %d, %d", clientPlayerStates[k].CurrentGoo, v.CurrentGoo)
		}
		if clientPlayerStates[k].GooVials != v.GooVials {
			return fmt.Errorf("GooVials do not match %d, %d", clientPlayerStates[k].GooVials, v.GooVials)
		}
		if err := compareDecks(clientPlayerStates[k].CardsInDeck, v.CardsInDeck, t); err != nil {
			return fmt.Errorf("Cards in deck do not match")
		}
		if err := compareDecks(clientPlayerStates[k].CardsInHand, v.CardsInHand, t); err != nil {
			return fmt.Errorf("Cards in hand do not match")
		}
		if err := compareDecks(clientPlayerStates[k].CardsInPlay, v.CardsInPlay, t); err != nil {
			return fmt.Errorf("Cards in play do not match")
		}
		// The client state is wrong, we skip this check for now
		/*if err := compareDecks(clientPlayerStates[k].CardsInGraveyard, v.CardsInGraveyard, t); err != nil {
			return fmt.Errorf("Cards in graveyard do not match")
		}*/
	}
	return nil
}

func compareDecks(d1 []*zb_data.CardInstance, d2 []*zb_data.CardInstance, t *testing.T) error {
	if len(d1) != len(d2) {
		return fmt.Errorf("Number of cards are not equal %d, %d", len(d1), len(d2))
	}

	for i := 0; i < len(d1); i++ {
		if err := compareCards(d1[i], d2[i]); err != nil {
			return err
		}
	}
	return nil
}

func compareCards(c1 *zb_data.CardInstance, c2 *zb_data.CardInstance) error {
	if c1.Instance.Defense != c2.Instance.Defense {
		return fmt.Errorf("defenses are not equal %d, %d\n", c1.Instance.Defense, c2.Instance.Defense)
	}

	if c1.Instance.Damage != c2.Instance.Damage {
		return fmt.Errorf("attacks are not equal %d, %d\n", c1.Instance.Damage, c2.Instance.Damage)
	}

	if c1.Instance.Cost != c2.Instance.Cost {
		return fmt.Errorf("goocost are not equal %d, %d\n", c1.Instance.Cost, c2.Instance.Cost)
	}

	return nil
}

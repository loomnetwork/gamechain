package battleground

import (
	battleground_proto "github.com/loomnetwork/gamechain/battleground/proto"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	assert "github.com/stretchr/testify/require"
)

func TestAbilityReplaceUnitsWithTypeOnStrongerOnes(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setupInitFromFile(c, pubKeyHexString, &addr, &ctx, t)

	player1 := "player-1"
	player2 := "player-2"

	deck0 := &zb_data.Deck{
		Id:         0,
		OverlordId: 1,
		Name:       "Default",
		Cards: []*zb_data.DeckCard{
			{CardKey: battleground_proto.CardKey{MouldId: 133}, Amount: 5},
			{CardKey: battleground_proto.CardKey{MouldId: 136}, Amount: 10},
		},
	}

	deck1 := &zb_data.Deck{
		Id:         0,
		OverlordId: 1,
		Name:       "Default",
		Cards: []*zb_data.DeckCard{
			{CardKey: battleground_proto.CardKey{MouldId: 11}, Amount: 15},
		},
	}

	t.Run("Play vortex should replace all water zombie with stronger ones", func(t *testing.T) {
		players := []*zb_data.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck1},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb_data.Card{
			Name:    "Vortex",
			Defense: 3,
			Damage:  2,
			Cost: 1,
			Faction: zb_enums.Faction_Water,
			Abilities: []*zb_data.AbilityData{
				{
					Ability: zb_enums.AbilityType_ReplaceUnitsWithTypeOnStrongerOnes,
					Trigger: zb_enums.AbilityTrigger_Entry,
				},
			},
		}
		instance0 := &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 2},
			Instance:   newCardInstanceSpecificDataFromCardDetails(card0),
			Prototype:  proto.Clone(card0).(*zb_data.Card),
			AbilitiesInstances: []*zb_data.CardAbilityInstance{
				&zb_data.CardAbilityInstance{
					IsActive: true,
					Trigger:  card0.Abilities[0].Trigger,
					AbilityType: &zb_data.CardAbilityInstance_ReplaceUnitsWithTypeOnStrongerOnes{
						ReplaceUnitsWithTypeOnStrongerOnes: &zb_data.CardAbilityReplaceUnitsWithTypeOnStrongerOnes{
							Faction: card0.Faction,
						},
					},
				},
			},
			Zone:  zb_enums.Zone_HAND,
			Owner: player1,
		}
		instance1 := &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 3},
			Prototype:  &zb_data.Card{},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 5,
				Damage:  4,
				Cost: 3,
				Faction: zb_enums.Faction_Water,
			},
			Owner: player1,
		}
		instance2 := &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 4},
			Prototype:  &zb_data.Card{},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 5,
				Damage:  4,
				Cost: 3,
				Faction: zb_enums.Faction_Fire,
			},
			Owner: player1,
		}
		gp.State.PlayerStates[0].CardsInHand = append(gp.State.PlayerStates[0].CardsInPlay, instance0)
		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance1)
		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance2)
		// gp.DebugState()
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
		gp.DebugState()
		assert.Equal(t, 3, len(gp.State.PlayerStates[0].CardsInPlay))
		assert.Equal(t, false, gp.State.PlayerStates[0].CardsInPlay[2].AbilitiesInstances[0].IsActive)
		assert.Equal(t, 1, len(gp.actionOutcomes[0].GetReplaceUnitsWithTypeOnStrongerOnes().NewCardInstances))
	})
}

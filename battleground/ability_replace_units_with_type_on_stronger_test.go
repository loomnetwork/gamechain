package battleground

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/stretchr/testify/assert"
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
		OverlordId: 2,
		Name:       "Default",
		Cards: []*zb_data.DeckCard{
			{MouldId: 133, Amount: 5},
			{MouldId: 136, Amount: 10},
		},
	}

	deck1 := &zb_data.Deck{
		Id:         0,
		OverlordId: 2,
		Name:       "Default",
		Cards: []*zb_data.DeckCard{
			{MouldId: 11, Amount: 15},
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

		card0 := &zb.Card{
			Name:    "Vortex",
			Defense: 3,
			Damage:  2,
			Cost: 1,
			Faction: zb.Faction_Water,
			Abilities: []*zb.AbilityData{
				{
					Ability: zb.AbilityType_ReplaceUnitsWithTypeOnStrongerOnes,
					Trigger: zb_enums.AbilityTrigger_Entry,
				},
			},
		}
		instance0 := &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 2},
			Instance:   newCardInstanceSpecificDataFromCardDetails(card0),
			Prototype:  proto.Clone(card0).(*zb.Card),
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
			Prototype:  &zb.Card{},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 5,
				Damage:  4,
				Cost: 3,
				Faction: zb.Faction_Water,
			},
			Owner: player1,
		}
		instance2 := &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 4},
			Prototype:  &zb.Card{},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 5,
				Damage:  4,
				Cost: 3,
				Faction: zb.Faction_Fire,
			},
			Owner: player1,
		}
		gp.State.PlayerStates[0].CardsInHand = append(gp.State.PlayerStates[0].CardsInPlay, instance0)
		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance1)
		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance2)
		// gp.DebugState()
		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardPlay{
				CardPlay: &zb.PlayerActionCardPlay{
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

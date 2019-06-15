package battleground

import (
	battleground_proto "github.com/loomnetwork/gamechain/battleground/proto"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	"testing"

	"github.com/gogo/protobuf/proto"
	loom "github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	assert "github.com/stretchr/testify/require"
)

func TestAbilityDevourZombieAndCombineStats(t *testing.T) {
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
			{CardKey: battleground_proto.CardKey{MouldId: 53}, Amount: 2},
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

	t.Run("DevourZombieAndCombineStat is active when enter the field, devouring a target", func(t *testing.T) {
		players := []*zb_data.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb_data.Card{
			Defense: 4,
			Damage:  2,
			Abilities: []*zb_data.AbilityData{
				{
					Ability: zb_enums.AbilityType_DevourZombiesAndCombineStats,
					Trigger: zb_enums.AbilityTrigger_Entry,
				},
			},
		}
		instance0 := &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 1},
			Instance:   newCardInstanceSpecificDataFromCardDetails(card0),
			Prototype:  proto.Clone(card0).(*zb_data.Card),
			AbilitiesInstances: []*zb_data.CardAbilityInstance{
				&zb_data.CardAbilityInstance{
					Trigger: card0.Abilities[0].Trigger,
					AbilityType: &zb_data.CardAbilityInstance_DevourZombieAndCombineStats{
						DevourZombieAndCombineStats: &zb_data.CardAbilityDevourZombieAndCombineStats{
							Faction: card0.Faction,
						},
					},
					IsActive: true,
				},
			},
		}
		instance1 := &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 2},
			Prototype:  &zb_data.Card{},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 2,
				Damage:  1,
			},
		}
		instance2 := &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 3},
			Prototype:  &zb_data.Card{},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 2,
				Damage:  1,
			},
		}

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance0, instance1, instance2)

		assert.Equal(t, int(3), len(gp.State.PlayerStates[0].CardsInPlay))

		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardAbilityUsed,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardAbilityUsed{
				CardAbilityUsed: &zb_data.PlayerActionCardAbilityUsed{
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
		assert.Equal(t, int(2), len(gp.State.PlayerStates[0].CardsInPlay))
		assert.Equal(t, int(1), len(gp.State.PlayerStates[0].CardsInGraveyard))
		assert.Equal(t, int32(6), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, int32(3), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Damage)

		// Try to use the ability again but this time it should not work
		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardAbilityUsed,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardAbilityUsed{
				CardAbilityUsed: &zb_data.PlayerActionCardAbilityUsed{
					Card: &zb_data.InstanceId{Id: 1},
					Targets: []*zb_data.Unit{
						&zb_data.Unit{
							InstanceId: &zb_data.InstanceId{Id: 2},
						},
					},
				},
			},
		})
		assert.Equal(t, int(2), len(gp.State.PlayerStates[0].CardsInPlay))
		assert.Equal(t, int(1), len(gp.State.PlayerStates[0].CardsInGraveyard))
		assert.Equal(t, int32(6), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, int32(3), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Damage)
	})

	t.Run("DevourZombieAndCombineStat is active when enter the field, devouring all ally zombies", func(t *testing.T) {
		players := []*zb_data.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb_data.Card{
			Defense: 4,
			Damage:  2,
			Abilities: []*zb_data.AbilityData{
				{
					Ability: zb_enums.AbilityType_DevourZombiesAndCombineStats,
					Trigger: zb_enums.AbilityTrigger_Entry,
				},
			},
		}
		instance0 := &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 1},
			Instance:   newCardInstanceSpecificDataFromCardDetails(card0),
			Prototype:  proto.Clone(card0).(*zb_data.Card),
			AbilitiesInstances: []*zb_data.CardAbilityInstance{
				&zb_data.CardAbilityInstance{
					Trigger: card0.Abilities[0].Trigger,
					AbilityType: &zb_data.CardAbilityInstance_DevourZombieAndCombineStats{
						DevourZombieAndCombineStats: &zb_data.CardAbilityDevourZombieAndCombineStats{},
					},
					IsActive: true,
				},
			},
		}
		instance1 := &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 2},
			Prototype:  &zb_data.Card{},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 2,
				Damage:  1,
			},
		}
		instance2 := &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 3},
			Prototype:  &zb_data.Card{},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 2,
				Damage:  1,
			},
		}

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance0, instance1, instance2)

		assert.Equal(t, int(3), len(gp.State.PlayerStates[0].CardsInPlay))

		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardAbilityUsed,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardAbilityUsed{
				CardAbilityUsed: &zb_data.PlayerActionCardAbilityUsed{
					Card: &zb_data.InstanceId{Id: 1},
					Targets: []*zb_data.Unit{
						&zb_data.Unit{
							InstanceId: &zb_data.InstanceId{Id: 2},
						},
						&zb_data.Unit{
							InstanceId: &zb_data.InstanceId{Id: 3},
						},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, int(1), len(gp.State.PlayerStates[0].CardsInPlay))
		assert.Equal(t, int(2), len(gp.State.PlayerStates[0].CardsInGraveyard))
		assert.Equal(t, int32(8), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, int32(4), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Damage)

	})
}

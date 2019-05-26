package battleground

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/stretchr/testify/assert"
)

func TestAbilityChangeStat(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setupInitFromFile(c, pubKeyHexString, &addr, &ctx, t)

	player1 := "player-1"
	player2 := "player-2"

	deck0 := &zb.Deck{
		Id:         0,
		OverlordId: 2,
		Name:       "Default",
		Cards: []*zb.DeckCard{
			{MouldId: 90, Amount: 2},
			{MouldId: 91, Amount: 2},
			{MouldId: 96, Amount: 2},
			{MouldId: 3, Amount: 2},
			{MouldId: 2, Amount: 2},
			{MouldId: 92, Amount: 2},
			{MouldId: 1, Amount: 1},
			{MouldId: 93, Amount: 1},
			{MouldId: 7, Amount: 1},
			{MouldId: 94, Amount: 1},
			{MouldId: 5, Amount: 1},
		},
	}

	t.Run("ChangeStat is activated when attacking a card", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb.Card{
			Defense: 5,
			Damage:  2,
			Abilities: []*zb.AbilityData{
				{
					Ability: zb.AbilityType_ChangeStat,
					Trigger: zb.AbilityTrigger_Attack,
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
					AbilityType: &zb.CardAbilityInstance_ChangeStat{
						ChangeStat: &zb.CardAbilityChangeStat{
							StatAdjustment: -1,
							Stat:           zb.Stat_Damage,
						},
					},
				},
				&zb.CardAbilityInstance{
					IsActive: true,
					Trigger:  card0.Abilities[0].Trigger,
					AbilityType: &zb.CardAbilityInstance_ChangeStat{
						ChangeStat: &zb.CardAbilityChangeStat{
							StatAdjustment: -1,
							Stat:           zb.Stat_Defense,
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
				Damage:  1,
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
		assert.Equal(t, int32(3), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, int32(1), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Damage)
		assert.Equal(t, int32(3), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, int32(1), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Damage)
	})

	t.Run("ChangeStat is activated when attacking overlord", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb.Card{
			Defense: 2,
			Damage:  3,
			Abilities: []*zb.AbilityData{
				{
					Ability: zb.AbilityType_ChangeStat,
					Trigger: zb.AbilityTrigger_Attack,
				},
			},
		}
		instance0 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 2},
			Instance:   newCardInstanceSpecificDataFromCardDetails(card0),
			Prototype:  proto.Clone(card0).(*zb.Card),
			AbilitiesInstances: []*zb.CardAbilityInstance{
				&zb.CardAbilityInstance{
					IsActive: true,
					Trigger:  card0.Abilities[0].Trigger,
					AbilityType: &zb.CardAbilityInstance_ChangeStat{
						ChangeStat: &zb.CardAbilityChangeStat{
							StatAdjustment: -1,
							Stat:           zb.Stat_Damage,
						},
					},
				},
				&zb.CardAbilityInstance{
					IsActive: true,
					Trigger:  card0.Abilities[0].Trigger,
					AbilityType: &zb.CardAbilityInstance_ChangeStat{
						ChangeStat: &zb.CardAbilityChangeStat{
							StatAdjustment: -1,
							Stat:           zb.Stat_Defense,
						},
					},
				},
			},
		}

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance0)

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
		assert.Equal(t, int32(1), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, int32(2), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Damage)

		assert.Equal(t, int32(47), gp.State.PlayerStates[1].Defense)
	})

}

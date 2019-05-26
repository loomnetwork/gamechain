package battleground

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/stretchr/testify/assert"
)

func TestAbilityDealDamageToThisAndAdjacentUnits(t *testing.T) {
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

	t.Run("DealDamageToThisAndAdjacentUnits should attack adjacent cards", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb.Card{
			Name:    "WiZp",
			Defense: 3,
			Damage:  2,
			Abilities: []*zb.AbilityData{
				{
					Ability: zb.AbilityType_DealDamageToThisAndAdjacentUnits,
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
					AbilityType: &zb.CardAbilityInstance_DealDamageToThisAndAdjacentUnits{
						DealDamageToThisAndAdjacentUnits: &zb.CardAbilityDealDamageToThisAndAdjacentUnits{
							AdjacentDamage: 2,
						},
					},
				},
			},
			Owner:      player1,
			OwnerIndex: 0,
		}
		instance1 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 3},
			Prototype: &zb.Card{
				Name: "target1",
			},
			Instance: &zb.CardInstanceSpecificData{
				Damage:  1,
				Defense: 1,
			},
			Owner:      player2,
			OwnerIndex: 1,
		}
		instance2 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 4},
			Prototype: &zb.Card{
				Name: "target2",
			},
			Instance: &zb.CardInstanceSpecificData{
				Damage:  1,
				Defense: 5,
			},
			Owner:      player2,
			OwnerIndex: 1,
		}
		instance3 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 5},
			Prototype: &zb.Card{
				Name: "target3",
			},
			Instance: &zb.CardInstanceSpecificData{
				Damage:  3,
				Defense: 3,
			},
			Owner:      player2,
			OwnerIndex: 1,
		}

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance0)
		gp.State.PlayerStates[1].CardsInPlay = []*zb.CardInstance{instance1, instance2, instance3}

		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardAttack,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardAttack{
				CardAttack: &zb.PlayerActionCardAttack{
					Attacker: &zb.InstanceId{Id: 2},
					Target: &zb.Unit{
						InstanceId: &zb.InstanceId{Id: 4},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(2), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Damage)
		assert.Equal(t, int32(2), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, 2, len(gp.State.PlayerStates[1].CardsInPlay), "player2 cards in play should be 2")
		assert.Equal(t, int32(1), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Damage)
		assert.Equal(t, int32(3), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, int32(3), gp.State.PlayerStates[1].CardsInPlay[1].Instance.Damage)
		assert.Equal(t, int32(1), gp.State.PlayerStates[1].CardsInPlay[1].Instance.Defense)
	})

	t.Run("DealDamageToThisAndAdjacentUnits should attack left card", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb.Card{
			Name:    "WiZp",
			Defense: 3,
			Damage:  2,
			Abilities: []*zb.AbilityData{
				{
					Ability: zb.AbilityType_DealDamageToThisAndAdjacentUnits,
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
					AbilityType: &zb.CardAbilityInstance_DealDamageToThisAndAdjacentUnits{
						DealDamageToThisAndAdjacentUnits: &zb.CardAbilityDealDamageToThisAndAdjacentUnits{
							AdjacentDamage: 2,
						},
					},
				},
			},
			Owner:      player1,
			OwnerIndex: 0,
		}
		instance1 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 3},
			Prototype: &zb.Card{
				Name: "target1",
			},
			Instance: &zb.CardInstanceSpecificData{
				Damage:  1,
				Defense: 1,
			},
			Owner:      player2,
			OwnerIndex: 1,
		}
		instance2 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 4},
			Prototype: &zb.Card{
				Name: "target2",
			},
			Instance: &zb.CardInstanceSpecificData{
				Damage:  1,
				Defense: 5,
			},
			Owner:      player2,
			OwnerIndex: 1,
		}

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance0)
		gp.State.PlayerStates[1].CardsInPlay = []*zb.CardInstance{instance1, instance2}

		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardAttack,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardAttack{
				CardAttack: &zb.PlayerActionCardAttack{
					Attacker: &zb.InstanceId{Id: 2},
					Target: &zb.Unit{
						InstanceId: &zb.InstanceId{Id: 4},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(2), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Damage)
		assert.Equal(t, int32(2), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, 1, len(gp.State.PlayerStates[1].CardsInPlay), "player2 cards in play should be 1")
		assert.Equal(t, int32(1), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Damage)
		assert.Equal(t, int32(3), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Defense)
	})
}

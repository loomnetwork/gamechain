package battleground

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/stretchr/testify/assert"
)

func TestAbilityReanimate(t *testing.T) {
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

	t.Run("Reanimate ability get activated when attacker death", func(t *testing.T) {
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
			Attack:  2,
			Abilities: []*zb.CardAbility{
				{
					Type:    zb.CardAbilityType_ReanimateUnit,
					Trigger: zb.CardAbilityTrigger_Death,
				},
				{
					Type:    zb.CardAbilityType_ReplaceUnitsWithTypeOnStrongerOnes,
					Trigger: zb.CardAbilityTrigger_Entry,
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
					AbilityType: &zb.CardAbilityInstance_Reanimate{
						Reanimate: &zb.CardAbilityReanimate{
							DefaultAttack:  card0.Attack,
							DefaultDefense: card0.Defense,
						},
					},
				},
				&zb.CardAbilityInstance{
					IsActive: true,
					Trigger:  card0.Abilities[1].Trigger,
					AbilityType: &zb.CardAbilityInstance_ReplaceUnitsWithTypeOnStrongerOnes{
						ReplaceUnitsWithTypeOnStrongerOnes: &zb.CardAbilityReplaceUnitsWithTypeOnStrongerOnes{
							Faction: zb.Faction_Water,
						},
					},
				},
			},
			Owner: player1,
		}
		instance1 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 3},
			Prototype: &zb.Card{
				Name: "target",
			},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 5,
				Attack:  4,
			},
			Owner: player2,
		}

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance0)
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, instance1)
		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardAttack,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardAttack{
				CardAttack: &zb.PlayerActionCardAttack{
					Attacker: &zb.InstanceId{Id: 2},
					Target: &zb.Unit{
						InstanceId: &zb.InstanceId{Id: 3},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(2), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Attack)
		assert.Equal(t, int32(3), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, 1, len(gp.State.PlayerStates[0].CardsInPlay[0].AbilitiesInstances), "renaimate should not be in new instance")
		assert.NotNil(t, gp.State.PlayerStates[0].CardsInPlay[0].AbilitiesInstances[0].GetReplaceUnitsWithTypeOnStrongerOnes, "ReplaceUnitsWithTypeOnStrongerOne should be on the new instance")
		assert.Equal(t, int32(2), gp.actionOutcomes[0].GetReanimate().NewCardInstance.Instance.Attack)
		assert.Equal(t, int32(3), gp.actionOutcomes[0].GetReanimate().NewCardInstance.Instance.Defense)
	})
}

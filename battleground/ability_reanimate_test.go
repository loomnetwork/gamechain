package battleground

import (
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	"testing"

	"github.com/gogo/protobuf/proto"
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

	deck0 := &zb_data.Deck{
		Id:         0,
		OverlordId: 2,
		Name:       "Default",
		Cards: []*zb_data.DeckCard{
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

	t.Run("Reanimate ability get activated when attacker death", func(t *testing.T) {
		players := []*zb_data.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb_data.Card{
			Name:    "WiZp",
			Defense: 3,
			Damage:  2,
			Abilities: []*zb_data.AbilityData{
				{
					Ability: zb_enums.AbilityType_ReanimateUnit,
					Trigger: zb_enums.AbilityTrigger_Death,
				},
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
					AbilityType: &zb_data.CardAbilityInstance_Reanimate{
						Reanimate: &zb_data.CardAbilityReanimate{
							DefaultDamage:  card0.Damage,
							DefaultDefense: card0.Defense,
						},
					},
				},
				&zb_data.CardAbilityInstance{
					IsActive: true,
					Trigger:  card0.Abilities[1].Trigger,
					AbilityType: &zb_data.CardAbilityInstance_ReplaceUnitsWithTypeOnStrongerOnes{
						ReplaceUnitsWithTypeOnStrongerOnes: &zb_data.CardAbilityReplaceUnitsWithTypeOnStrongerOnes{
							Faction: zb_enums.Faction_Water,
						},
					},
				},
			},
			Owner: player1,
		}
		instance1 := &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 3},
			Prototype: &zb_data.Card{
				Name: "target",
			},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 5,
				Damage:  4,
			},
			Owner: player2,
		}

		gp.State.PlayerStates[0].CardsInPlay = append(gp.State.PlayerStates[0].CardsInPlay, instance0)
		gp.State.PlayerStates[1].CardsInPlay = append(gp.State.PlayerStates[1].CardsInPlay, instance1)
		err = gp.AddAction(&zb_data.PlayerAction{
			ActionType: zb_enums.PlayerActionType_CardAttack,
			PlayerId:   player1,
			Action: &zb_data.PlayerAction_CardAttack{
				CardAttack: &zb_data.PlayerActionCardAttack{
					Attacker: &zb_data.InstanceId{Id: 2},
					Target: &zb_data.Unit{
						InstanceId: &zb_data.InstanceId{Id: 3},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(2), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Damage)
		assert.Equal(t, int32(3), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, 1, len(gp.State.PlayerStates[0].CardsInPlay[0].AbilitiesInstances), "renaimate should not be in new instance")
		assert.NotNil(t, gp.State.PlayerStates[0].CardsInPlay[0].AbilitiesInstances[0].GetReplaceUnitsWithTypeOnStrongerOnes, "ReplaceUnitsWithTypeOnStrongerOne should be on the new instance")
		assert.Equal(t, int32(2), gp.actionOutcomes[0].GetReanimate().NewCardInstance.Instance.Damage)
		assert.Equal(t, int32(3), gp.actionOutcomes[0].GetReanimate().NewCardInstance.Instance.Defense)
	})
}

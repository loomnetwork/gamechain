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

func TestAbilityRage(t *testing.T) {
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

	t.Run("Rage ability works", func(t *testing.T) {
		players := []*zb_data.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		instance0 := &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 2},
			Prototype: &zb_data.Card{
				Name: "attacker",
			},
			Instance: &zb_data.CardInstanceSpecificData{
				Defense: 5,
				Damage:  1,
			},
		}

		card1 := &zb_data.Card{
			Defense: 5,
			Damage:  2,
			Abilities: []*zb_data.AbilityData{
				{
					Ability: zb_enums.AbilityType_Rage,
					Trigger: zb_enums.AbilityTrigger_GotDamage,
					Value:   2,
				},
			},
		}
		instance1 := &zb_data.CardInstance{
			InstanceId: &zb_data.InstanceId{Id: 3},
			Instance:   newCardInstanceSpecificDataFromCardDetails(card1),
			Prototype:  proto.Clone(card1).(*zb_data.Card),
			AbilitiesInstances: []*zb_data.CardAbilityInstance{
				&zb_data.CardAbilityInstance{
					IsActive: true,
					Trigger:  card1.Abilities[0].Trigger,
					AbilityType: &zb_data.CardAbilityInstance_Rage{
						Rage: &zb_data.CardAbilityRage{
							AddedDamage: card1.Abilities[0].Value,
						},
					},
				},
			},
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
		assert.Equal(t, int32(4), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Damage)
		assert.Equal(t, int32(4), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, int32(3), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, int32(4), gp.actionOutcomes[0].GetRage().NewDamage)
		assert.Equal(t, int32(3), gp.actionOutcomes[0].GetRage().InstanceId.Id)
	})
}

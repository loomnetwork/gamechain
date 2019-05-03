package battleground

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/stretchr/testify/assert"
)

func TestAbilityAttackOverlord(t *testing.T) {
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

	t.Run("Player overlord is damaged when the card is played", func(t *testing.T) {
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
					Ability: zb.AbilityType_AttackOverlord,
					Trigger: zb.AbilityTrigger_Entry,
				},
			},
		}
		instance0 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 100},
			Instance:   newCardInstanceSpecificDataFromCardDetails(card0),
			Prototype:  proto.Clone(card0).(*zb.Card),
			AbilitiesInstances: []*zb.CardAbilityInstance{
				&zb.CardAbilityInstance{
					IsActive: true,
					Trigger:  card0.Abilities[0].Trigger,
					AbilityType: &zb.CardAbilityInstance_AttackOverlord{
						AttackOverlord: &zb.CardAbilityAttackOverlord{
							Damage:     2,
							WasApplied: false,
						},
					},
				},
			},
		}

		gp.State.PlayerStates[0].CardsInHand = append(gp.State.PlayerStates[0].CardsInHand, instance0)

		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardPlay{
				CardPlay: &zb.PlayerActionCardPlay{
					Card: &zb.InstanceId{Id: 100},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(48), gp.State.PlayerStates[0].Defense)

		instance1 := &zb.CardInstance{
			InstanceId:         &zb.InstanceId{Id: 101},
			Instance:           newCardInstanceSpecificDataFromCardDetails(card0),
			Prototype:          proto.Clone(card0).(*zb.Card),
			AbilitiesInstances: []*zb.CardAbilityInstance{},
		}

		gp.State.PlayerStates[0].CardsInHand = append(gp.State.PlayerStates[0].CardsInHand, instance1)

		err = gp.AddAction(&zb.PlayerAction{
			ActionType: zb.PlayerActionType_CardPlay,
			PlayerId:   player1,
			Action: &zb.PlayerAction_CardPlay{
				CardPlay: &zb.PlayerActionCardPlay{
					Card: &zb.InstanceId{Id: 101},
				},
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(48), gp.State.PlayerStates[0].Defense)

	})
}

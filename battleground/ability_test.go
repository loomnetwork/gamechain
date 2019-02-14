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
			Attack:  2,
			Abilities: []*zb.CardAbility{
				{
					Type:    zb.CardAbilityType_ChangeStat,
					Trigger: zb.CardAbilityTrigger_Attack,
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
							Stat:           zb.StatType_Attack,
						},
					},
				},
				&zb.CardAbilityInstance{
					IsActive: true,
					Trigger:  card0.Abilities[0].Trigger,
					AbilityType: &zb.CardAbilityInstance_ChangeStat{
						ChangeStat: &zb.CardAbilityChangeStat{
							StatAdjustment: -1,
							Stat:           zb.StatType_Defense,
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
				Attack:  1,
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
		assert.Equal(t, int32(1), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Attack)
		assert.Equal(t, int32(3), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, int32(1), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Attack)
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
			Attack:  3,
			Abilities: []*zb.CardAbility{
				{
					Type:    zb.CardAbilityType_ChangeStat,
					Trigger: zb.CardAbilityTrigger_Attack,
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
							Stat:           zb.StatType_Attack,
						},
					},
				},
				&zb.CardAbilityInstance{
					IsActive: true,
					Trigger:  card0.Abilities[0].Trigger,
					AbilityType: &zb.CardAbilityInstance_ChangeStat{
						ChangeStat: &zb.CardAbilityChangeStat{
							StatAdjustment: -1,
							Stat:           zb.StatType_Defense,
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
		assert.Equal(t, int32(2), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Attack)

		assert.Equal(t, int32(17), gp.State.PlayerStates[1].Defense)
	})

}

func TestAbilityAttackOverlord(t *testing.T) {
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
			Attack:  2,
			Abilities: []*zb.CardAbility{
				{
					Type:    zb.CardAbilityType_AttackOverlord,
					Trigger: zb.CardAbilityTrigger_Entry,
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
		assert.Equal(t, int32(18), gp.State.PlayerStates[0].Defense)

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
		assert.Equal(t, int32(18), gp.State.PlayerStates[0].Defense)

	})
}

func TestAbilityAdditionalDamageToHeavyInAttack(t *testing.T) {
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

	t.Run("AdditionalDamageToHeavyInAttack ability does trigger when target is heavy", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb.Card{
			Defense: 5,
			Attack:  2,
			Abilities: []*zb.CardAbility{
				{
					Type:    zb.CardAbilityType_AdditionalDamageToHeavyInAttack,
					Trigger: zb.CardAbilityTrigger_Attack,
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
					AbilityType: &zb.CardAbilityInstance_AdditionalDamageToHeavyInAttack{
						AdditionalDamageToHeavyInAttack: &zb.CardAbilityAdditionalDamageToHeavyInAttack{
							AddedAttack: 2,
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
				Attack:  1,
				Type:    zb.CreatureType_Heavy,
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
		assert.Equal(t, int32(1), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Defense)
	})

	t.Run("AdditionalDamageToHeavyInAttack ability does NOT trigger when target is NOT heavy", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb.Card{
			Defense: 5,
			Attack:  2,
			Abilities: []*zb.CardAbility{
				{
					Type:    zb.CardAbilityType_AdditionalDamageToHeavyInAttack,
					Trigger: zb.CardAbilityTrigger_Attack,
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
					AbilityType: &zb.CardAbilityInstance_AdditionalDamageToHeavyInAttack{
						AdditionalDamageToHeavyInAttack: &zb.CardAbilityAdditionalDamageToHeavyInAttack{
							AddedAttack: 2,
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
				Attack:  1,
				Type:    zb.CreatureType_Feral,
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
		assert.Equal(t, int32(3), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Defense)
	})
}

func TestAbilityRage(t *testing.T) {
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

	t.Run("Rage ability works", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb.Card{
			Defense: 5,
			Attack:  2,
			Abilities: []*zb.CardAbility{
				{
					Type:    zb.CardAbilityType_Rage,
					Trigger: zb.CardAbilityTrigger_GotDamage,
					Value:   2,
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
					AbilityType: &zb.CardAbilityInstance_Rage{
						Rage: &zb.CardAbilityRage{
							AddedAttack: card0.Abilities[0].Value,
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
				Attack:  1,
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
		assert.Equal(t, int32(4), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Attack)
		assert.Equal(t, int32(4), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, int32(1), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, int32(4), gp.actionOutcomes[0].GetRage().NewAttack)
		assert.Equal(t, int32(1), gp.actionOutcomes[0].GetRage().InstanceId.Id)
	})
}

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
			Name:    "Reanimate",
			Defense: 3,
			Attack:  2,
			Abilities: []*zb.CardAbility{
				{
					Type:    zb.CardAbilityType_ReanimateUnit,
					Trigger: zb.CardAbilityTrigger_Death,
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
					AbilityType: &zb.CardAbilityInstance_Reanimate{
						Reanimate: &zb.CardAbilityReanimate{
							DefaultAttack:  card0.Attack,
							DefaultDefense: card0.Defense,
						},
					},
				},
			},
			Owner: player1,
		}
		instance1 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 2},
			Prototype:  &zb.Card{},
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
					Attacker: &zb.InstanceId{Id: 1},
					Target: &zb.Unit{
						InstanceId: &zb.InstanceId{Id: 2},
					},
				},
			},
		})

		assert.Nil(t, err)
		assert.Equal(t, int32(2), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Attack)
		assert.Equal(t, int32(3), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, false, gp.State.PlayerStates[0].CardsInPlay[0].AbilitiesInstances[0].IsActive)
		assert.Equal(t, int32(2), gp.actionOutcomes[0].GetReanimate().NewCardInstance.Instance.Attack)
		assert.Equal(t, int32(3), gp.actionOutcomes[0].GetReanimate().NewCardInstance.Instance.Defense)
	})
}

func TestAbilityPriorityAttack(t *testing.T) {
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

	t.Run("PriorityAttack ability does not trigger when target is not killed", func(t *testing.T) {
		players := []*zb.PlayerState{
			{Id: player1, Deck: deck0},
			{Id: player2, Deck: deck0},
		}
		seed := int64(0)
		gp, err := NewGamePlay(ctx, 3, "v1", players, seed, nil, true, nil)
		assert.Nil(t, err)

		card0 := &zb.Card{
			Defense: 5,
			Attack:  2,
			Abilities: []*zb.CardAbility{
				{
					Type: zb.CardAbilityType_PriorityAttack,
				},
			},
		}
		instance0 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 1},
			Instance:   newCardInstanceSpecificDataFromCardDetails(card0),
			Prototype:  proto.Clone(card0).(*zb.Card),
		}
		instance1 := &zb.CardInstance{
			InstanceId: &zb.InstanceId{Id: 2},
			Prototype:  &zb.Card{},
			Instance: &zb.CardInstanceSpecificData{
				Defense: 5,
				Attack:  1,
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
		assert.Equal(t, int32(4), gp.State.PlayerStates[0].CardsInPlay[0].Instance.Defense)
		assert.Equal(t, int32(3), gp.State.PlayerStates[1].CardsInPlay[0].Instance.Defense)
	})
}

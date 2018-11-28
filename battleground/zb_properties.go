package battleground

import (
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
)

type CardInstance struct {
	*zb.CardInstance
}

type CardAbility interface {
	defenseChangedHandler(card *CardInstance) []*zb.PlayerActionOutcome
}

type CardAbilityRage struct {
	*zb.CardAbilityRage
}

type CardAbilityPriorityAttack struct {
	*zb.CardAbilityPriorityAttack
}

type abilityInstanceFn func(game *Gameplay, ability CardAbility, card *CardInstance) []*zb.PlayerActionOutcome

// SetDefence will set the card's defense value and call the ability's defenseChangedHandler
func (card *CardInstance) SetDefense(game *Gameplay, defense int32) {
	card.tryInitAbilitiesInstances()
	card.Instance.Defense = defense

	callAbilityInstancesFunc(game, card, func(game *Gameplay, ability CardAbility, card *CardInstance) []*zb.PlayerActionOutcome {
		return ability.defenseChangedHandler(card)
	})

	fmt.Printf("\n\ngame.actionOutcomes: %v\n\n", game.actionOutcomes)
}

func callAbilityInstancesFunc(game *Gameplay, card *CardInstance, fn abilityInstanceFn) {
	for _, abilityInstanceRaw := range card.AbilitiesInstances {
		var abilityInstance CardAbility
		fmt.Println(abilityInstanceRaw.AbilityType)
		switch abilityType := abilityInstanceRaw.AbilityType.(type) {
		case *zb.CardAbilityInstance_Rage:
			abilityInstance = &CardAbilityRage{abilityType.Rage}
		case *zb.CardAbilityInstance_PriorityAttack:
			abilityInstance = &CardAbilityPriorityAttack{abilityType.PriorityAttack}
		default:
			fmt.Println("CardAbilityInstance has unexpected type %T", abilityType)
		}

		outcomes := fn(game, abilityInstance, card)
		if outcomes != nil {
			for _, outcome := range outcomes {
				game.actionOutcomes = append(game.actionOutcomes, outcome)
			}
		}
	}
}

func (card *CardInstance) initAbilityInstances() {
	if card.Prototype.Abilities == nil {
		return
	}

	for _, abilityInstanceRaw := range card.Prototype.Abilities {
		switch abilityInstanceRaw.Type {
		case zb.CardAbilityType_Rage:
			card.AbilitiesInstances = append(card.AbilitiesInstances, &zb.CardAbilityInstance{
				AbilityType: &zb.CardAbilityInstance_Rage{
					Rage: &zb.CardAbilityRage{
						WasApplied:  false,
						AddedAttack: abilityInstanceRaw.Value,
					},
				},
			})
		case zb.CardAbilityType_PriorityAttack:
			card.AbilitiesInstances = append(card.AbilitiesInstances, &zb.CardAbilityInstance{
				AbilityType: &zb.CardAbilityInstance_PriorityAttack{
					PriorityAttack: &zb.CardAbilityPriorityAttack{},
				},
			})
		default:
			fmt.Printf("CardAbility.Type has unexpected value %d", abilityInstanceRaw.Type)
		}
	}
}

func (card *CardInstance) tryInitAbilitiesInstances() {
	if !card.AbilitiesInstancesInitialized {
		card.initAbilityInstances()
		card.AbilitiesInstancesInitialized = true
	}
}

/* ability specific hooks */
// Rage
func (rage *CardAbilityRage) defenseChangedHandler(card *CardInstance) []*zb.PlayerActionOutcome {
	if !rage.WasApplied {
		if card.Instance.Defense < card.Prototype.Defense {
			rage.WasApplied = true
			card.Instance.Attack += rage.AddedAttack

			return []*zb.PlayerActionOutcome{
				{
					Outcome: &zb.PlayerActionOutcome_Rage{
						Rage: &zb.PlayerActionOutcome_CardAbilityRageOutcome{
							InstanceId: card.InstanceId,
							NewAttack:  card.Instance.Attack,
						},
					},
				},
			}
		}
	} else if card.Instance.Defense >= card.Prototype.Defense {
		rage.WasApplied = false
		card.Instance.Attack -= rage.AddedAttack

		return []*zb.PlayerActionOutcome{
			{
				Outcome: &zb.PlayerActionOutcome_Rage{
					Rage: &zb.PlayerActionOutcome_CardAbilityRageOutcome{
						InstanceId: card.InstanceId,
						NewAttack:  card.Instance.Attack,
					},
				},
			},
		}
	}

	return nil
}

// Priority Attack
func (priorityAttack *CardAbilityPriorityAttack) defenseChangedHandler(card *CardInstance) []*zb.PlayerActionOutcome {
	// reset the card's defense
	fmt.Println("prAttack")
	// TODO: need the older value here
	card.Instance.Defense = card.Prototype.Defense
	return []*zb.PlayerActionOutcome{}
}

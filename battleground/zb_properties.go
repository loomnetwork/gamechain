package battleground

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb"
)

type CardInstance struct {
	*zb.CardInstance
}

type CardAbility interface {
	defenseChangedHandler(card *CardInstance)
}

type CardAbilityRage struct {
	*zb.CardAbilityRage
}

func (card *CardInstance) SetDefense(defense int32) {
	card.tryInitAbilitiesInstances()
	fmt.Println(card)
	card.Instance.Defense = defense
	fmt.Println(card)

	callAbilityInstancesFunc(card, func(ability CardAbility, card *CardInstance) {
		ability.defenseChangedHandler(card)
	})
}

func (rage *CardAbilityRage) defenseChangedHandler(card *CardInstance) {
	if !rage.WasApplied {
		if card.Instance.Defense < card.Prototype.Defense {
			rage.WasApplied = true
			card.Instance.Attack += rage.AddedAttack
		}
	} else {
		if card.Instance.Defense >= card.Prototype.Defense {
			rage.WasApplied = false
			card.Instance.Attack -= rage.AddedAttack
		}
	}
}

type abilityInstanceFn func(ability CardAbility, card *CardInstance)

func callAbilityInstancesFunc(card *CardInstance, fn abilityInstanceFn) {
	for _, abilityInstanceRaw := range card.AbilitiesInstances {
		var abilityInstance CardAbility
		switch abilityType := abilityInstanceRaw.AbilityType.(type) {
		case *zb.CardAbilityInstance_Rage:
			abilityInstance = &CardAbilityRage{abilityType.Rage}
		default:
			panic(fmt.Errorf("CardAbilityInstance has unexpected type %T", abilityType))
		}

		fn(abilityInstance, card)
	}
}

func (card *CardInstance) initAbilityInstances() {
	for _, abilityInstanceRaw := range card.Instance.Abilities {
		switch abilityInstanceRaw.Type {
		case zb.CardAbilityType_Rage:
			card.AbilitiesInstances = append(card.AbilitiesInstances, &zb.CardAbilityInstance{
				AbilityType: &zb.CardAbilityInstance_Rage{
					Rage: &zb.CardAbilityRage{
						WasApplied: false,
						AddedAttack: abilityInstanceRaw.Value,
					},
				},
			})
		default:
			panic(fmt.Errorf("CardAbility.Type has unexpected value %d", abilityInstanceRaw.Type))
		}
	}
}

func (card *CardInstance) tryInitAbilitiesInstances() {
	if !card.AbilitiesInstancesInitialized {
		card.initAbilityInstances()
		card.AbilitiesInstancesInitialized = true
	}
}
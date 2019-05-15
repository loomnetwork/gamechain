package battleground

import (
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
)

// rage ability
// description:
//     when being attacked, if not dies, increase card damage
type rage struct {
	*CardInstance
	cardAbility *zb_data.CardAbilityRage
}

var _ Ability = &rage{}

func NewRage(card *CardInstance, cardAbility *zb_data.CardAbilityRage) *rage {
	return &rage{
		CardInstance: card,
		cardAbility:  cardAbility,
	}
}

func (c *rage) Apply(gameplay *Gameplay) error {
	rage := c.cardAbility
	if !rage.WasApplied {
		if c.Instance.Defense < c.Prototype.Defense {
			// change state
			rage.WasApplied = true
			c.Instance.Damage += rage.AddedDamage
			// generate outcome
			gameplay.actionOutcomes = append(gameplay.actionOutcomes, &zb_data.PlayerActionOutcome{
				Outcome: &zb_data.PlayerActionOutcome_Rage{
					Rage: &zb_data.PlayerActionOutcome_CardAbilityRageOutcome{
						InstanceId: c.InstanceId,
						NewDamage:  c.Instance.Damage,
					},
				},
			})
		}
	} else if c.Instance.Defense >= c.Prototype.Defense {
		rage.WasApplied = false
		c.Instance.Damage -= rage.AddedDamage
		gameplay.actionOutcomes = append(gameplay.actionOutcomes, &zb_data.PlayerActionOutcome{
			Outcome: &zb_data.PlayerActionOutcome_Rage{
				Rage: &zb_data.PlayerActionOutcome_CardAbilityRageOutcome{
					InstanceId: c.InstanceId,
					NewDamage:  c.Instance.Damage,
				},
			},
		})
	}

	return nil
}

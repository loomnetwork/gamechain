package battleground

import (
	"github.com/loomnetwork/gamechain/types/zb"
)

// rage ability
// description:
//     when being attacked, if not dies, increase card attack
type rage struct {
	*CardInstance
	cardAbility *zb.CardAbilityRage
}

var _ Ability = &reanimate{}

func NewRage(card *CardInstance, cardAbility *zb.CardAbilityRage) *rage {
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
			c.Instance.Attack += rage.AddedAttack
			// generate outcome
			gameplay.actionOutcomes = append(gameplay.actionOutcomes, &zb.PlayerActionOutcome{
				Outcome: &zb.PlayerActionOutcome_Rage{
					Rage: &zb.PlayerActionOutcome_CardAbilityRageOutcome{
						InstanceId: c.InstanceId,
						NewAttack:  c.Instance.Attack,
					},
				},
			})
		}
	} else if c.Instance.Defense >= c.Prototype.Defense {
		rage.WasApplied = false
		c.Instance.Attack -= rage.AddedAttack
		gameplay.actionOutcomes = append(gameplay.actionOutcomes, &zb.PlayerActionOutcome{
			Outcome: &zb.PlayerActionOutcome_Rage{
				Rage: &zb.PlayerActionOutcome_CardAbilityRageOutcome{
					InstanceId: c.InstanceId,
					NewAttack:  c.Instance.Attack,
				},
			},
		})
	}

	return nil
}

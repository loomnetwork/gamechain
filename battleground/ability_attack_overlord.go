package battleground

import (
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
)

// attackOverlord ability
// description:
//     able to attack overlord
type attackOverlord struct {
	*CardInstance
	cardAbility *zb_data.CardAbilityAttackOverlord
}

var _ Ability = &attackOverlord{}

func NewAttackOverlord(card *CardInstance, cardAbility *zb_data.CardAbilityAttackOverlord) *attackOverlord {
	return &attackOverlord{
		CardInstance: card,
		cardAbility:  cardAbility,
	}
}

func (c *attackOverlord) Apply(gameplay *Gameplay) error {
	attackOverlord := c.cardAbility
	if !attackOverlord.WasApplied {
		// damage player overlord
		gameplay.activePlayer().Defense -= attackOverlord.Damage
		attackOverlord.WasApplied = true
		// generate outcome
		gameplay.actionOutcomes = append(gameplay.actionOutcomes, &zb_data.PlayerActionOutcome{
			Outcome: &zb_data.PlayerActionOutcome_AttackOverlord{
				AttackOverlord: &zb_data.PlayerActionOutcome_CardAbilityAttackOverlordOutcome{
					InstanceId: gameplay.activePlayer().InstanceId,
					NewDefense: gameplay.activePlayer().Defense,
				},
			},
		})
	}
	return nil
}

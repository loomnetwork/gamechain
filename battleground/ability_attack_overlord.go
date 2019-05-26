package battleground

import "github.com/loomnetwork/gamechain/types/zb"

// attackOverlord ability
// description:
//     able to attack overlord
type attackOverlord struct {
	*CardInstance
	cardAbility *zb.CardAbilityAttackOverlord
}

var _ Ability = &attackOverlord{}

func NewAttackOverlord(card *CardInstance, cardAbility *zb.CardAbilityAttackOverlord) *attackOverlord {
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
		gameplay.actionOutcomes = append(gameplay.actionOutcomes, &zb.PlayerActionOutcome{
			Outcome: &zb.PlayerActionOutcome_AttackOverlord{
				AttackOverlord: &zb.PlayerActionOutcome_CardAbilityAttackOverlordOutcome{
					InstanceId: gameplay.activePlayer().InstanceId,
					NewDefense: gameplay.activePlayer().Defense,
				},
			},
		})
	}
	return nil
}

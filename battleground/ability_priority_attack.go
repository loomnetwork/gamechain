package battleground

import (
	"github.com/loomnetwork/gamechain/types/zb"
)

// priority attack ability
// description:
//     reset the card's defense to the value before the attack, only if the opponent card dies
type priorityAttack struct {
	*CardInstance
	cardAbility *zb_data.CardAbilityPriorityAttack
}

var _ Ability = &priorityAttack{}

func NewPriorityAttack(card *CardInstance, cardAbility *zb_data.CardAbilityPriorityAttack) *priorityAttack {
	return &priorityAttack{
		CardInstance: card,
		cardAbility:  cardAbility,
	}
}

func (c *priorityAttack) Apply(gameplay *Gameplay) error {
	priorityAttack := c.cardAbility
	c.Instance.Defense = priorityAttack.AttackerOldDefense
	return nil
}

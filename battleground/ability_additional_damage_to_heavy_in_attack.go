package battleground

import "github.com/loomnetwork/gamechain/types/zb"

// additionalDamgeToHeavyInAttack ability
// description:
//     If the card is heavy, add addtional attack to defense
type additionalDamgeToHeavyInAttack struct {
	*CardInstance
	cardAbility *zb.CardAbilityAdditionalDamageToHeavyInAttack
}

var _ Ability = &attackOverlord{}

func NewAdditionalDamgeToHeavyInAttack(card *CardInstance, cardAbility *zb.CardAbilityAdditionalDamageToHeavyInAttack) *additionalDamgeToHeavyInAttack {
	return &additionalDamgeToHeavyInAttack{
		CardInstance: card,
		cardAbility:  cardAbility,
	}
}

func (c *additionalDamgeToHeavyInAttack) Apply(gameplay *Gameplay) error {
	additionalDamageToHeavyInAttack := c.cardAbility
	if c.Instance.Type == zb.CreatureType_Heavy {
		c.Instance.Defense -= additionalDamageToHeavyInAttack.AddedAttack
	}
	return nil
}

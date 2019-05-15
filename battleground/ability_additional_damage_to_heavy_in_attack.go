package battleground

import "github.com/loomnetwork/gamechain/types/zb"

// additionalDamgeToHeavyInAttack ability
// description:
//     If the card is heavy, add addtional damage to defense
type additionalDamgeToHeavyInAttack struct {
	*CardInstance
	cardAbility *zb_data.CardAbilityAdditionalDamageToHeavyInAttack
	target      *CardInstance
}

var _ Ability = &additionalDamgeToHeavyInAttack{}

func NewAdditionalDamgeToHeavyInAttack(card *CardInstance, cardAbility *zb_data.CardAbilityAdditionalDamageToHeavyInAttack, target *CardInstance) *additionalDamgeToHeavyInAttack {
	return &additionalDamgeToHeavyInAttack{
		CardInstance: card,
		cardAbility:  cardAbility,
		target:       target,
	}
}

func (c *additionalDamgeToHeavyInAttack) Apply(gameplay *Gameplay) error {
	additionalDamageToHeavyInAttack := c.cardAbility
	if c.target.Instance.Type == zb.CardType_Heavy {
		c.target.Instance.Defense -= additionalDamageToHeavyInAttack.AddedDamage
	}
	return nil
}

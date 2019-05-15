package battleground

import "github.com/loomnetwork/gamechain/types/zb/zb_data"
import "github.com/loomnetwork/gamechain/types/zb/zb_enums"

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
	if c.target.Instance.Type == zb_enums.CardType_Heavy {
		c.target.Instance.Defense -= additionalDamageToHeavyInAttack.AddedDamage
	}
	return nil
}

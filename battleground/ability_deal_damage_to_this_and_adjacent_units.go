package battleground

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
)

// dealDamageToThisAndAdjacentUnits ability
// description:
//     Deals damage to any zombies to the left and right of the target.
type dealDamageToThisAndAdjacentUnits struct {
	*CardInstance
	cardAbility *zb_data.CardAbilityDealDamageToThisAndAdjacentUnits
	target      *CardInstance
}

var _ Ability = &dealDamageToThisAndAdjacentUnits{}

func NewDealDamageToThisAndAdjacentUnits(card *CardInstance, cardAbility *zb_data.CardAbilityDealDamageToThisAndAdjacentUnits, target *CardInstance) *dealDamageToThisAndAdjacentUnits {
	return &dealDamageToThisAndAdjacentUnits{
		CardInstance: card,
		cardAbility:  cardAbility,
		target:       target,
	}
}

func (c *dealDamageToThisAndAdjacentUnits) Apply(gameplay *Gameplay) error {
	target := c.target
	// find left and right of the target
	owner := target.Player()
	if owner == nil {
		return fmt.Errorf("no owner for card instance %d", c.InstanceId)
	}

	// if only one or zero card in play, do nothing
	if len(owner.CardsInPlay) <= 1 {
		return nil
	}

	index := -1
	for i, card := range owner.CardsInPlay {
		if proto.Equal(card.InstanceId, target.InstanceId) {
			index = i
			break
		}
	}
	if index < 0 {
		return fmt.Errorf("card not found in play %d", c.InstanceId)
	}

	var left, right *zb_data.CardInstance
	if index > 0 {
		left = owner.CardsInPlay[index-1]
	}
	if index+1 < len(owner.CardsInPlay) {
		right = owner.CardsInPlay[index+1]
	}

	// apply adjacent damage
	if left != nil {
		left.Instance.Defense -= c.cardAbility.AdjacentDamage
		cardInstance := NewCardInstance(left, gameplay)
		if err := cardInstance.OnBeingAttacked(c.CardInstance); err != nil {
			return err
		}
	}

	if right != nil {
		right.Instance.Defense -= c.cardAbility.AdjacentDamage
		cardInstance := NewCardInstance(right, gameplay)
		if err := cardInstance.OnBeingAttacked(c.CardInstance); err != nil {
			return err
		}
	}

	// TODO: add action outcome

	return nil
}

package battleground

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
)

// reanimate ability
// description:
//     When zombie dies, return it to play with default Dmg, Def and effects
type reanimate struct {
	*CardInstance
	cardAbility *zb_data.CardAbilityReanimate
}

var _ Ability = &reanimate{}

func NewReanimate(card *CardInstance, cardAbility *zb_data.CardAbilityReanimate) *reanimate {
	return &reanimate{
		CardInstance: card,
		cardAbility:  cardAbility,
	}
}

func (c *reanimate) Apply(gameplay *Gameplay) error {
	owner := c.Player()
	if owner == nil {
		return fmt.Errorf("no owner for card instance %d", c.InstanceId)
	}
	if err := c.MoveZone(zb_enums.Zone_PLAY, zb_enums.Zone_GRAVEYARD); err != nil {
		return err
	}

	state := gameplay.State
	reanimate := c.cardAbility
	newInstance := proto.Clone(c.CardInstance.CardInstance).(*zb_data.CardInstance)
	// filtter out reanimate ability
	var newAbilityInstances []*zb_data.CardAbilityInstance
	for _, ability := range newInstance.AbilitiesInstances {
		switch ability.AbilityType.(type) {
		case *zb_data.CardAbilityInstance_Reanimate:
			// do not add reanimate ability
		default:
			newAbilityInstances = append(newAbilityInstances, ability)
		}
	}
	newInstance.AbilitiesInstances = newAbilityInstances
	newInstance.InstanceId.Id = state.NextInstanceId
	newInstance.Instance.Damage = reanimate.DefaultDamage
	newInstance.Instance.Defense = reanimate.DefaultDefense
	state.NextInstanceId++

	owner.CardsInGraveyard = append(owner.CardsInGraveyard, newInstance)
	newcardInstance := NewCardInstance(newInstance, gameplay)
	if err := newcardInstance.MoveZone(zb_enums.Zone_GRAVEYARD, zb_enums.Zone_PLAY); err != nil {
		return err
	}
	// just only trigger once
	reanimate.NewInstance = newInstance

	// generated outcome
	gameplay.actionOutcomes = append(gameplay.actionOutcomes, &zb_data.PlayerActionOutcome{
		Outcome: &zb_data.PlayerActionOutcome_Reanimate{
			Reanimate: &zb_data.PlayerActionOutcome_CardAbilityReanimateOutcome{
				NewCardInstance: newInstance,
			},
		},
	})

	return nil
}

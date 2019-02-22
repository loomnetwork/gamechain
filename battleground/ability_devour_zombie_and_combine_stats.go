package battleground

import (
	"fmt"

	"github.com/loomnetwork/gamechain/types/zb"
)

// devourZombieAndCombineStats ability
// description:
//     TODO
type devourZombieAndCombineStats struct {
	*CardInstance
	cardAbility *zb.CardAbilityDevourZombieAndCombineStats
	target      *zb.InstanceId
}

var _ Ability = &devourZombieAndCombineStats{}

func NewDevourZombieAndCombineStats(card *CardInstance, cardAbility *zb.CardAbilityDevourZombieAndCombineStats, target *zb.InstanceId) *devourZombieAndCombineStats {
	return &devourZombieAndCombineStats{
		CardInstance: card,
		cardAbility:  cardAbility,
		target:       target,
	}
}

func (c *devourZombieAndCombineStats) Apply(gameplay *Gameplay) error {
	cardsInPlay := gameplay.activePlayer().CardsInPlay
	_, t, found := findCardInCardListByInstanceId(c.target, cardsInPlay)
	if !found {
		return fmt.Errorf("no owner for card instance %d in play", c.target)
	}

	targetInstance := NewCardInstance(t, gameplay)
	c.Instance.Defense += targetInstance.Instance.Defense
	c.Instance.Attack += targetInstance.Instance.Attack
	if err := targetInstance.MoveZone(zb.Zone_PLAY, zb.Zone_GRAVEYARD); err != nil {
		return err
	}

	// TODO: Implement outcome for frontend

	return nil
}

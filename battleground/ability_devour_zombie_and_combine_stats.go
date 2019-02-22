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
	targets     []*zb.Unit
}

var _ Ability = &devourZombieAndCombineStats{}

func NewDevourZombieAndCombineStats(card *CardInstance, cardAbility *zb.CardAbilityDevourZombieAndCombineStats, targets []*zb.Unit) *devourZombieAndCombineStats {
	return &devourZombieAndCombineStats{
		CardInstance: card,
		cardAbility:  cardAbility,
		targets:      targets,
	}
}

func (c *devourZombieAndCombineStats) Apply(gameplay *Gameplay) error {
	cardsInPlay := gameplay.activePlayer().CardsInPlay

	for _, target := range c.targets {
		_, t, found := findCardInCardListByInstanceId(target.InstanceId, cardsInPlay)
		if !found {
			return fmt.Errorf("no owner for card instance %d in play", target.InstanceId)
		}

		targetInstance := NewCardInstance(t, gameplay)
		c.Instance.Defense += targetInstance.Instance.Defense
		c.Instance.Attack += targetInstance.Instance.Attack
		if err := targetInstance.MoveZone(zb.Zone_PLAY, zb.Zone_GRAVEYARD); err != nil {
			return err
		}
	}

	// TODO: Implement outcome for frontend

	return nil
}

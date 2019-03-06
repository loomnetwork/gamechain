package battleground

import (
	"github.com/loomnetwork/gamechain/types/zb"
)

// devourZombieAndCombineStats ability
// description:
//     TODO
type devourZombieAndCombineStats struct {
	*CardInstance
	cardAbility *zb.CardAbilityDevourZombieAndCombineStats
	targets     []*CardInstance
}

var _ Ability = &devourZombieAndCombineStats{}

func NewDevourZombieAndCombineStats(card *CardInstance, cardAbility *zb.CardAbilityDevourZombieAndCombineStats, targets []*CardInstance) *devourZombieAndCombineStats {
	return &devourZombieAndCombineStats{
		CardInstance: card,
		cardAbility:  cardAbility,
		targets:      targets,
	}
}

func (c *devourZombieAndCombineStats) Apply(gameplay *Gameplay) error {
	for _, target := range c.targets {
		c.Instance.Defense += target.Instance.Defense
		c.Instance.Attack += target.Instance.Attack
		if err := target.MoveZone(zb.Zone_PLAY, zb.Zone_GRAVEYARD); err != nil {
			return err
		}
	}

	// TODO: Implement outcome for frontend

	return nil
}

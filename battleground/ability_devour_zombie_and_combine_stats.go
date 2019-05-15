package battleground

import (
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
)

// devourZombieAndCombineStats ability
// description:
//     TODO
type devourZombieAndCombineStats struct {
	*CardInstance
	cardAbility *zb_data.CardAbilityDevourZombieAndCombineStats
	targets     []*CardInstance
}

var _ Ability = &devourZombieAndCombineStats{}

func NewDevourZombieAndCombineStats(card *CardInstance, cardAbility *zb_data.CardAbilityDevourZombieAndCombineStats, targets []*CardInstance) *devourZombieAndCombineStats {
	return &devourZombieAndCombineStats{
		CardInstance: card,
		cardAbility:  cardAbility,
		targets:      targets,
	}
}

func (c *devourZombieAndCombineStats) Apply(gameplay *Gameplay) error {
	for _, target := range c.targets {
		c.Instance.Defense += target.Instance.Defense
		c.Instance.Damage += target.Instance.Damage
		if err := target.MoveZone(zb_enums.Zone_PLAY, zb_enums.Zone_GRAVEYARD); err != nil {
			return err
		}
	}

	// TODO: Implement outcome for frontend

	return nil
}

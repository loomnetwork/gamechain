package battleground

import (
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
)

// changeStat ability
// description:
//     TODO
type changeStat struct {
	*CardInstance
	cardAbility *zb_data.CardAbilityChangeStat
	target      *zb_data.InstanceId
}

var _ Ability = &changeStat{}

func NewChangeState(card *CardInstance, cardAbility *zb_data.CardAbilityChangeStat, target *zb_data.InstanceId) *changeStat {
	return &changeStat{
		CardInstance: card,
		cardAbility:  cardAbility,
		target:       target,
	}
}

func (c *changeStat) Apply(gameplay *Gameplay) error {
	changeStat := c.cardAbility
	// Once attacking, defense and damage values are decreased
	// TODO: generate change zone first
	if changeStat.Stat == zb_enums.Stat_Defense {
		c.Instance.Defense += changeStat.StatAdjustment
		// generate outcome
		gameplay.actionOutcomes = append(gameplay.actionOutcomes, &zb_data.PlayerActionOutcome{
			Outcome: &zb_data.PlayerActionOutcome_ChangeStat{
				ChangeStat: &zb_data.PlayerActionOutcome_CardAbilityChangeStatOutcome{
					InstanceId:       c.InstanceId,
					NewDefense:       c.Instance.Defense,
					Stat:             zb_enums.Stat_Defense,
					TargetInstanceId: c.target,
				},
			},
		})
	} else if changeStat.Stat == zb_enums.Stat_Damage {
		c.Instance.Damage += changeStat.StatAdjustment
		// generate outcome
		gameplay.actionOutcomes = append(gameplay.actionOutcomes, &zb_data.PlayerActionOutcome{
			Outcome: &zb_data.PlayerActionOutcome_ChangeStat{
				ChangeStat: &zb_data.PlayerActionOutcome_CardAbilityChangeStatOutcome{
					InstanceId:       c.InstanceId,
					NewDamage:        c.Instance.Damage,
					Stat:             zb_enums.Stat_Damage,
					TargetInstanceId: c.target,
				},
			},
		})
	}
	return nil
}

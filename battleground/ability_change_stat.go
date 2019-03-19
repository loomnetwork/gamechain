package battleground

import "github.com/loomnetwork/gamechain/types/zb"

// changeStat ability
// description:
//     TODO
type changeStat struct {
	*CardInstance
	cardAbility *zb.CardAbilityChangeStat
	target      *zb.InstanceId
}

var _ Ability = &changeStat{}

func NewChangeState(card *CardInstance, cardAbility *zb.CardAbilityChangeStat, target *zb.InstanceId) *changeStat {
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
	if changeStat.Stat == zb.StatType_Defense {
		c.Instance.Defense += changeStat.StatAdjustment
		// generate outcome
		gameplay.actionOutcomes = append(gameplay.actionOutcomes, &zb.PlayerActionOutcome{
			Outcome: &zb.PlayerActionOutcome_ChangeStat{
				ChangeStat: &zb.PlayerActionOutcome_CardAbilityChangeStatOutcome{
					InstanceId:       c.InstanceId,
					NewDefense:       c.Instance.Defense,
					Stat:             zb.StatType_Defense,
					TargetInstanceId: c.target,
				},
			},
		})
	} else if changeStat.Stat == zb.StatType_Damage {
		c.Instance.Damage += changeStat.StatAdjustment
		// generate outcome
		gameplay.actionOutcomes = append(gameplay.actionOutcomes, &zb.PlayerActionOutcome{
			Outcome: &zb.PlayerActionOutcome_ChangeStat{
				ChangeStat: &zb.PlayerActionOutcome_CardAbilityChangeStatOutcome{
					InstanceId:       c.InstanceId,
					NewDamage:        c.Instance.Damage,
					Stat:             zb.StatType_Damage,
					TargetInstanceId: c.target,
				},
			},
		})
	}
	return nil
}

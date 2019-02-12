package battleground

import (
	"github.com/loomnetwork/gamechain/types/zb"
)

type CardInstance struct {
	*zb.CardInstance
	Outcomes []*zb.PlayerActionOutcome
	Gameplay *Gameplay
}

func NewCardInstance(instance *zb.CardInstance, gameplay *Gameplay) *CardInstance {
	return &CardInstance{
		CardInstance: instance,
		Gameplay:     gameplay,
	}
}

func (c *CardInstance) Attack(target *CardInstance) {
	old := c.Instance.Defense
	c.Instance.Defense = c.Instance.Defense - target.Instance.Attack
	c.OnDefenseChange(old, c.Instance.Defense)

	old = target.Instance.Defense
	target.Instance.Defense = target.Instance.Defense - c.Instance.Attack
	target.OnDefenseChange(old, target.Instance.Defense)

	if c.Instance.Defense <= 0 {
		c.OnDeath(target)
	}

	if target.Instance.Defense <= 0 {
		target.OnDeath(c)
	}
}

func (c *CardInstance) GotDamage(attacker *CardInstance) {
}

func (c *CardInstance) OnDeath(attacker *CardInstance) {
	// trigger target ability
	for _, ai := range c.AbilitiesInstances {
		if ai.Trigger == zb.CardAbilityTrigger_Death {
			switch ability := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_Reanimate:
				reanimate := ability.Reanimate
				// When zombie dies, return it to play with default Atk, Def and effects
				// TODO: generate change zone first
				c.Instance.Attack = reanimate.Attack
				c.Instance.Defense = reanimate.Defence
				ai.IsActive = false
			}
		}
	}

	// trigger attacker ability
	for _, ai := range attacker.AbilitiesInstances {
		if ai.Trigger == zb.CardAbilityTrigger_Attack {
			// activate ability
			switch ability := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_PriorityAttack:
				priorityAttack := ability.PriorityAttack
				// reset the card's defense to the value before the attack, only if the opponent card dies
				attacker.Instance.Defense = priorityAttack.AttackerOldDefense
			}
		}
	}
}

func (c *CardInstance) OnDefenseChange(oldValue, newValue int32) {
	// trigger ability
	for _, ai := range c.AbilitiesInstances {
		if ai.Trigger == zb.CardAbilityTrigger_GotDamage {
			// activate ability
			switch ability := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_Rage:
				rage := ability.Rage
				if !rage.WasApplied {
					if c.Instance.Defense < c.Prototype.Defense {
						// change state
						rage.WasApplied = true
						c.Instance.Attack += rage.AddedAttack
						// generate outcome
						c.Gameplay.actionOutcomes = append(c.Gameplay.actionOutcomes, &zb.PlayerActionOutcome{
							Outcome: &zb.PlayerActionOutcome_Rage{
								Rage: &zb.PlayerActionOutcome_CardAbilityRageOutcome{
									InstanceId: c.InstanceId,
									NewAttack:  c.Instance.Attack,
								},
							},
						})
					}
				} else if c.Instance.Defense >= c.Prototype.Defense {
					rage.WasApplied = false
					c.Instance.Attack -= rage.AddedAttack
					c.Gameplay.actionOutcomes = append(c.Gameplay.actionOutcomes, &zb.PlayerActionOutcome{
						Outcome: &zb.PlayerActionOutcome_Rage{
							Rage: &zb.PlayerActionOutcome_CardAbilityRageOutcome{
								InstanceId: c.InstanceId,
								NewAttack:  c.Instance.Attack,
							},
						},
					})
				}
			}
		}
	}
}

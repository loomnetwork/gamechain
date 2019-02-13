package battleground

import (
	"github.com/gogo/protobuf/proto"
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

func (c *CardInstance) Play() {
	c.OnPlay()
}

func (c *CardInstance) Attack(target *CardInstance) {
	old := c.Instance.Defense
	c.Instance.Defense = c.Instance.Defense - target.Instance.Attack
	c.OnDefenseChange(old, c.Instance.Defense)

	old = target.Instance.Defense
	target.Instance.Defense = target.Instance.Defense - c.Instance.Attack
	target.OnBeingAttacked(c)
	target.OnDefenseChange(old, target.Instance.Defense)
	c.AfterAttacking(target)

	if c.Instance.Defense <= 0 {
		c.OnDeath(target)
	}

	if target.Instance.Defense <= 0 {
		target.OnDeath(c)
	}
}

func (c *CardInstance) GotDamage(attacker *CardInstance) {
}

func (c *CardInstance) AfterAttacking(target *CardInstance) {
	for _, ai := range c.AbilitiesInstances {
		if ai.Trigger == zb.CardAbilityTrigger_Attack {
			switch ability := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_ChangeStat:
				changeStat := ability.ChangeStat
				// Once attacking, defense and attack values are decreased
				// TODO: generate change zone first
				c.Instance.Defense -= changeStat.DecreasedValue
				c.Instance.Attack -= changeStat.DecreasedValue
			}
		}
	}
}

func (c *CardInstance) OnBeingAttacked(attacker *CardInstance) {
	for _, ai := range attacker.AbilitiesInstances {
		if ai.Trigger == zb.CardAbilityTrigger_Attack {
			switch ability := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_AdditionalDamageToHeavyInAttack:
				additionalDamageToHeavyInAttack := ability.AdditionalDamageToHeavyInAttack
				// If the target is heavy, add addtional attack to defense
				// TODO: generate change zone first
				if c.Instance.Type == zb.CreatureType_Heavy {
					c.Instance.Defense -= additionalDamageToHeavyInAttack.AddedAttack
				}
			}
		}
	}
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

	// after apply ability, update zone if the card instance is really dead
	if c.Instance.Defense <= 0 {
		c.MoveZone(zb.Zone_PLAY, zb.Zone_GRAVEYARD)
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

func (c *CardInstance) OnPlay() {
	c.MoveZone(zb.Zone_HAND, zb.Zone_PLAY)
}

func (c *CardInstance) MoveZone(from, to zb.ZoneType) {
	var cardInstance *zb.CardInstance
	var cardIndex int
	var owner *zb.PlayerState
	for i := 0; i < len(c.Gameplay.State.PlayerStates); i++ {
		for j, card := range c.Gameplay.State.PlayerStates[i].CardsInPlay {
			if proto.Equal(card.InstanceId, c.InstanceId) {
				cardInstance = card
				cardIndex = j
				owner = c.Gameplay.State.PlayerStates[i]
				break
			}
		}
	}

	if cardInstance != nil {
		if from == zb.Zone_PLAY && to == zb.Zone_GRAVEYARD {
			// move from play to graveyard
			owner.CardsInPlay = append(owner.CardsInPlay[:cardIndex], owner.CardsInPlay[cardIndex+1:]...)
			owner.CardsInGraveyard = append(owner.CardsInGraveyard, cardInstance)
		}

		// FIX ME
		// else if from == zb.Zone_HAND && to == zb.Zone_PLAY {
		// 	// move from hand to play
		// 	owner.CardsInPlay = append(owner.CardsInPlay, cardInstance)
		// 	activeCardsInHand := owner.CardsInHand
		// 	activeCardsInHand = append(activeCardsInHand[:cardIndex], activeCardsInHand[cardIndex+1:]...)
		// 	owner.CardsInHand = activeCardsInHand
		// }

	}
}

func (c *CardInstance) AttackOverload(target *zb.PlayerState, attacker *zb.PlayerState) {
	target.Defense -= c.Instance.Attack

	if target.Defense <= 0 {
		c.Gameplay.State.Winner = attacker.Id
		c.Gameplay.State.IsEnded = true
	}
}

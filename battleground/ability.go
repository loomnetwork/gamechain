package battleground

import (
	"fmt"

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

	if c.Instance.Defense <= 0 {
		c.OnDeath(target)
	}

	if target.Instance.Defense <= 0 {
		target.OnDeath(c)
	}
}

func (c *CardInstance) GotDamage(attacker *CardInstance) {
}

func (c *CardInstance) OnBeingAttacked(attacker *CardInstance) {
	for _, ai := range attacker.AbilitiesInstances {
		if ai.Trigger == zb.CardAbilityTrigger_Attack {
			switch ability := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_AdditionalDamageToHeavyInAttack:
				additionalDamageToHeavyInAttack := ability.AdditionalDamageToHeavyInAttack
				// When zombie dies, return it to play with default Atk, Def and effects
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
		if !ai.IsActive {
			continue
		}
		if ai.Trigger == zb.CardAbilityTrigger_Death {
			switch ability := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_Reanimate:
				reanimate := ability.Reanimate
				// When zombie dies, return it to play with default Atk, Def and effects
				// TODO: generate change zone first
				c.Instance.Attack = reanimate.Attack
				c.Instance.Defense = reanimate.Defense
				// just only trigger once
				ai.IsActive = false

				// generated outcome
				c.Gameplay.actionOutcomes = append(c.Gameplay.actionOutcomes, &zb.PlayerActionOutcome{
					Outcome: &zb.PlayerActionOutcome_Reanimate{
						Reanimate: &zb.PlayerActionOutcome_CardAbilityReanimateOutcome{
							InstanceId: c.InstanceId,
							Attack:     c.Instance.Attack,
							Defense:    c.Instance.Defense,
						},
					},
				})
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

func (c *CardInstance) OnPlay() error {
	return c.MoveZone(zb.Zone_HAND, zb.Zone_PLAY)
}

func (c *CardInstance) MoveZone(from, to zb.ZoneType) error {
	var cardInstance *zb.CardInstance
	var cardIndex int
	var owner *zb.PlayerState
	// move from play to graveyard
	if from == zb.Zone_PLAY && to == zb.Zone_GRAVEYARD {
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
		if cardInstance == nil {
			return fmt.Errorf("card instance id %s not found in play", c.InstanceId)
		}
		owner.CardsInPlay = append(owner.CardsInPlay[:cardIndex], owner.CardsInPlay[cardIndex+1:]...)
		owner.CardsInGraveyard = append(owner.CardsInGraveyard, cardInstance)
		c.Zone = zb.Zone_GRAVEYARD
	}
	// move from hand to play
	if from == zb.Zone_HAND && to == zb.Zone_PLAY {
		for i := 0; i < len(c.Gameplay.State.PlayerStates); i++ {
			for j, card := range c.Gameplay.State.PlayerStates[i].CardsInHand {
				if proto.Equal(card.InstanceId, c.InstanceId) {
					cardInstance = card
					cardIndex = j
					owner = c.Gameplay.State.PlayerStates[i]
					break
				}
			}
		}
		if cardInstance == nil {
			return fmt.Errorf("card instance id %s not found in play", c.InstanceId)
		}
		owner.CardsInHand = append(owner.CardsInHand[:cardIndex], owner.CardsInHand[cardIndex+1:]...)
		owner.CardsInPlay = append(owner.CardsInPlay, cardInstance)
		c.Zone = zb.Zone_PLAY
	}
	// move from hand to play
	if from == zb.Zone_HAND && to == zb.Zone_DECK {
		for i := 0; i < len(c.Gameplay.State.PlayerStates); i++ {
			for j, card := range c.Gameplay.State.PlayerStates[i].CardsInHand {
				if proto.Equal(card.InstanceId, c.InstanceId) {
					cardInstance = card
					cardIndex = j
					owner = c.Gameplay.State.PlayerStates[i]
					break
				}
			}
		}
		if cardInstance == nil {
			return fmt.Errorf("card instance id %s not found in play", c.InstanceId)
		}
		owner.CardsInHand = append(owner.CardsInHand[:cardIndex], owner.CardsInHand[cardIndex+1:]...)
		owner.CardsInDeck = append(owner.CardsInDeck, cardInstance)
		c.Zone = zb.Zone_DECK
	}
	// move from deck to hand
	if from == zb.Zone_DECK && to == zb.Zone_HAND {
		for i := 0; i < len(c.Gameplay.State.PlayerStates); i++ {
			for j, card := range c.Gameplay.State.PlayerStates[i].CardsInDeck {
				if proto.Equal(card.InstanceId, c.InstanceId) {
					cardInstance = card
					cardIndex = j
					owner = c.Gameplay.State.PlayerStates[i]
					break
				}
			}
		}
		if cardInstance == nil {
			return fmt.Errorf("card instance id %s not found in play", c.InstanceId)
		}
		owner.CardsInDeck = append(owner.CardsInDeck[:cardIndex], owner.CardsInDeck[cardIndex+1:]...)
		owner.CardsInHand = append(owner.CardsInHand, cardInstance)
		c.Zone = zb.Zone_HAND
	}

	return nil
}

func (c *CardInstance) AttackOverload(target *zb.PlayerState, attacker *zb.PlayerState) {
	target.Defense -= c.Instance.Attack

	if target.Defense <= 0 {
		c.Gameplay.State.Winner = attacker.Id
		c.Gameplay.State.IsEnded = true
	}
}

func (c *CardInstance) Mulligan() error {

	var owner *zb.PlayerState
	for i := 0; i < len(c.Gameplay.State.PlayerStates); i++ {
		for _, card := range c.Gameplay.State.PlayerStates[i].CardsInHand {
			if proto.Equal(card.InstanceId, c.InstanceId) {
				owner = c.Gameplay.State.PlayerStates[i]
				break
			}
		}
	}

	if owner == nil {
		return fmt.Errorf("no owner for card instance %d", c.InstanceId)
	}
	// also draw a new card from deck to hand
	if len(owner.CardsInDeck) == 0 {
		return fmt.Errorf("no card in deck to be drawn")
	}

	if err := c.MoveZone(zb.Zone_HAND, zb.Zone_DECK); err != nil {
		return err
	}

	owner.MulliganCards = append(owner.MulliganCards, c.CardInstance)

	newcard := owner.CardsInDeck[0]
	newCardInstance := NewCardInstance(newcard, c.Gameplay)
	newCardInstance.MoveZone(zb.Zone_DECK, zb.Zone_HAND)
	return nil
}

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
	c.AfterAttacking()

	if c.Instance.Defense <= 0 {
		c.OnDeath(target)
	}

	if target.Instance.Defense <= 0 {
		target.OnDeath(c)
	}
}

func (c *CardInstance) GotDamage(attacker *CardInstance) {
}

func (c *CardInstance) AfterAttacking() {
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
		if !ai.IsActive {
			continue
		}
		if ai.Trigger == zb.CardAbilityTrigger_Death {
			switch ability := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_Reanimate:
				// When zombie dies, return it to play with default Atk, Def and effects

				// find the new instance id
				var nextInstanceID int32
				for _, playerState := range c.Gameplay.State.PlayerStates {
					for _, card := range playerState.CardsInPlay {
						if card.InstanceId.Id > nextInstanceID {
							nextInstanceID = card.InstanceId.Id
						}
					}
					for _, card := range playerState.CardsInHand {
						if card.InstanceId.Id > nextInstanceID {
							nextInstanceID = card.InstanceId.Id
						}
					}
					for _, card := range playerState.CardsInDeck {
						if card.InstanceId.Id > nextInstanceID {
							nextInstanceID = card.InstanceId.Id
						}
					}
					for _, card := range playerState.CardsInGraveyard {
						if card.InstanceId.Id > nextInstanceID {
							nextInstanceID = card.InstanceId.Id
						}
					}
				}
				nextInstanceID++
				c.MoveZone(zb.Zone_PLAY, zb.Zone_GRAVEYARD)
				ai.IsActive = false
				reanimate := ability.Reanimate
				newInstance := proto.Clone(c.CardInstance).(*zb.CardInstance)
				// remove ability
				var newAbilities []*zb.CardAbilityInstance
				for _, ability := range newInstance.AbilitiesInstances {
					if ability.AbilityType != ai.AbilityType {
						newAbilities = append(newAbilities, ability)
					}
				}
				newInstance.AbilitiesInstances = newAbilities
				newInstance.InstanceId.Id = nextInstanceID
				newInstance.Instance.Attack = reanimate.DefaultAttack
				newInstance.Instance.Defense = reanimate.DefaultDefense
				// FIX ME: better way to do this?
				var activePlayer *zb.PlayerState
				for _, playerState := range c.Gameplay.State.PlayerStates {
					if playerState.Id == newInstance.Owner {
						activePlayer = playerState
						break
					}
				}
				if activePlayer == nil {
					panic("want not nil activePlayer")
				}
				activePlayer.CardsInGraveyard = append(activePlayer.CardsInGraveyard, newInstance)
				newcardInstance := NewCardInstance(newInstance, c.Gameplay)
				newcardInstance.MoveZone(zb.Zone_GRAVEYARD, zb.Zone_PLAY)
				// just only trigger once
				reanimate.NewInstance = newInstance

				// generated outcome
				c.Gameplay.actionOutcomes = append(c.Gameplay.actionOutcomes, &zb.PlayerActionOutcome{
					Outcome: &zb.PlayerActionOutcome_Reanimate{
						Reanimate: &zb.PlayerActionOutcome_CardAbilityReanimateOutcome{
							NewCardInstance: newInstance,
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
		if cardIndex == 0 {
			owner.CardsInPlay = owner.CardsInPlay[cardIndex+1:]
		} else {
			owner.CardsInPlay = append(owner.CardsInPlay[:cardIndex], owner.CardsInPlay[cardIndex+1:]...)
		}
		owner.CardsInGraveyard = append(owner.CardsInGraveyard, cardInstance)
		c.Zone = zb.Zone_GRAVEYARD
	} else if from == zb.Zone_GRAVEYARD && to == zb.Zone_PLAY {
		for i := 0; i < len(c.Gameplay.State.PlayerStates); i++ {
			for j, card := range c.Gameplay.State.PlayerStates[i].CardsInGraveyard {
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
		if cardIndex == 0 {
			owner.CardsInGraveyard = owner.CardsInGraveyard[cardIndex+1:]
		} else {
			owner.CardsInGraveyard = append(owner.CardsInGraveyard[:cardIndex], owner.CardsInGraveyard[cardIndex+1:]...)
		}
		owner.CardsInPlay = append(owner.CardsInPlay, cardInstance)
		c.Zone = zb.Zone_PLAY
	} else if from == zb.Zone_HAND && to == zb.Zone_PLAY {
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
		if cardIndex == 0 {
			owner.CardsInHand = owner.CardsInHand[cardIndex+1:]
		} else {
			owner.CardsInHand = append(owner.CardsInHand[:cardIndex], owner.CardsInHand[cardIndex+1:]...)
		}
		owner.CardsInPlay = append(owner.CardsInPlay, cardInstance)
		c.Zone = zb.Zone_PLAY
	} else if from == zb.Zone_HAND && to == zb.Zone_DECK {
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
		if cardIndex == 0 {
			owner.CardsInHand = owner.CardsInHand[cardIndex+1:]
		} else {
			owner.CardsInHand = append(owner.CardsInHand[:cardIndex], owner.CardsInHand[cardIndex+1:]...)
		}
		owner.CardsInDeck = append(owner.CardsInDeck, cardInstance)
		c.Zone = zb.Zone_DECK
	} else if from == zb.Zone_DECK && to == zb.Zone_HAND {
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
		if cardIndex == 0 {
			owner.CardsInDeck = owner.CardsInDeck[cardIndex+1:]
		} else {
			owner.CardsInDeck = append(owner.CardsInDeck[:cardIndex], owner.CardsInDeck[cardIndex+1:]...)
		}
		owner.CardsInHand = append(owner.CardsInHand, cardInstance)
		c.Zone = zb.Zone_HAND
	} else {
		return fmt.Errorf("invalid move zone from %v to %v", from, to)
	}

	return nil
}

func (c *CardInstance) AttackOverlord(target *zb.PlayerState, attacker *zb.PlayerState) {
	c.Gameplay.debugf("Attack Overlord")
	target.Defense -= c.Instance.Attack

	if target.Defense <= 0 {
		c.Gameplay.State.Winner = attacker.Id
		c.Gameplay.State.IsEnded = true
	}
	c.AfterAttacking()
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

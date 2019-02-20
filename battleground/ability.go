package battleground

import (
	"fmt"
	"math/rand"

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
	c.AfterAttacking(target.InstanceId)

	if c.Instance.Defense <= 0 {
		c.OnDeath(target)
	}

	if target.Instance.Defense <= 0 {
		target.OnDeath(c)
	}
}

func (c *CardInstance) GotDamage(attacker *CardInstance) {
}

func (c *CardInstance) AfterAttacking(targetInstanceId *zb.InstanceId) {
	for _, ai := range c.AbilitiesInstances {
		if ai.Trigger == zb.CardAbilityTrigger_Attack {
			switch ability := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_ChangeStat:
				changeStat := ability.ChangeStat
				// Once attacking, defense and attack values are decreased
				// TODO: generate change zone first
				if changeStat.Stat == zb.StatType_Defense {
					c.Instance.Defense += changeStat.StatAdjustment
					// generated outcome
					c.Gameplay.actionOutcomes = append(c.Gameplay.actionOutcomes, &zb.PlayerActionOutcome{
						Outcome: &zb.PlayerActionOutcome_ChangeStat{
							ChangeStat: &zb.PlayerActionOutcome_CardAbilityChangeStatOutcome{
								InstanceId:       c.InstanceId,
								NewDefense:       c.Instance.Defense,
								Stat:             zb.StatType_Defense,
								TargetInstanceId: targetInstanceId,
							},
						},
					})
				} else if changeStat.Stat == zb.StatType_Attack {
					c.Instance.Attack += changeStat.StatAdjustment
					// generated outcome
					c.Gameplay.actionOutcomes = append(c.Gameplay.actionOutcomes, &zb.PlayerActionOutcome{
						Outcome: &zb.PlayerActionOutcome_ChangeStat{
							ChangeStat: &zb.PlayerActionOutcome_CardAbilityChangeStatOutcome{
								InstanceId:       c.InstanceId,
								NewAttack:        c.Instance.Attack,
								Stat:             zb.StatType_Attack,
								TargetInstanceId: targetInstanceId,
							},
						},
					})
				}

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
				newInstance.InstanceId.Id = c.Gameplay.State.NextInstanceId
				newInstance.Instance.Attack = reanimate.DefaultAttack
				newInstance.Instance.Defense = reanimate.DefaultDefense
				c.Gameplay.State.NextInstanceId++
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
	// trigger card ability on play
	for _, ai := range c.AbilitiesInstances {
		if ai.Trigger == zb.CardAbilityTrigger_Entry {
			// activate ability
			switch ability := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_AttackOverlord:
				attackOverlord := ability.AttackOverlord
				if !attackOverlord.WasApplied {
					// damage player overlord
					c.Gameplay.activePlayer().Defense -= attackOverlord.Damage
					attackOverlord.WasApplied = true
					c.Gameplay.actionOutcomes = append(c.Gameplay.actionOutcomes, &zb.PlayerActionOutcome{
						Outcome: &zb.PlayerActionOutcome_AttackOverlord{
							AttackOverlord: &zb.PlayerActionOutcome_CardAbilityAttackOverlordOutcome{
								InstanceId: c.Gameplay.activePlayer().InstanceId,
								NewDefense: c.Gameplay.activePlayer().Defense,
							},
						},
					})
				}
			case *zb.CardAbilityInstance_ReplaceUnitsWithTypeOnStrongerOnes:
				owner := c.owner()
				if owner == nil {
					return fmt.Errorf("no owner for card instance %d", c.InstanceId)
				}
				// find the cards in card library with same types as cards in plays
				var toReplaceCards []*zb.CardInstance
				var oldInstanceIds []*zb.InstanceId
				for _, card := range owner.CardsInPlay {
					if c.Instance.Set == card.Instance.Set && !proto.Equal(c.InstanceId, card.InstanceId) {
						toReplaceCards = append(toReplaceCards, card)
						oldInstanceIds = append(oldInstanceIds, card.InstanceId)
					}
				}

				// continue if there is no same type card in play
				if len(toReplaceCards) == 0 {
					continue
				}

				sameTypeStrongerFn := func(cardLibrary *zb.CardList, target *zb.CardInstance) []*zb.Card {
					var sameTypeStrongerCards []*zb.Card
					for _, card := range cardLibrary.Cards {
						if card.Set == target.Instance.Set && card.GooCost > target.Instance.GooCost {
							sameTypeStrongerCards = append(sameTypeStrongerCards, card)
						}
					}
					return sameTypeStrongerCards
				}

				var newcardInstances []*zb.PlayerActionOutcome_CardAbilityReplaceUnitsWithTypeOnStrongerOnes_NewCardInstance
				for i, card := range toReplaceCards {
					sameTypeStrongerCards := sameTypeStrongerFn(c.Gameplay.cardLibrary, card)
					if len(sameTypeStrongerCards) == 0 {
						continue
					}
					var r = rand.New(rand.NewSource(c.Gameplay.State.RandomSeed))
					randomCardIndex := r.Perm(len(sameTypeStrongerCards))

					// create new instance from card
					newcard := sameTypeStrongerCards[randomCardIndex[i]]
					instanceid := &zb.InstanceId{Id: c.Gameplay.State.NextInstanceId}
					c.Gameplay.State.NextInstanceId++
					newinstance := newCardInstanceFromCardDetails(newcard, instanceid, c.Owner, c.OwnerIndex)
					newinstance.Zone = zb.Zone_PLAY
					newcardInstances = append(newcardInstances, &zb.PlayerActionOutcome_CardAbilityReplaceUnitsWithTypeOnStrongerOnes_NewCardInstance{
						CardInstance: newinstance,
						Position:     int32(i),
					})
				}

				// remove card from card in play
				var newCardsInplay []*zb.CardInstance
				for _, card := range owner.CardsInPlay {
					for _, toreplace := range toReplaceCards {
						if !proto.Equal(toreplace.InstanceId, card.InstanceId) {
							newCardsInplay = append(newCardsInplay, card)
						}
					}
				}

				// append card in play
				// TODO: maybe we don't append, we just replace?
				for _, card := range newcardInstances {
					newCardsInplay = append(newCardsInplay, card.CardInstance)
				}
				// set cardinplay to gamestate
				c.owner().CardsInPlay = newCardsInplay

				ai.IsActive = false
				// outcome
				c.Gameplay.actionOutcomes = append(c.Gameplay.actionOutcomes, &zb.PlayerActionOutcome{
					Outcome: &zb.PlayerActionOutcome_ReplaceUnitsWithTypeOnStrongerOnes{
						ReplaceUnitsWithTypeOnStrongerOnes: &zb.PlayerActionOutcome_CardAbilityReplaceUnitsWithTypeOnStrongerOnes{
							NewCardInstances: newcardInstances,
							OldInstanceIds:   oldInstanceIds,
						},
					},
				})
			}
		}
	}

	return c.MoveZone(zb.Zone_HAND, zb.Zone_PLAY)
}

func (c *CardInstance) MoveZone(from, to zb.ZoneType) error {
	if int(c.OwnerIndex) > len(c.Gameplay.State.PlayerStates)-1 {
		return fmt.Errorf("Invalid owner index: %d", c.OwnerIndex)
	}

	owner := c.Gameplay.State.PlayerStates[c.OwnerIndex]
	var err error

	switch {
	case from == zb.Zone_PLAY && to == zb.Zone_GRAVEYARD:
		owner.CardsInPlay, owner.CardsInGraveyard, err = moveCard(c, owner.CardsInPlay, owner.CardsInGraveyard, zb.Zone_GRAVEYARD)
	case from == zb.Zone_GRAVEYARD && to == zb.Zone_PLAY:
		owner.CardsInGraveyard, owner.CardsInPlay, err = moveCard(c, owner.CardsInGraveyard, owner.CardsInPlay, zb.Zone_PLAY)
	case from == zb.Zone_HAND && to == zb.Zone_PLAY:
		owner.CardsInHand, owner.CardsInPlay, err = moveCard(c, owner.CardsInHand, owner.CardsInPlay, zb.Zone_PLAY)
	case from == zb.Zone_HAND && to == zb.Zone_DECK:
		owner.CardsInHand, owner.CardsInDeck, err = moveCard(c, owner.CardsInHand, owner.CardsInDeck, zb.Zone_DECK)
	case from == zb.Zone_DECK && to == zb.Zone_HAND:
		owner.CardsInDeck, owner.CardsInHand, err = moveCard(c, owner.CardsInDeck, owner.CardsInHand, zb.Zone_HAND)
	default:
	}

	return err
}

func moveCard(c *CardInstance, from, to []*zb.CardInstance, zone zb.ZoneType) ([]*zb.CardInstance, []*zb.CardInstance, error) {
	var cardInstance *zb.CardInstance
	var cardIndex int
	for i, card := range from {
		if proto.Equal(card.InstanceId, c.InstanceId) {
			cardInstance = card
			cardIndex = i
			break
		}
	}

	if cardInstance == nil {
		return from, to, fmt.Errorf("card instance id %s not found in play", c.InstanceId)
	}
	if cardIndex == 0 {
		from = from[cardIndex+1:]
	} else {
		from = append(from[:cardIndex], from[cardIndex+1:]...)
	}
	to = append(to, cardInstance)
	c.Zone = zone
	return from, to, nil
}

func (c *CardInstance) AttackOverlord(target *zb.PlayerState, attacker *zb.PlayerState) {
	target.Defense -= c.Instance.Attack
	c.Gameplay.actionOutcomes = append(c.Gameplay.actionOutcomes, &zb.PlayerActionOutcome{
		Outcome: &zb.PlayerActionOutcome_CardAttack{
			CardAttack: &zb.PlayerActionOutcome_CardAttackOutcome{
				AttackerInstanceId: c.InstanceId,
				TargetInstanceId:   target.InstanceId,
				AttackerNewDefense: c.Instance.Defense,
				TargetNewDefense:   target.Defense,
			},
		},
	})

	if target.Defense <= 0 {
		c.Gameplay.State.Winner = attacker.Id
		c.Gameplay.State.IsEnded = true
	}
	c.AfterAttacking(target.InstanceId)
}

func (c *CardInstance) Mulligan() error {
	owner := c.owner()
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

func (c *CardInstance) owner() *zb.PlayerState {
	if int(c.OwnerIndex) > len(c.Gameplay.State.PlayerStates)-1 {
		return nil
	}
	return c.Gameplay.State.PlayerStates[c.OwnerIndex]
}

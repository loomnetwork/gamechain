package battleground

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
)

type Ability interface {
	Apply(gameplay *Gameplay) error
}

type CardInstance struct {
	*zb.CardInstance
	Gameplay *Gameplay
}

func NewCardInstance(instance *zb.CardInstance, gameplay *Gameplay) *CardInstance {
	return &CardInstance{
		CardInstance: instance,
		Gameplay:     gameplay,
	}
}

func (c *CardInstance) Play() error {
	return c.OnPlay()
}

func (c *CardInstance) UseAbility(targets []*CardInstance) error {
	return c.OnAbilityUsed(targets)
}

func (c *CardInstance) OnAbilityUsed(targets []*CardInstance) error {
	var ability Ability
	for _, ai := range c.AbilitiesInstances {
		if ai.Trigger == zb.CardAbilityTrigger_Entry {
			switch abilityInstance := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_DevourZombieAndCombineStats:
				if ai.IsActive {
					devourZombieAndCombineStats := abilityInstance.DevourZombieAndCombineStats
					ability = NewDevourZombieAndCombineStats(c, devourZombieAndCombineStats, targets)
					if err := ability.Apply(c.Gameplay); err != nil {
						return err
					}
					ai.IsActive = false
				}
			}
		}
	}
	return nil
}

func (c *CardInstance) Attack(target *CardInstance) error {
	c.Instance.Defense = c.Instance.Defense - target.Instance.Damage
	target.Instance.Defense = target.Instance.Defense - c.Instance.Damage

	if err := c.OnAttack(target); err != nil {
		return err
	}

	if err := target.OnBeingAttacked(c); err != nil {
		return err
	}

	if err := c.AfterAttacking(target.InstanceId); err != nil {
		return err
	}

	if c.Instance.Defense <= 0 {
		if err := c.OnDeath(target); err != nil {
			return err
		}
	}

	return nil
}

// OnAttack trigger ability when the card attacks a target
func (c *CardInstance) OnAttack(target *CardInstance) error {
	var ability Ability
	for _, ai := range c.AbilitiesInstances {
		if ai.Trigger == zb.CardAbilityTrigger_Attack {
			switch abilityInstance := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_AdditionalDamageToHeavyInAttack:
				additionalDamageToHeavyInAttack := abilityInstance.AdditionalDamageToHeavyInAttack
				ab := NewAdditionalDamgeToHeavyInAttack(c, additionalDamageToHeavyInAttack, target)
				if err := ab.Apply(c.Gameplay); err != nil {
					return err
				}
			case *zb.CardAbilityInstance_DealDamageToThisAndAdjacentUnits:
				dealDamageToThisAndAdjacentUnits := abilityInstance.DealDamageToThisAndAdjacentUnits
				ability = NewDealDamageToThisAndAdjacentUnits(c, dealDamageToThisAndAdjacentUnits, target)
				if err := ability.Apply(c.Gameplay); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (c *CardInstance) AfterAttacking(target *zb.InstanceId) error {
	for _, ai := range c.AbilitiesInstances {
		if ai.Trigger == zb.CardAbilityTrigger_Attack {
			switch ability := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_ChangeStat:
				changeStat := ability.ChangeStat
				ab := NewChangeState(c, changeStat, target)
				if err := ab.Apply(c.Gameplay); err != nil {
					return err
				}
				ai.IsActive = false
			}
		}
	}
	return nil
}

func (c *CardInstance) OnBeingAttacked(attacker *CardInstance) error {
	var ability Ability
	for _, ai := range c.AbilitiesInstances {
		if ai.Trigger == zb.CardAbilityTrigger_GotDamage {
			switch abilityInstance := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_Rage:
				rage := abilityInstance.Rage
				ability = NewRage(c, rage)
				if err := ability.Apply(c.Gameplay); err != nil {
					return err
				}
			}
		}
	}

	if c.Instance.Defense <= 0 {
		if err := c.OnDeath(c); err != nil {
			return err
		}
	}

	return nil
}

func (c *CardInstance) OnDeath(attacker *CardInstance) error {
	// trigger target ability on death
	var ability Ability
	for _, ai := range c.AbilitiesInstances {
		if !ai.IsActive {
			continue
		}
		if ai.Trigger == zb.CardAbilityTrigger_Death {
			switch abilityInstance := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_Reanimate:
				reanimate := abilityInstance.Reanimate
				ability = NewReanimate(c, reanimate)
				if err := ability.Apply(c.Gameplay); err != nil {
					return err
				}
				ai.IsActive = false
			}
		}
	}

	// trigger attacker ability
	for _, ai := range attacker.AbilitiesInstances {
		if ai.Trigger == zb.CardAbilityTrigger_Attack {
			switch abilityInstance := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_PriorityAttack:
				priorityAttack := abilityInstance.PriorityAttack
				ability = NewPriorityAttack(attacker, priorityAttack)
				if err := ability.Apply(c.Gameplay); err != nil {
					return err
				}
			}
		}
	}

	// after apply ability, update zone if the card instance is really dead and not moved to graveyard
	if c.Instance.Defense <= 0 && c.Zone != zb.Zone_GRAVEYARD {
		if err := c.MoveZone(zb.Zone_PLAY, zb.Zone_GRAVEYARD); err != nil {
			return err
		}
	}

	return nil
}

func (c *CardInstance) OnPlay() error {
	// trigger card ability on play
	var ability Ability
	for _, ai := range c.AbilitiesInstances {
		if ai.Trigger == zb.CardAbilityTrigger_Entry {
			switch abilityInstance := ai.AbilityType.(type) {
			case *zb.CardAbilityInstance_AttackOverlord:
				attackOverlord := abilityInstance.AttackOverlord
				ability = NewAttackOverlord(c, attackOverlord)
				if err := ability.Apply(c.Gameplay); err != nil {
					return err
				}
			case *zb.CardAbilityInstance_ReplaceUnitsWithTypeOnStrongerOnes:
				replaceUnitsWithTypeOnStrongerOnes := abilityInstance.ReplaceUnitsWithTypeOnStrongerOnes
				ability = NewReplaceUnitsWithTypeOnStrongerOnes(c, replaceUnitsWithTypeOnStrongerOnes, c.Gameplay.cardLibrary)
				if err := ability.Apply(c.Gameplay); err != nil {
					return err
				}
				ai.IsActive = false
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
		return fmt.Errorf("invalid moing from %v to %v", from, to)
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

func (c *CardInstance) AttackOverlord(target *zb.PlayerState, attacker *zb.PlayerState) error {
	target.Defense -= c.Instance.Damage

	if target.Defense <= 0 {
		c.Gameplay.State.Winner = attacker.Id
		c.Gameplay.State.IsEnded = true
		return nil
	}
	return c.AfterAttacking(target.InstanceId)
}

func (c *CardInstance) Mulligan() error {
	owner := c.Player()
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
	return newCardInstance.MoveZone(zb.Zone_DECK, zb.Zone_HAND)
}

func (c *CardInstance) Player() *zb.PlayerState {
	if int(c.OwnerIndex) > len(c.Gameplay.State.PlayerStates)-1 {
		return nil
	}
	return c.Gameplay.State.PlayerStates[c.OwnerIndex]
}

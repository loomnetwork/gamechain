package battleground

import (
	"fmt"
	"math/rand"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
)

type replaceUnitsWithTypeOnStrongerOnes struct {
	*CardInstance
	cardAbility *zb.CardAbilityReplaceUnitsWithTypeOnStrongerOnes
	cardlibrary *zb.CardList
}

var _ Ability = &replaceUnitsWithTypeOnStrongerOnes{}

func NewReplaceUnitsWithTypeOnStrongerOnes(card *CardInstance, cardAbility *zb.CardAbilityReplaceUnitsWithTypeOnStrongerOnes, cardlibrary *zb.CardList) *replaceUnitsWithTypeOnStrongerOnes {
	return &replaceUnitsWithTypeOnStrongerOnes{
		CardInstance: card,
		cardAbility:  cardAbility,
		cardlibrary:  cardlibrary,
	}
}

func (c *replaceUnitsWithTypeOnStrongerOnes) Apply(gameplay *Gameplay) error {
	owner := c.Player()
	if owner == nil {
		return fmt.Errorf("no owner for card instance %d", c.InstanceId)
	}
	// find the cards in card library with same types as cards in plays
	var toReplaceCards []*zb.PlayerActionOutcome_CardAbilityReplaceUnitsWithTypeOnStrongerOnesOutcome_NewCardInstance
	var oldInstanceIds []*zb.InstanceId
	for i, card := range owner.CardsInPlay {
		if c.Instance.Faction == card.Instance.Faction && !proto.Equal(c.InstanceId, card.InstanceId) {
			toReplaceCards = append(toReplaceCards, &zb.PlayerActionOutcome_CardAbilityReplaceUnitsWithTypeOnStrongerOnesOutcome_NewCardInstance{
				CardInstance: card,
				Position:     int32(i),
			})
			oldInstanceIds = append(oldInstanceIds, card.InstanceId)
		}
	}

	// do nothing if there is no same type card in play
	if len(toReplaceCards) == 0 {
		return nil
	}

	sameTypeStrongerFn := func(cardLibrary *zb.CardList, target *zb.CardInstance) []*zb.Card {
		var sameTypeStrongerCards []*zb.Card
		for _, card := range cardLibrary.Cards {
			if card.Faction == target.Instance.Faction && card.GooCost > target.Instance.GooCost {
				sameTypeStrongerCards = append(sameTypeStrongerCards, card)
			}
		}
		return sameTypeStrongerCards
	}

	state := gameplay.State

	var newcardInstances []*zb.PlayerActionOutcome_CardAbilityReplaceUnitsWithTypeOnStrongerOnesOutcome_NewCardInstance
	for i, card := range toReplaceCards {
		sameTypeStrongerCards := sameTypeStrongerFn(c.cardlibrary, card.CardInstance)
		if len(sameTypeStrongerCards) == 0 {
			continue
		}
		var r = rand.New(rand.NewSource(state.RandomSeed))
		randomCardIndex := r.Perm(len(sameTypeStrongerCards))

		// create new instance from card
		newcard := sameTypeStrongerCards[randomCardIndex[i]]
		instanceid := &zb.InstanceId{Id: state.NextInstanceId}
		state.NextInstanceId++
		newinstance := newCardInstanceFromCardDetails(newcard, instanceid, c.Owner, c.OwnerIndex)
		newinstance.Zone = zb.Zone_PLAY
		newcardInstances = append(newcardInstances, &zb.PlayerActionOutcome_CardAbilityReplaceUnitsWithTypeOnStrongerOnesOutcome_NewCardInstance{
			CardInstance: newinstance,
			Position:     card.Position,
		})
	}

	// replace card in play
	for i := 0; i < len(owner.CardsInPlay); i++ {
		for j := 0; j < len(newcardInstances); j++ {
			newcard := newcardInstances[j]
			if i == int(newcard.Position) {
				owner.CardsInPlay[i] = newcard.CardInstance
			}
		}
	}

	// outcome
	gameplay.actionOutcomes = append(gameplay.actionOutcomes, &zb.PlayerActionOutcome{
		Outcome: &zb.PlayerActionOutcome_ReplaceUnitsWithTypeOnStrongerOnes{
			ReplaceUnitsWithTypeOnStrongerOnes: &zb.PlayerActionOutcome_CardAbilityReplaceUnitsWithTypeOnStrongerOnesOutcome{
				NewCardInstances: newcardInstances,
				OldInstanceIds:   oldInstanceIds,
			},
		},
	})
	return nil
}

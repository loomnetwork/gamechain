package battleground

import (
	"fmt"
	"math/rand"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
)

type replaceUnitsWithTypeOnStrongerOnes struct {
	*CardInstance
	cardAbility *zb_data.CardAbilityReplaceUnitsWithTypeOnStrongerOnes
	cardlibrary *zb_data.CardList
}

var _ Ability = &replaceUnitsWithTypeOnStrongerOnes{}

func NewReplaceUnitsWithTypeOnStrongerOnes(card *CardInstance, cardAbility *zb_data.CardAbilityReplaceUnitsWithTypeOnStrongerOnes, cardlibrary *zb_data.CardList) *replaceUnitsWithTypeOnStrongerOnes {
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
	var toReplaceCards []*zb_data.PlayerActionOutcome_CardAbilityReplaceUnitsWithTypeOnStrongerOnesOutcome_NewCardInstance
	var oldInstanceIds []*zb_data.InstanceId
	for i, card := range owner.CardsInPlay {
		if c.Instance.Faction == card.Instance.Faction && !proto.Equal(c.InstanceId, card.InstanceId) {
			toReplaceCards = append(toReplaceCards, &zb_data.PlayerActionOutcome_CardAbilityReplaceUnitsWithTypeOnStrongerOnesOutcome_NewCardInstance{
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

	sameTypeStrongerFn := func(cardLibrary *zb_data.CardList, target *zb_data.CardInstance) []*zb.Card {
		var sameTypeStrongerCards []*zb.Card
		for _, card := range cardLibrary.Cards {
			if card.Faction == target.Instance.Faction && card.Cost > target.Instance.Cost {
				sameTypeStrongerCards = append(sameTypeStrongerCards, card)
			}
		}
		return sameTypeStrongerCards
	}

	state := gameplay.State

	var newcardInstances []*zb_data.PlayerActionOutcome_CardAbilityReplaceUnitsWithTypeOnStrongerOnesOutcome_NewCardInstance
	for i, card := range toReplaceCards {
		sameTypeStrongerCards := sameTypeStrongerFn(c.cardlibrary, card.CardInstance)
		if len(sameTypeStrongerCards) == 0 {
			continue
		}
		var r = rand.New(rand.NewSource(state.RandomSeed))
		randomCardIndex := r.Perm(len(sameTypeStrongerCards))

		// create new instance from card
		newcard := sameTypeStrongerCards[randomCardIndex[i]]
		instanceid := &zb_data.InstanceId{Id: state.NextInstanceId}
		state.NextInstanceId++
		newinstance := newCardInstanceFromCardDetails(newcard, instanceid, c.Owner, c.OwnerIndex)
		newinstance.Zone = zb_enums.Zone_PLAY
		newcardInstances = append(newcardInstances, &zb_data.PlayerActionOutcome_CardAbilityReplaceUnitsWithTypeOnStrongerOnesOutcome_NewCardInstance{
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
	gameplay.actionOutcomes = append(gameplay.actionOutcomes, &zb_data.PlayerActionOutcome{
		Outcome: &zb_data.PlayerActionOutcome_ReplaceUnitsWithTypeOnStrongerOnes{
			ReplaceUnitsWithTypeOnStrongerOnes: &zb_data.PlayerActionOutcome_CardAbilityReplaceUnitsWithTypeOnStrongerOnesOutcome{
				NewCardInstances: newcardInstances,
				OldInstanceIds:   oldInstanceIds,
			},
		},
	})
	return nil
}

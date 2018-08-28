package battleground

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/loomnetwork/zombie_battleground/types/zb"
)

const (
	MaxDeckNameChar = 48
)

var (
	ErrDeckNameExists  = errors.New("deck name already exists")
	ErrDeckNameEmpty   = errors.New("deck name cannot be empty")
	ErrDeckMustNotNil  = errors.New("deck must not be nil")
	ErrDeckNameTooLong = fmt.Errorf("deck name is more than %d characters", MaxDeckNameChar)
)

func validateDeckCollections(userCollections []*zb.CardCollection, deckCollections []*zb.CardCollection) error {
	maxAmountMap := make(map[string]int64)
	for _, collection := range userCollections {
		maxAmountMap[collection.CardName] = collection.Amount
	}

	var errorString = ""
	for _, collection := range deckCollections {
		cardAmount, ok := maxAmountMap[collection.CardName]
		if !ok {
			return fmt.Errorf("cannot add card %s", collection.CardName)
		}
		if cardAmount < collection.Amount {
			errorString += fmt.Sprintf("%s: %d ", collection.CardName, cardAmount)
		}
	}

	if errorString != "" {
		return fmt.Errorf("cannot add more than maximum for these cards: %s", errorString)
	}
	return nil
}

func validateDeckName(deckList []*zb.Deck, validatedDeck *zb.Deck) error {
	validatedDeck.Name = strings.TrimSpace(validatedDeck.Name)
	if len(validatedDeck.Name) == 0 {
		return ErrDeckNameEmpty
	}
	if utf8.RuneCountInString(validatedDeck.Name) > MaxDeckNameChar {
		return ErrDeckNameTooLong
	}
	for _, deck := range deckList {
		// Skip name validation for same deck id - support renaming deck
		if deck.Id == validatedDeck.Id {
			continue
		}
		if strings.EqualFold(deck.Name, validatedDeck.Name) {
			return ErrDeckNameExists
		}
	}
	return nil
}

func mergeDeckSets(deckSet1 []*zb.Deck, deckSet2 []*zb.Deck) []*zb.Deck {
	deckMap := make(map[int64]*zb.Deck)

	for _, deck := range deckSet1 {
		deckMap[deck.Id] = deck
	}

	for _, deck := range deckSet2 {
		deckMap[deck.Id] = deck
	}

	newArray := make([]*zb.Deck, len(deckMap))

	i := 0
	for j := len(deckSet2) - 1; j >= 0; j -= 1 {
		deck := deckSet2[j]

		newDeck, ok := deckMap[deck.Id]
		if !ok {
			continue
		}

		newArray[i] = newDeck
		i++

		delete(deckMap, deck.Id)
	}

	for j := len(deckSet1) - 1; j >= 0; j -= 1 {
		deck := deckSet1[j]

		newDeck, ok := deckMap[deck.Id]
		if !ok {
			continue
		}

		newArray[i] = newDeck
		i++

		delete(deckMap, deck.Id)
	}

	return newArray
}

func getHeroById(heroList []*zb.Hero, heroId int64) *zb.Hero {
	for _, hero := range heroList {
		if hero.HeroId == heroId {
			return hero
		}
	}
	return nil
}

func validateDeckHero(heroList []*zb.Hero, heroID int64) error {
	// check if the user has hero
	if getHeroById(heroList, heroID) != nil {
		return nil
	}
	return fmt.Errorf("hero: %d cannot be part of deck, since it is not owned by User", heroID)
}

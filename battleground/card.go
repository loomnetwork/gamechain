package battleground

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/loomnetwork/zombie_battleground/types/zb"
)

func validateDeckCollections(userCollections []*zb.CardCollection, deckCollections []*zb.CardCollection) error {
	maxAmountMap := make(map[string]int64)

	for _, collection := range userCollections {
		maxAmountMap[collection.CardName] = collection.Amount
	}

	var errorString = ""
	for _, collection := range deckCollections {
		if maxAmountMap[collection.CardName] < collection.Amount {
			errorString += fmt.Sprintf("%s: %d ", collection.CardName, maxAmountMap[collection.CardName])
		}
	}

	if errorString != "" {
		return errors.New("Cannot add more than maximum for these cards: " + errorString)
	} else {
		return nil
	}
}

func validateDeckName(deckList []*zb.Deck, validatedDeck *zb.Deck) error {
	validatedDeck.Name = strings.TrimSpace(validatedDeck.Name)
	if len(validatedDeck.Name) == 0 {
		return errors.New("deck name can't be empty")
	}

	if utf8.RuneCountInString(validatedDeck.Name) > 48 {
		return errors.New("deck name is more than 48 characters")
	}

	for _, deck := range deckList {
		// Skip name validation for same deck id
		if deck.Id == validatedDeck.Id {
			continue
		}

		if strings.EqualFold(deck.Name, validatedDeck.Name) {
			return errors.New("deck name already exists")
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

func editDeck(deckSet []*zb.Deck, deck *zb.Deck) error {
	deckToEdit := getDeckById(deckSet, deck.Id)

	if deckToEdit == nil {
		return fmt.Errorf("Unable to find deck: %d", deck.Id)
	}

	deckToEdit.Name = deck.Name
	deckToEdit.Cards = deck.Cards
	deckToEdit.HeroId = deck.HeroId

	return nil
}

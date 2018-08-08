package battleground

import (
	"fmt"
	"errors"

	"github.com/loomnetwork/zombie_battleground/types/zb"
)

func validateDeckCollections(userCollections []*zb.CardCollection, deckCollections []*zb.CardCollection) error {
	maxAmountMap := make(map[string]int64)

	for _, collection := range userCollections {
		maxAmountMap[collection.CardName] = collection.Amount
	}

	var error error = nil
	var errorString = ""
	for _, collection := range deckCollections {
		if maxAmountMap[collection.CardName] < collection.Amount {
			errorString += fmt.Sprintf("you cannot add more than %d for your card name: %s\n", maxAmountMap[collection.CardName], collection.CardName)
		}
	}

	if errorString != "" {
		error = errors.New(errorString)
	}

	return error
}

func mergeDeckSets(deckSet1 []*zb.Deck, deckSet2 []*zb.Deck) []*zb.Deck {
	deckMap := make(map[string]*zb.Deck)

	for _, deck := range deckSet1 {
		deckMap[deck.Name] = deck
	}

	for _, deck := range deckSet2 {
		deckMap[deck.Name] = deck
	}

	newArray := make([]*zb.Deck, len(deckMap))

	i := 0
	for j := len(deckSet2) - 1; j >= 0; j -= 1 {
		deck := deckSet2[j]

		newDeck, ok := deckMap[deck.Name]
		if !ok {
			continue
		}

		newArray[i] = newDeck
		i++

		delete(deckMap, deck.Name)
	}

	for j := len(deckSet1) - 1; j >= 0; j -= 1 {
		deck := deckSet1[j]

		newDeck, ok := deckMap[deck.Name]
		if !ok {
			continue
		}

		newArray[i] = newDeck
		i++

		delete(deckMap, deck.Name)
	}

	return newArray
}

func editDeck(deckSet []*zb.Deck, deck *zb.Deck) error {
	var deckToEdit *zb.Deck

	for _, deckInSet := range deckSet {
		if deck.Name == deckInSet.Name {
			deckToEdit = deckInSet
			break
		}
	}

	if deckToEdit == nil {
		return fmt.Errorf("Unable to find deck: %s", deck.Name)
	}

	deckToEdit.Cards = deck.Cards
	deckToEdit.HeroId = deck.HeroId

	return nil
}

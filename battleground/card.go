package battleground

import (
	"fmt"

	"github.com/loomnetwork/zombie_battleground/types/zb"
)

func validateDeckAddition(collections []*zb.CardCollection, collectionList []*zb.CardCollectionList) error {
	maxAmountMap := make(map[int64]int64)

	for _, card := range collections {
		maxAmountMap[card.CardId] = card.Amount
	}

	for _, collist := range collectionList {
		for _, card := range collist.Cards {
			if maxAmountMap[card.CardId] < card.Amount {
				return fmt.Errorf("you cannot add more than %d for your card id: %d", maxAmountMap[card.CardId], card.CardId)
			}
		}
	}

	return nil
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

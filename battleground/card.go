package battleground

import (
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"math/rand"
	"strings"
	"unicode/utf8"

	"github.com/loomnetwork/gamechain/types/zb"
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

func validateCardLibrary(cards []*zb.Card, deckCards []*zb.DeckCard) error {
	cardmap := make(map[string]interface{})
	for _, card := range cards {
		cardmap[card.Name] = struct{}{}
	}
	for _, collection := range deckCards {
		if _, ok := cardmap[collection.CardName]; !ok {
			return fmt.Errorf("card %s not found in card library", collection.CardName)
		}
	}
	return nil
}

func validateDeckCollections(userCollections []*zb.CardCollectionCard, deckCollections []*zb.CardCollectionCard) error {
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

func getOverlordById(overlordList []*zb.Overlord, overlordId int64) *zb.Overlord {
	for _, overlord := range overlordList {
		if overlord.OverlordId == overlordId {
			return overlord
		}
	}
	return nil
}

func validateDeckOverlord(overlordList []*zb.Overlord, overlordID int64) error {
	// check if the user has overlord
	if getOverlordById(overlordList, overlordID) != nil {
		return nil
	}
	return fmt.Errorf("overlord: %d cannot be part of deck, since it is not owned by User", overlordID)
}

func shuffleCardInDeck(deck []*zb.CardInstance, seed int64, playerIndex int) []*zb.CardInstance {
	r := rand.New(rand.NewSource(seed + int64(playerIndex)))
	for i := 0; i < len(deck); i++ {
		n := r.Intn(i + 1)
		// do a swap
		if i != n {
			deck[n], deck[i] = deck[i], deck[n]
		}
	}
	return deck
}

func drawFromCardList(cardlist []*zb.Card, n int) (cards []*zb.Card, renaming []*zb.Card) {
	var i int
	for i = 0; i < n; i++ {
		if i > len(cardlist)-1 {
			break
		}
		cards = append(cards, cardlist[i])
	}
	// update cardlist
	renaming = cardlist[i:]
	return
}

func findCardInCardListByName(card *zb.CardInstance, cards []*zb.CardInstance) (int, *zb.CardInstance, bool) {
	for i, c := range cards {
		if card.Prototype.Name == c.Prototype.Name {
			return i, c, true
		}
	}
	return -1, nil, false
}

func findCardInCardListByInstanceId(instanceId *zb.InstanceId, cards []*zb.CardInstance) (int, *zb.CardInstance, bool) {
	for i, c := range cards {
		if proto.Equal(instanceId, c.InstanceId) {
			return i, c, true
		}
	}
	return -1, nil, false
}
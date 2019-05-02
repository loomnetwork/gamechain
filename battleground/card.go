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

func validateCardLibraryCards(cardLibrary []*zb.Card) error {
	for _, card := range cardLibrary {
		if card.PictureTransform == nil || card.PictureTransform.Position == nil || card.PictureTransform.Scale == nil {
			return fmt.Errorf("card '%s' (mould id %d) missing value for PictureTransform field", card.Name, card.MouldId)
		}
	}

	return nil
}

func validateDeckCards(cardLibrary []*zb.Card, deckCards []*zb.DeckCard) error {
	cardMap := make(map[int64]interface{})
	for _, card := range cardLibrary {
		cardMap[card.MouldId] = struct{}{}
	}
	for _, deckCard := range deckCards {
		if deckCard.MouldId <= 0 {
			return fmt.Errorf("mould id not set for card %s", deckCard.CardNameDeprecated)
		}

		if _, ok := cardMap[deckCard.MouldId]; !ok {
			return fmt.Errorf("card with mould id %d not found in card library", deckCard.MouldId)
		}
	}
	return nil
}

func validateDeckCollections(userCollections []*zb.CardCollectionCard, deckCollections []*zb.CardCollectionCard) error {
	maxAmountMap := make(map[int64]int64)
	for _, collection := range userCollections {
		maxAmountMap[collection.MouldId] = collection.Amount
	}

	var errorString = ""
	for _, collection := range deckCollections {
		cardAmount, ok := maxAmountMap[collection.MouldId]
		if !ok {
			return fmt.Errorf("cannot add card %d", collection.MouldId)
		}
		if cardAmount < collection.Amount {
			errorString += fmt.Sprintf("%d: %d ", collection.MouldId, cardAmount)
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
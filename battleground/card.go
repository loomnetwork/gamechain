package battleground

import (
	"errors"
	"fmt"
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

func validateCardLibrary(cards []*zb.Card, deckCollections []*zb.CardCollection) error {
	cardmap := make(map[string]interface{})
	for _, card := range cards {
		cardmap[card.Name] = struct{}{}
	}
	for _, collection := range deckCollections {
		if _, ok := cardmap[collection.CardName]; !ok {
			return fmt.Errorf("card %s not found in card library", collection.CardName)
		}
	}
	return nil
}

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

func cardInstanceFromDeck(deck *zb.Deck) (cards []*zb.CardInstance) {
	for _, collection := range deck.Cards {
		for i := int64(0); i < collection.Amount; i++ {
			// TODO: finalize that fields are needed
			cards = append(cards, &zb.CardInstance{
				Prototype:  &zb.CardPrototype{Name: collection.CardName},
				InstanceId: 3, //TODO
				Attack:     4, //TODO
				Defence:    5, //TODO
			})
		}
	}
	return
}

func shuffleCardInDeck(deck *zb.Deck, seed int64) []*zb.CardInstance {
	cards := cardInstanceFromDeck(deck)
	r := rand.New(rand.NewSource(seed))
	for i := 0; i < len(cards); i++ {
		n := r.Intn(i + 1)
		// do a swap
		if i != n {
			cards[n], cards[i] = cards[i], cards[n]
		}
	}
	return cards
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

func containCardInCardList(card *zb.CardInstance, cards []*zb.CardInstance) (*zb.CardInstance, bool) {
	for _, c := range cards {
		if card.Prototype.Name == c.Prototype.Name {
			return c, true
		}
	}
	return nil, false
}

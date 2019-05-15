package battleground

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"unicode/utf8"

	"github.com/gogo/protobuf/proto"

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
	existingCardsSet, err := getMouldIdToCardMap(cardLibrary)
	if err != nil {
		return err
	}

	for _, card := range cardLibrary {
		if card.MouldId <= 0 {
			return fmt.Errorf("mould id not set for card %s", card.Name)
		}

		if card.PictureTransform == nil || card.PictureTransform.Position == nil || card.PictureTransform.Scale == nil {
			return fmt.Errorf("card '%s' (mould id %d) missing value for PictureTransform field", card.Name, card.MouldId)
		}

		err = validateSourceMouldId(card, existingCardsSet)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateDeckCards(cardLibrary []*zb.Card, deckCards []*zb_data.DeckCard) error {
	mouldIdToCard, err := getMouldIdToCardMap(cardLibrary)
	if err != nil {
		return err
	}

	for _, deckCard := range deckCards {
		if deckCard.MouldId <= 0 {
			return fmt.Errorf("mould id not set for card %s", deckCard.CardNameDeprecated)
		}

		if deckCard.CardNameDeprecated != "" {
			return fmt.Errorf("card %d has non-empty name '%s', must be empty", deckCard.MouldId, deckCard.CardNameDeprecated)
		}

		if _, ok := mouldIdToCard[deckCard.MouldId]; !ok {
			return fmt.Errorf("card with mould id %d not found in card library", deckCard.MouldId)
		}
	}
	return nil
}

func validateDeckCollections(userCollections []*zb_data.CardCollectionCard, deckCollections []*zb_data.CardCollectionCard) error {
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

func validateDeckName(deckList []*zb_data.Deck, validatedDeck *zb_data.Deck) error {
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

func shuffleCardInDeck(deck []*zb_data.CardInstance, seed int64, playerIndex int) []*zb_data.CardInstance {
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

func findCardInCardListByName(card *zb_data.CardInstance, cards []*zb_data.CardInstance) (int, *zb_data.CardInstance, bool) {
	for i, c := range cards {
		if card.Prototype.Name == c.Prototype.Name {
			return i, c, true
		}
	}
	return -1, nil, false
}

func findCardInCardListByInstanceId(instanceId *zb_data.InstanceId, cards []*zb_data.CardInstance) (int, *zb_data.CardInstance, bool) {
	for i, c := range cards {
		if proto.Equal(instanceId, c.InstanceId) {
			return i, c, true
		}
	}
	return -1, nil, false
}

func getMouldIdToCardMap(cardLibrary []*zb.Card) (map[int64]*zb.Card, error) {
	existingCardsSet := make(map[int64]*zb.Card)
	for _, card := range cardLibrary {
		_, exists := existingCardsSet[card.MouldId]
		if !exists {
			existingCardsSet[card.MouldId] = card
		} else {
			return nil, fmt.Errorf("more than one card has mould ID %d, this is not allowed", card.MouldId)
		}
	}

	return existingCardsSet, nil
}

func applySourceMouldIdAndOverrides(card *zb.Card, mouldIdToCard map[int64]*zb.Card) error {
	if card.SourceMouldId <= 0 {
		return nil
	}

	sourceMouldId := card.SourceMouldId
	overrides := card.Overrides
	sourceCard, exists := mouldIdToCard[sourceMouldId]
	if !exists {
		return fmt.Errorf("source card with mould id %d not found", sourceMouldId)
	}

	card.Reset()
	proto.Merge(card, sourceCard)

	if overrides == nil {
		return nil
	}

	// per-field merge
	if overrides.Kind != nil {
		card.Kind = overrides.Kind.Value
	}

	if overrides.Faction != nil {
		card.Faction = overrides.Faction.Value
	}

	if overrides.Name != nil {
		card.Name = overrides.Name.Value
	}

	if overrides.Description != nil {
		card.Description = overrides.Description.Value
	}

	if overrides.FlavorText != nil {
		card.FlavorText = overrides.FlavorText.Value
	}

	if overrides.Picture != nil {
		card.Picture = overrides.Picture.Value
	}

	if overrides.Rank != nil {
		card.Rank = overrides.Rank.Value
	}

	if overrides.Type != nil {
		card.Type = overrides.Type.Value
	}

	if overrides.Frame != nil {
		card.Frame = overrides.Frame.Value
	}

	if overrides.Damage != nil {
		card.Damage = overrides.Damage.Value
	}

	if overrides.Defense != nil {
		card.Defense = overrides.Defense.Value
	}

	if overrides.Cost != nil {
		card.Cost = overrides.Cost.Value
	}

	if overrides.PictureTransform != nil {
		card.PictureTransform = overrides.PictureTransform
	}

	if len(overrides.Abilities) > 0 {
		card.Abilities = overrides.Abilities
	}

	if overrides.UniqueAnimation != nil {
		card.UniqueAnimation = overrides.UniqueAnimation.Value
	}

	if overrides.Hidden != nil {
		card.Hidden = overrides.Hidden.Value
	}

	return nil
}

func validateSourceMouldId(card *zb.Card, mouldIdToCard map[int64]*zb.Card) error {
	if card.SourceMouldId <= 0 {
		return nil
	}

	if card.SourceMouldId == card.MouldId {
		return fmt.Errorf(
			"card '%s' (mould id %d) has sourceMouldId equal to mouldId",
			card.Name,
			card.MouldId,
		)
	}

	sourceCard, exists := mouldIdToCard[card.SourceMouldId]
	if !exists {
		return fmt.Errorf(
			"card '%s' (mould id %d) has sourceMouldId %d, but such mould id is not found",
			card.Name,
			card.MouldId,
			card.SourceMouldId,
		)
	}

	if sourceCard.SourceMouldId > 0 {
		return fmt.Errorf("source card %d can't have sourceMouldId set itself", sourceCard.SourceMouldId)
	}

	return nil
}

package battleground

import (
	"fmt"
	battleground_proto "github.com/loomnetwork/gamechain/battleground/proto"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	"github.com/pkg/errors"
	"math/rand"
	"strings"
	"unicode/utf8"

	"github.com/gogo/protobuf/proto"
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

func validateCardLibraryCards(cardLibrary []*zb_data.Card) error {
	existingCardsSet, err := getCardKeyToCardMap(cardLibrary)
	if err != nil {
		return err
	}

	for _, card := range cardLibrary {
		if card.CardKey.MouldId <= 0 {
			return fmt.Errorf("mould id not set for card %s", card.Name)
		}

		if card.CardKey.Variant == zb_enums.CardVariant_Standard {
			if card.PictureTransform == nil || card.PictureTransform.Position == nil || card.PictureTransform.Scale == nil {
				return fmt.Errorf("card '%s' (card key %d) missing value for PictureTransform field", card.Name, card.CardKey.MouldId)
			}
		}

		err = validateCardEdition(card, existingCardsSet)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateDeckCards(cardLibrary []*zb_data.Card, deckCards []*zb_data.DeckCard) error {
	cardKeyToCard, err := getCardKeyToCardMap(cardLibrary)
	if err != nil {
		return err
	}

	for _, deckCard := range deckCards {
		if deckCard.CardKey.MouldId <= 0 {
			return fmt.Errorf("mould id not set for card [%v]", deckCard.CardKey)
		}

		if _, ok := cardKeyToCard[deckCard.CardKey]; !ok {
			return fmt.Errorf("card with cardKey [%s] not found in card library", deckCard.CardKey.String())
		}
	}
	return nil
}

func validateDeck(isEditDeck bool, cardLibrary *zb_data.CardList, deck *zb_data.Deck, deckList []*zb_data.Deck, overlords []*zb_data.OverlordUserInstance) error {
	// validate version on card library
	if err := validateDeckCards(cardLibrary.Cards, deck.Cards); err != nil {
		return errors.Wrap(err, "error validating deck cards")
	}

	// Since the server side does not have any knowleadge on user's collection, we skip this logic on the server side for now.
	// TODO: Turn on the check when the server side knows user's collection
	// validating against default card collection
	// var defaultCollection zb.CardCollectionList
	// if err := ctx.Get(MakeVersionedKey(req.Version, defaultCollectionKey), &defaultCollection); err != nil {
	// 	return nil, errors.Wrapf(err, "unable to get default collectionlist")
	// }
	// // make sure the given cards and amount must be a subset of user's cards
	// if err := validateDeckCollections(defaultCollection.Cards, req.Deck.Cards); err != nil {
	// 	return nil, err
	// }

	if err := validateDeckName(deckList, deck); err != nil {
		return errors.Wrap(err, "error validating deck name")
	}

	return nil
}

func validateDeckCollections(userCollections []*zb_data.CardCollectionCard, deckCollections []*zb_data.CardCollectionCard) error {
	maxAmountMap := make(map[battleground_proto.CardKey]int64)
	for _, collection := range userCollections {
		maxAmountMap[collection.CardKey] = collection.Amount
	}

	var errorString = ""
	for _, collection := range deckCollections {
		cardAmount, ok := maxAmountMap[collection.CardKey]
		if !ok {
			return fmt.Errorf("cannot add card [%s]", collection.CardKey.String())
		}
		if cardAmount < collection.Amount {
			errorString += fmt.Sprintf("[%s]: %d ", collection.CardKey.String(), cardAmount)
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
		// Skip same-name validation for same deck id - support renaming deck
		if deck.Id == validatedDeck.Id {
			continue
		}
		if strings.EqualFold(deck.Name, validatedDeck.Name) {
			return ErrDeckNameExists
		}
	}
	return nil
}

func getOverlordUserDataByPrototypeId(overlordsUserData []*zb_data.OverlordUserData, overlordPrototypeId int64) (*zb_data.OverlordUserData, bool){
	for _, overlordUserData := range overlordsUserData {
		if overlordUserData.PrototypeId == overlordPrototypeId {
			return overlordUserData, true
		}
	}
	return nil, false
}

func getOverlordUserInstanceByPrototypeId(overlordsUserInstance []*zb_data.OverlordUserInstance, overlordPrototypeId int64) (*zb_data.OverlordUserInstance, bool) {
	for _, overlordUserInstance := range overlordsUserInstance {
		if overlordUserInstance.Prototype.Id == overlordPrototypeId {
			return overlordUserInstance, true
		}
	}
	return nil, false
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

func drawFromCardList(cardlist []*zb_data.Card, n int) (cards []*zb_data.Card, renaming []*zb_data.Card) {
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

func getCardKeyToCardMap(cardLibrary []*zb_data.Card) (map[battleground_proto.CardKey]*zb_data.Card, error) {
	existingCardsSet := make(map[battleground_proto.CardKey]*zb_data.Card)
	for _, card := range cardLibrary {
		_, exists := existingCardsSet[card.CardKey]
		if !exists {
			existingCardsSet[card.CardKey] = card
		} else {
			return nil, fmt.Errorf("more than one card has cardKey [%s], this is not allowed", card.CardKey.String())
		}
	}

	return existingCardsSet, nil
}

func applySourceMouldIdAndOverrides(card *zb_data.Card, cardKeyToCard map[battleground_proto.CardKey]*zb_data.Card) error {
	if card.CardKey.Variant == zb_enums.CardVariant_Standard {
		return nil
	}

	cardKey := card.CardKey
	sourceCardKey := battleground_proto.CardKey{
		MouldId: card.CardKey.MouldId,
		Variant: zb_enums.CardVariant_Standard,
	}

	overrides := card.Overrides
	sourceCard, exists := cardKeyToCard[sourceCardKey]
	if !exists {
		return fmt.Errorf("source card with cardKey [%s] not found", sourceCardKey.String())
	}

	card.Reset()
	proto.Merge(card, sourceCard)
	card.CardKey = cardKey

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

func validateCardEdition(card *zb_data.Card, cardKeyToCard map[battleground_proto.CardKey]*zb_data.Card) error {
	if card.CardKey.Variant == zb_enums.CardVariant_Standard {
		return nil
	}

	sourceCardKey := battleground_proto.CardKey{
		MouldId: card.CardKey.MouldId,
		Variant: zb_enums.CardVariant_Standard,
	}

	_, exists := cardKeyToCard[sourceCardKey]
	if !exists {
		return fmt.Errorf(
			"card '%s' (cardKey [%s]) has edition %s, but Normal edition is not found for such mouldId",
			card.Name,
			card.CardKey.String(),
			zb_enums.CardVariant_Enum_name[int32(card.CardKey.Variant)],
		)
	}

	return nil
}

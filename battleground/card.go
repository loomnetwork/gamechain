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
				return fmt.Errorf("card '%s' (card key %s) missing value for PictureTransform field", card.Name, card.CardKey.String())
			}
		}

		// FIXME: add check
		/*if card.Type == zb_enums.CardType_Undefined {
			return fmt.Errorf("type is not set for card '%s' (card key %s)", card.Name, card.CardKey.String())
		}*/

		if card.Kind == zb_enums.CardKind_Undefined {
			return fmt.Errorf("kind is not set for card '%s' (card key %s)", card.Name, card.CardKey.String())
		}

		if card.Rank == zb_enums.CreatureRank_Undefined {
			return fmt.Errorf("rank is not set for card '%s' (card key %s)", card.Name, card.CardKey.String())
		}

		if card.Faction == zb_enums.Faction_None {
			return fmt.Errorf("faction is not set for card '%s' (card key %s)", card.Name, card.CardKey.String())
		}

		err = validateCardVariant(card, existingCardsSet)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateDeckAgainstCardLibrary(cardLibrary []*zb_data.Card, deckCards []*zb_data.DeckCard) error {
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

func validateDeck(isEditDeck bool, cardLibrary *zb_data.CardList, userCardCollection []*zb_data.CardCollectionCard, deck *zb_data.Deck, deckList []*zb_data.Deck, overlords []*zb_data.OverlordUserInstance) error {
	if err := validateDeckAgainstUserCardCollection(userCardCollection, deck.Cards); err != nil {
		return errors.Wrap(err, "error validating deck cards")
	}

	if err := validateDeckName(deckList, deck); err != nil {
		return errors.Wrap(err, "error validating deck name")
	}

	// TODO: check for unlocked overlord skills

	return nil
}

func validateDeckAgainstUserCardCollection(userCardCollection []*zb_data.CardCollectionCard, deckCards []*zb_data.DeckCard) error {
	amountMap := make(map[battleground_proto.CardKey]int64)
	for _, collectionCard := range userCardCollection {
		amountMap[collectionCard.CardKey] = collectionCard.Amount
	}

	var errorString = ""
	for _, collection := range deckCards {
		cardAmount, ok := amountMap[collection.CardKey]
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

func findCardInCardListByCardKey(card *zb_data.CardInstance, cards []*zb_data.CardInstance) (int, *zb_data.CardInstance, bool) {
	for i, c := range cards {
		if card.Prototype.CardKey == c.Prototype.CardKey {
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

func validateCardVariant(card *zb_data.Card, cardKeyToCard map[battleground_proto.CardKey]*zb_data.Card) error {
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
			"card '%s' (cardKey [%s]) has variant %s, but Normal variant is not found for such mouldId",
			card.Name,
			card.CardKey.String(),
			zb_enums.CardVariant_Enum_name[int32(card.CardKey.Variant)],
		)
	}

	return nil
}

func cardKeyFromCardTokenId(cardTokenId int64) battleground_proto.CardKey {
	return battleground_proto.CardKey{
		MouldId: cardTokenId / 10,
		Variant: zb_enums.CardVariant_Enum(cardTokenId % 10),
	}
}

func newCardInstanceSpecificDataFromCardDetails(cardDetails *zb_data.Card) *zb_data.CardInstanceSpecificData {
	cardDetails = proto.Clone(cardDetails).(*zb_data.Card)
	return &zb_data.CardInstanceSpecificData{
		Damage:    cardDetails.Damage,
		Defense:   cardDetails.Defense,
		Type:      cardDetails.Type,
		Faction:   cardDetails.Faction,
		Cost:      cardDetails.Cost,
		Abilities: cardDetails.Abilities,
	}
}

func newCardInstanceFromCardDetails(cardDetails *zb_data.Card, instanceID *zb_data.InstanceId, owner string, ownerIndex int32) *zb_data.CardInstance {
	instance := newCardInstanceSpecificDataFromCardDetails(cardDetails)
	var abilities []*zb_data.CardAbilityInstance
	for _, raw := range cardDetails.Abilities {
		switch raw.Ability {
		case zb_enums.AbilityType_Rage:
			abilities = append(abilities, &zb_data.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb_data.CardAbilityInstance_Rage{
					Rage: &zb_data.CardAbilityRage{
						AddedDamage: raw.Value,
					},
				},
			})
		case zb_enums.AbilityType_PriorityAttack:
			abilities = append(abilities, &zb_data.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb_data.CardAbilityInstance_PriorityAttack{
					PriorityAttack: &zb_data.CardAbilityPriorityAttack{},
				},
			})
		case zb_enums.AbilityType_ReanimateUnit:
			abilities = append(abilities, &zb_data.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb_data.CardAbilityInstance_Reanimate{
					Reanimate: &zb_data.CardAbilityReanimate{
						DefaultDamage:  cardDetails.Damage,
						DefaultDefense: cardDetails.Defense,
					},
				},
			})
		case zb_enums.AbilityType_ChangeStat:
			abilities = append(abilities, &zb_data.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb_data.CardAbilityInstance_ChangeStat{
					ChangeStat: &zb_data.CardAbilityChangeStat{
						StatAdjustment: raw.Value,
						Stat:           raw.Stat,
					},
				},
			})
		case zb_enums.AbilityType_AttackOverlord:
			abilities = append(abilities, &zb_data.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb_data.CardAbilityInstance_AttackOverlord{
					AttackOverlord: &zb_data.CardAbilityAttackOverlord{
						Damage: raw.Value,
					},
				},
			})
		case zb_enums.AbilityType_ReplaceUnitsWithTypeOnStrongerOnes:
			abilities = append(abilities, &zb_data.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb_data.CardAbilityInstance_ReplaceUnitsWithTypeOnStrongerOnes{
					ReplaceUnitsWithTypeOnStrongerOnes: &zb_data.CardAbilityReplaceUnitsWithTypeOnStrongerOnes{
						Faction: cardDetails.Faction,
					},
				},
			})
		case zb_enums.AbilityType_DealDamageToThisAndAdjacentUnits:
			abilities = append(abilities, &zb_data.CardAbilityInstance{
				IsActive: true,
				Trigger:  raw.Trigger,
				AbilityType: &zb_data.CardAbilityInstance_DealDamageToThisAndAdjacentUnits{
					DealDamageToThisAndAdjacentUnits: &zb_data.CardAbilityDealDamageToThisAndAdjacentUnits{
						AdjacentDamage: cardDetails.Damage,
					},
				},
			})
		}

	}
	return &zb_data.CardInstance{
		InstanceId:         proto.Clone(instanceID).(*zb_data.InstanceId),
		Owner:              owner,
		Prototype:          proto.Clone(cardDetails).(*zb_data.Card),
		Instance:           instance,
		AbilitiesInstances: abilities,
		Zone:               zb_enums.Zone_DECK, // default to deck
		OwnerIndex:         ownerIndex,
	}
}

func getInstanceIdsFromCardInstances(cards []*zb_data.CardInstance) []*zb_data.InstanceId {
	var instanceIds = make([]*zb_data.InstanceId, len(cards), len(cards))
	for i := range cards {
		instanceIds[i] = cards[i].InstanceId
	}

	return instanceIds
}

func populateDeckCards(cardLibrary *zb_data.CardList, playerStates []*zb_data.PlayerState, useBackendGameLogic bool) error {
	for playerIndex, playerState := range playerStates {
		deck := playerState.Deck
		if deck == nil {
			return fmt.Errorf("no card deck fro player %s", playerState.Id)
		}
		for _, cardAmounts := range deck.Cards {
			for i := int64(0); i < cardAmounts.Amount; i++ {
				cardDetails, err := getCardByCardKey(cardLibrary, cardAmounts.CardKey)
				if err != nil {
					return fmt.Errorf("unable to get card [%s] from card library: %s", cardAmounts.CardKey.String(), err.Error())
				}

				cardInstance := newCardInstanceFromCardDetails(
					cardDetails,
					nil,
					playerState.Id,
					int32(playerIndex),
				)
				playerState.CardsInDeck = append(playerState.CardsInDeck, cardInstance)
			}
		}
	}

	removeUnsupportedCardFeatures(useBackendGameLogic, playerStates)

	return nil
}

func removeUnsupportedCardFeatures(useBackendGameLogic bool, playerStates []*zb_data.PlayerState) {
	if !useBackendGameLogic {
		return
	}

	for _, playerState := range playerStates {
		filteredCards := make([]*zb_data.CardInstance, 0, 0)

		for _, card := range playerState.CardsInDeck {
			filteredAbilities := make([]*zb_data.AbilityData, 0, 0)
			for _, ability := range card.Prototype.Abilities {
				switch ability.Ability {
				case zb_enums.AbilityType_Rage:
					fallthrough
				case zb_enums.AbilityType_PriorityAttack:
					fallthrough
				case zb_enums.AbilityType_ReanimateUnit:
					fallthrough
				case zb_enums.AbilityType_ChangeStat:
					fallthrough
				case zb_enums.AbilityType_AttackOverlord:
					fallthrough
				case zb_enums.AbilityType_ReplaceUnitsWithTypeOnStrongerOnes:
					filteredAbilities = append(filteredAbilities, ability)
				default:
					fmt.Printf("Unsupported AbilityType value %s, removed (card '%s')\n", zb_enums.AbilityType_Enum_name[int32(ability.Ability)], card.Prototype.Name)
				}
			}

			card.Prototype.Abilities = filteredAbilities

			switch card.Prototype.Type {
			case zb_enums.CardType_Feral:
				fallthrough
			case zb_enums.CardType_Heavy:
				fmt.Printf("Unsupported CardType value %s, fallback to WALKER (card %s)\n", zb_enums.CardType_Enum_name[int32(card.Prototype.Type)], card.Prototype.Name)
				card.Prototype.Type = zb_enums.CardType_Walker
			}

			switch card.Instance.Type {
			case zb_enums.CardType_Feral:
				fallthrough
			case zb_enums.CardType_Heavy:
				fmt.Printf("Unsupported CardType value %s, fallback to WALKER (card %s)\n", zb_enums.CardType_Enum_name[int32(card.Instance.Type)], card.Prototype.Name)
				card.Instance.Type = zb_enums.CardType_Walker
			}

			switch card.Prototype.Kind {
			case zb_enums.CardKind_Creature:
				filteredCards = append(filteredCards, card)
			default:
				fmt.Printf("Unsupported CardKind value %s, removed (card '%s')\n", zb_enums.CardKind_Enum_name[int32(card.Prototype.Kind)], card.Prototype.Name)
			}

			switch card.Prototype.Rank {
			case zb_enums.CreatureRank_Officer:
				fallthrough
			case zb_enums.CreatureRank_Commander:
				fallthrough
			case zb_enums.CreatureRank_General:
				fmt.Printf("Unsupported CreatureRank value %s, fallback to MINION (card %s)\n", zb_enums.CreatureRank_Enum_name[int32(card.Prototype.Rank)], card.Prototype.Name)
				card.Prototype.Rank = zb_enums.CreatureRank_Minion
			}
		}

		playerState.CardsInDeck = filteredCards
	}
}

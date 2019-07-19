package battleground

import (
	"encoding/hex"
	"fmt"
	"github.com/gogo/protobuf/proto"
	battleground_proto "github.com/loomnetwork/gamechain/battleground/proto"
	orctype "github.com/loomnetwork/gamechain/types/oracle"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/common"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/pkg/errors"
	"sort"
	"strings"
)

type CardAmountChangeItem struct {
	CardKey      battleground_proto.CardKey
	AmountChange int64
}

func getMaxAmountOfCardInDeck(card *zb_data.Card) (int, error) {
	const cardItemMaxCopies = 2
	const cardMinionMaxCopies = 4
	const cardOfficerMaxCopies = 2
	const cardCommanderMaxCopies = 2
	const cardGeneralMaxCopies = 1

	var maxCopies int
	if card.Faction == zb_enums.Faction_Item {
		maxCopies = cardItemMaxCopies
	} else {
		switch card.Rank {
		case zb_enums.CreatureRank_Minion:
			maxCopies = cardMinionMaxCopies
		case zb_enums.CreatureRank_Officer:
			maxCopies = cardOfficerMaxCopies
		case zb_enums.CreatureRank_Commander:
			maxCopies = cardCommanderMaxCopies
		case zb_enums.CreatureRank_General:
			maxCopies = cardGeneralMaxCopies
		default:
			return 0, fmt.Errorf("no rule to determine max copies of card %s", card.CardKey.String())
		}
	}

	return maxCopies, nil
}

func (z *ZombieBattleground) processOracleCommandResponseGetUserFullCardCollection(
	ctx contract.Context,
	userAddress loom.Address,
	ownedCards []*orctype.RawCardCollectionCard,
	plasmachainBlockHeight uint64,
) error {
	ctx.Logger().Debug(
		"processOracleCommandResponseGetUserFullCardCollection",
		"userAddress", userAddress.String(),
		"len(ownedCards)", len(ownedCards),
		"plasmachainBlockHeight", plasmachainBlockHeight,
	)

	if plasmachainBlockHeight == 0 {
		return errors.Wrap(errors.New("plasmachainBlockHeight == 0"), "processOracleCommandResponseGetUserFullCardCollection")
	}

	configuration, err := loadContractConfiguration(ctx)
	if err != nil {
		return err
	}

	if configuration.CardCollectionSyncDataVersion == "" {
		return fmt.Errorf("configuration.CardCollectionSyncDataVersion is not set")
	}

	found, userId, err := loadUserIdByAddress(ctx, userAddress)
	if err != nil {
		return err
	}

	if !found {
		ctx.Logger().Info("failed to get user id by address", "address", userAddress.String())
		return nil
	}

	// New collection is equal to default collection + owned cards
	cardCollection, err := loadDefaultCardCollection(ctx, configuration.CardCollectionSyncDataVersion)
	if err != nil {
		return err
	}

	// Add owned cards
	cardAmountChanges := make([]CardAmountChangeItem, len(ownedCards), len(ownedCards))
	for index, rawCard := range ownedCards {
		cardAmountChanges[index] = CardAmountChangeItem{
			CardKey:      cardKeyFromCardTokenId(rawCard.CardTokenId.Value.Int64()),
			AmountChange: rawCard.Amount.Value.Int64(),
		}
	}

	cardCollection.Cards, err =
		z.syncAndSaveCardAmountChangesToCollectionAndUpdateDecks(
			ctx,
			cardCollection.Cards,
			cardAmountChanges,
			configuration.CardCollectionSyncDataVersion,
			userId,
		)

	if err != nil {
		return err
	}

	persistentData, err := loadUserPersistentData(ctx, userId)
	if err != nil {
		return err
	}

	ctx.Logger().Debug(
		"setting LastFullCardCollectionSyncPlasmachainBlockHeight",
		"userAddress", userAddress.String(),
		"userId", userId,
		"plasmachainBlockHeight", plasmachainBlockHeight,
	)
	persistentData.LastFullCardCollectionSyncPlasmachainBlockHeight = plasmachainBlockHeight
	err = saveUserPersistentData(ctx, userId, persistentData)
	if err != nil {
		return err
	}

	// Emit event
	emitMsg := createUserEventBase(userId)
	emitMsg.Event = &zb_calls.UserEvent_FullCardCollectionSync{
		FullCardCollectionSync: &zb_calls.UserEvent_FullCardCollectionSyncEvent{},
	}

	data, err := proto.Marshal(emitMsg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal FullCardCollectionSyncEvent")
	}
	ctx.EmitTopics(data, createUserEventTopic(userId))

	return nil
}

func loadUserCardCollection(ctx contract.StaticContext, version string, userID string) ([]*zb_data.CardCollectionCard, error) {
	configuration, err := loadContractConfiguration(ctx)
	if err != nil {
		return nil, err
	}

	cardLibrary, err := loadCardLibrary(ctx, version)
	if err != nil {
		return nil, err
	}

	var collectionCards []*zb_data.CardCollectionCard
	if !configuration.UseCardLibraryAsUserCollection {
		cardKeyToCardMap, err := getCardKeyToCardMap(cardLibrary.Cards)
		if err != nil {
			return nil, err
		}

		collectionCardsRaw, err := loadUserCardCollectionRaw(ctx, userID)
		if err != nil {
			return nil, err
		}

		// Filter out cards not in card library
		knownCollectionCards := make([]*zb_data.CardCollectionCard, 0)
		for _, collectionCard := range collectionCardsRaw.Cards {
			_, exists := cardKeyToCardMap[collectionCard.CardKey]
			if exists {
				knownCollectionCards = append(knownCollectionCards, collectionCard)
			}
		}

		collectionCards = knownCollectionCards
	} else {
		collectionCards, err = generateFullCardCollection(cardLibrary, true)
		if err != nil {
			return nil, err
		}
	}

	return collectionCards, nil
}

func generateFullCardCollection(cardLibrary *zb_data.CardList, onlyStandardCardVariant bool) ([]*zb_data.CardCollectionCard, error) {
	// Construct fake collection with max count of every standard card
	collectionCards := []*zb_data.CardCollectionCard{}
	for _, card := range cardLibrary.Cards {
		if card.Hidden {
			continue
		}

		if onlyStandardCardVariant && card.CardKey.Variant != zb_enums.CardVariant_Standard {
			continue
		}

		maxAmountOfCardInDeck, err := getMaxAmountOfCardInDeck(card)
		if err != nil {
			return nil, err
		}

		collectionCards = append(collectionCards, &zb_data.CardCollectionCard{
			CardKey: card.CardKey,
			Amount:  int64(maxAmountOfCardInDeck),
		})
	}

	return collectionCards, nil
}

func cardCollectionToCardKeyToAmountMap(collectionCards []*zb_data.CardCollectionCard) cardKeyToAmountMap {
	cardKeyToAmountMap := cardKeyToAmountMap{}
	for _, collectionCard := range collectionCards {
		cardKeyToAmountMap[collectionCard.CardKey] = collectionCard.Amount
	}

	return cardKeyToAmountMap
}

func (z *ZombieBattleground) syncCardAmountChangesToCollection(collectionCards []*zb_data.CardCollectionCard, changes []CardAmountChangeItem) []*zb_data.CardCollectionCard {
	// We are allowing unknown cards to be added.
	// This is to handle the case of user owning a card not existing on gamechain yet.

	cardKeyToAmountMap := cardCollectionToCardKeyToAmountMap(collectionCards)

	// Apply changes
	for _, change := range changes {
		currentAmount, _ := cardKeyToAmountMap[change.CardKey]
		currentAmount += change.AmountChange
		cardKeyToAmountMap[change.CardKey] = currentAmount
	}

	// Convert back to list
	collectionCards = []*zb_data.CardCollectionCard{}
	for cardKey, amount := range cardKeyToAmountMap {
		if amount == 0 {
			continue
		}

		collectionCards = append(collectionCards, &zb_data.CardCollectionCard{
			CardKey: cardKey,
			Amount:  amount,
		})
	}

	sort.SliceStable(collectionCards, func(i, j int) bool {
		return collectionCards[i].CardKey.Compare(&collectionCards[j].CardKey) < 0
	})

	return collectionCards
}

func (z *ZombieBattleground) syncAndSaveCardAmountChangesToCollectionAndUpdateDecks(
	ctx contract.Context,
	collectionCards []*zb_data.CardCollectionCard,
	changes []CardAmountChangeItem,
	version string,
	userId string,
) ([]*zb_data.CardCollectionCard, error) {
	collectionCards = z.syncCardAmountChangesToCollection(collectionCards, changes)

	if debugEnabled {
		fmt.Println("----------------- " + userId)
		for _, card := range collectionCards {
			fmt.Printf("card: %s, amount: %d\n", card.CardKey.String(), card.Amount)
		}
	}

	err := saveUserCardCollection(ctx, userId, &zb_data.CardCollectionList{
		Cards: collectionCards,
	})

	// Load and save decks to limit cards in deck to cards in collection
	deckList, err := loadDecks(ctx, version, userId)
	if err != nil {
		return nil, err
	}

	err = saveDecks(ctx, version, userId, deckList)
	if err != nil {
		return nil, err
	}

	return collectionCards, nil
}

// loads the list of card amount changes, merges with changes, saves the list back
func (z *ZombieBattleground) applyAddressToCardAmountsChangeMapDelta(ctx contract.Context, deltaAddressToCardKeyToAmountChangesMap addressToCardKeyToAmountChangesMap) error {
	// get sorted keys
	deltaAddresses := make([]string, 0, len(deltaAddressToCardKeyToAmountChangesMap))
	for deltaAddress := range deltaAddressToCardKeyToAmountChangesMap {
		deltaAddresses = append(deltaAddresses, deltaAddress)
	}
	sort.SliceStable(deltaAddresses, func(i, j int) bool {
		return strings.Compare(deltaAddresses[i], deltaAddresses[j]) < 0
	})

	for _, deltaAddress := range deltaAddresses {
		deltaCardKeyToAmountChangesMap := deltaAddressToCardKeyToAmountChangesMap[deltaAddress]

		// Get address from string
		deltaLocalAddress, _ := loom.LocalAddressFromHexString("0x" + deltaAddress)
		deltaAddress := loom.Address{
			ChainID: ctx.Message().Sender.ChainID,
			Local:   deltaLocalAddress,
		}

		// Load current list and convert to map
		cardAmountChangesContainer, err := loadPendingCardAmountChangesContainerByAddress(ctx, deltaAddress)
		if err != nil {
			return err
		}

		cardKeyToAmountChangeMap := cardKeyToAmountChangeMap{}
		for _, cardAmountChangeItem := range cardAmountChangesContainer.CardAmountChanges {
			cardKeyToAmountChangeMap[cardAmountChangeItem.CardKey] = int64(cardAmountChangeItem.AmountChange)
		}

		// Update map
		for cardKey, amountChange := range deltaCardKeyToAmountChangesMap {
			currentAmountChange, _ := cardKeyToAmountChangeMap[cardKey]
			cardKeyToAmountChangeMap[cardKey] = currentAmountChange + amountChange
		}

		// Convert map back to list and save
		// get sorted keys
		cardKeys := make([]battleground_proto.CardKey, 0, len(cardKeyToAmountChangeMap))
		for cardKey := range cardKeyToAmountChangeMap {
			cardKeys = append(cardKeys, cardKey)
		}
		sort.SliceStable(cardKeys, func(i, j int) bool {
			return cardKeys[i].Compare(&cardKeys[j]) < 0
		})

		cardAmountChangesContainer.CardAmountChanges = []*zb_data.CardAmountChangeItem{}
		for _, cardKey := range cardKeys {
			amountChange := cardKeyToAmountChangeMap[cardKey]
			cardAmountChangeItem := &zb_data.CardAmountChangeItem{
				CardKey:      cardKey,
				AmountChange: amountChange,
			}
			cardAmountChangesContainer.CardAmountChanges = append(cardAmountChangesContainer.CardAmountChanges, cardAmountChangeItem)
		}

		err = savePendingCardAmountChangeItemsContainerByAddress(ctx, deltaAddress, cardAmountChangesContainer)
		if err != nil {
			return err
		}

		// Try to apply the changes immediately
		userIdFound, err := z.applyPendingCardAmountChanges(ctx, deltaAddress)
		if err != nil {
			return err
		}

		if !userIdFound {
			ctx.Logger().Debug("failed to apply pending card amount changes, user id not found by address", "userAddress", deltaAddress.String())
		}
	}
	return nil
}

func (z *ZombieBattleground) applyPendingCardAmountChanges(ctx contract.Context, address loom.Address) (userIdFound bool, err error) {
	userIdFound, userId, err := loadUserIdByAddress(ctx, address)
	if err != nil {
		return userIdFound, err
	}

	if !userIdFound {
		return userIdFound, nil
	}

	configuration, err := loadContractConfiguration(ctx)
	if err != nil {
		return userIdFound, err
	}

	if configuration.CardCollectionSyncDataVersion == "" {
		return userIdFound, fmt.Errorf("configuration.CardCollectionSyncDataVersion is not set")
	}

	cardAmountChangesContainer, err := loadPendingCardAmountChangesContainerByAddress(ctx, address)
	if err != nil {
		return userIdFound, err
	}

	if len(cardAmountChangesContainer.CardAmountChanges) == 0 {
		return userIdFound, nil
	}

	// Load collection
	userCardCollection, err := loadUserCardCollectionRaw(ctx, userId)
	if err != nil {
		return userIdFound, err
	}

	// Add owned cards
	cardAmountChanges := make([]CardAmountChangeItem, len(cardAmountChangesContainer.CardAmountChanges), len(cardAmountChangesContainer.CardAmountChanges))
	for index, cardAmountChange := range cardAmountChangesContainer.CardAmountChanges {
		cardAmountChanges[index] = CardAmountChangeItem{
			CardKey:      cardAmountChange.CardKey,
			AmountChange: cardAmountChange.AmountChange,
		}
	}

	userCardCollection.Cards, err =
		z.syncAndSaveCardAmountChangesToCollectionAndUpdateDecks(
			ctx,
			userCardCollection.Cards,
			cardAmountChanges,
			configuration.CardCollectionSyncDataVersion,
			userId,
		)

	if err != nil {
		return userIdFound, err
	}

	// Clear pending changes
	cardAmountChangesContainer.CardAmountChanges = []*zb_data.CardAmountChangeItem{}
	err = savePendingCardAmountChangeItemsContainerByAddress(ctx, address, cardAmountChangesContainer)
	if err != nil {
		return userIdFound, err
	}

	// Emit event
	emitMsg := createUserEventBase(userId)
	emitMsg.Event = &zb_calls.UserEvent_AutoCardCollectionSync{
		AutoCardCollectionSync: &zb_calls.UserEvent_AutoCardCollectionSyncEvent{},
	}

	data, err := proto.Marshal(emitMsg)
	if err != nil {
		return userIdFound, errors.Wrap(err, "failed to marshal AutoCardCollectionSyncEvent")
	}
	ctx.EmitTopics(data, createUserEventTopic(userId))

	return userIdFound, nil
}

func (z *ZombieBattleground) updateCardAmountChangeToAddressToCardAmountsChangeMap(
	addressToCardAmountsChangeMap addressToCardKeyToAmountChangesMap,
	from common.LocalAddress,
	to common.LocalAddress,
	cardKey battleground_proto.CardKey,
	amount uint,
	zbgCardContractAddress loom.Address,
) error {
	// From
	fromHex := hex.EncodeToString(from)
	if fromHex != hex.EncodeToString(zbgCardContractAddress.Local) {
		fromCardAmountsChangeMap, exists := addressToCardAmountsChangeMap[fromHex]
		if !exists {
			fromCardAmountsChangeMap = cardKeyToAmountChangeMap{}
			addressToCardAmountsChangeMap[fromHex] = fromCardAmountsChangeMap
		}

		fromAmountValue, _ := fromCardAmountsChangeMap[cardKey]
		fromAmountValue = fromAmountValue - int64(amount)
		fromCardAmountsChangeMap[cardKey] = fromAmountValue
	}

	// To
	toHex := hex.EncodeToString(to)
	if toHex != hex.EncodeToString(zbgCardContractAddress.Local) {
		toCardAmountsChangeMap, exists := addressToCardAmountsChangeMap[toHex]
		if !exists {
			toCardAmountsChangeMap = cardKeyToAmountChangeMap{}
			addressToCardAmountsChangeMap[toHex] = toCardAmountsChangeMap
		}

		toAmountValue, _ := toCardAmountsChangeMap[cardKey]
		toAmountValue = toAmountValue + int64(amount)
		toCardAmountsChangeMap[cardKey] = toAmountValue
	}

	return nil
}

func fixDeckListCards(ctx contract.Context, deckList *zb_data.DeckList, version string, userID string) (changed bool, err error) {
	cardLibrary, err := loadCardLibrary(ctx, version)
	if err != nil {
		return false, errors.Wrap(err, "error fixing deck list")
	}

	cardKeyToCardMap, err := getCardKeyToCardMap(cardLibrary.Cards)
	if err != nil {
		return false, errors.Wrap(err, "error fixing deck list")
	}

	collectionCards, err := loadUserCardCollection(ctx, version, userID)
	if err != nil {
		return false, errors.Wrap(err, "error fixing deck list")
	}

	changed = false
	for _, deck := range deckList.Decks {
		if fixDeckCardVariants(deck, cardKeyToCardMap) {
			changed = true
		}

		if limitDeckByCardCollection(deck, collectionCards) {
			changed = true
		}
	}

	return changed, nil
}

func limitDeckByCardCollection(deck *zb_data.Deck, collectionCards []*zb_data.CardCollectionCard) (changed bool) {
	changed = false
	collectionCardKeyToAmountMap := cardCollectionToCardKeyToAmountMap(collectionCards)
	var newDeckCards = make([]*zb_data.DeckCard, 0)

	for _, deckCard := range deck.Cards {
		amountInCollection, existsInCollection := collectionCardKeyToAmountMap[deckCard.CardKey]
		if !existsInCollection {
			changed = true
			continue
		}

		if deckCard.Amount > amountInCollection {
			deckCard.Amount = amountInCollection
			changed = true
		}

		newDeckCards = append(newDeckCards, deckCard)
	}

	deck.Cards = newDeckCards
	return changed
}

func fixDeckCardVariants(deck *zb_data.Deck, cardKeyToCardMap map[battleground_proto.CardKey]*zb_data.Card) (changed bool) {
	var newDeckCards = make([]*zb_data.DeckCard, 0)
	var cardKeyToDeckCard = make(map[battleground_proto.CardKey]*zb_data.DeckCard)
	for _, deckCard := range deck.Cards {
		cardKeyToDeckCard[deckCard.CardKey] = deckCard
	}

	for _, deckCard := range deck.Cards {
		// Check if this specific variant of a card exists in card library
		_, variantExists := cardKeyToCardMap[deckCard.CardKey]
		if !variantExists {
			// If this variant is not in card library, try to fallback to Normal variant
			normalVariantCardKey := battleground_proto.CardKey{
				MouldId: deckCard.CardKey.MouldId,
				Variant: zb_enums.CardVariant_Standard,
			}

			_, normalVariantExists := cardKeyToCardMap[normalVariantCardKey]

			// If normal variant doesn't exist in card library too, just remove the card from the deck completely
			if !normalVariantExists {
				changed = true
			} else {
				normalVariantDeckCard, normalVariantDeckCardExists := cardKeyToDeckCard[normalVariantCardKey]
				// If normal variant exists in card library AND in the deck,
				// add the amount of the special variant to normal variant
				if normalVariantDeckCardExists {
					normalVariantDeckCard.Amount += deckCard.Amount
					changed = true
				} else {
					// If normal variant exists in card library, but NOT in the deck,
					// create a normal variant deck card and add special variant amount to it
					normalVariantDeckCard = &zb_data.DeckCard{
						CardKey: normalVariantCardKey,
						Amount:  deckCard.Amount,
					}

					newDeckCards = append(newDeckCards, normalVariantDeckCard)
					cardKeyToDeckCard[normalVariantCardKey] = normalVariantDeckCard
					changed = true
				}
			}
		} else {
			newDeckCards = append(newDeckCards, deckCard)
		}
	}

	deck.Cards = newDeckCards
	return changed
}

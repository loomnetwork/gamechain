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
	blockHeight uint64,
) error {
	ctx.Logger().Debug(
		"processOracleCommandResponseGetUserFullCardCollection",
		"userAddress", userAddress.String(),
		"len(ownedCards)", len(ownedCards),
	)

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

	// Take default collection and add owned cards to it
	cardCollection, err := loadDefaultCardCollection(ctx, configuration.CardCollectionSyncDataVersion)
	if err != nil {
		return err
	}

	for _, rawCard := range ownedCards {
		cardCollection.Cards = z.syncCardToCollection(
			ctx,
			userId,
			cardKeyFromCardTokenId(rawCard.CardTokenId.Value.Int64()),
			rawCard.Amount.Value.Int64(),
			cardCollection.Cards,
		)
	}

	sort.SliceStable(cardCollection.Cards, func(i, j int) bool {
		return cardCollection.Cards[i].CardKey.Compare(&cardCollection.Cards[j].CardKey) < 0
	})

	err = saveUserCardCollection(ctx, userId, cardCollection)
	if debugEnabled {
		fmt.Println("-----------------")
		for _, card := range cardCollection.Cards {
			fmt.Printf("card: %s, amount: %d\n", card.CardKey.String(), card.Amount)
		}
	}

	if err != nil {
		return err
	}

	persistentData, err := loadUserPersistentData(ctx, userId)
	if err != nil {
		return err
	}

	persistentData.LastFullCardCollectionSyncPlasmachainBlockHeight = blockHeight
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
	ctx.EmitTopics([]byte(data), createUserEventTopic(userId))

	return nil
}

func (z *ZombieBattleground) loadUserCardCollection(ctx contract.StaticContext, version string, userID string) ([]*zb_data.CardCollectionCard, error) {
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

		// Filter our cards not in card library
		knownCollectionCards := make([]*zb_data.CardCollectionCard, 0)
		for _, collectionCard := range collectionCardsRaw.Cards {
			_, exists := cardKeyToCardMap[collectionCard.CardKey]
			if exists {
				knownCollectionCards = append(knownCollectionCards, collectionCard)
			}
		}

		collectionCards = knownCollectionCards
	} else {
		// Construct fake collection with max count of every card
		for _, card := range cardLibrary.Cards {
			if card.Hidden {
				continue
			}

			if card.CardKey.Variant != zb_enums.CardVariant_Standard {
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
	}

	return collectionCards, nil
}

func (z *ZombieBattleground) syncCardToCollection(ctx contract.Context, userID string, cardKey battleground_proto.CardKey, amount int64, collectionCards []*zb_data.CardCollectionCard) []*zb_data.CardCollectionCard {
	// We are allowing unknown cards to be added.
	// This is to handle the case of user buying a card not existing on gamechain yet.

	// add to collection
	found := false
	for _, collectionCard := range collectionCards {
		if collectionCard.CardKey == cardKey {
			collectionCard.Amount += amount
			found = true
			break
		}
	}
	if !found {
		collectionCards = append(collectionCards, &zb_data.CardCollectionCard{
			CardKey: cardKey,
			Amount:  amount,
		})
	}

	return collectionCards
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

	// Apply changes
	for _, cardAmountChange := range cardAmountChangesContainer.CardAmountChanges {
		userCardCollection.Cards = z.syncCardToCollection(
			ctx,
			userId,
			cardAmountChange.CardKey,
			cardAmountChange.AmountChange,
			userCardCollection.Cards,
		)
	}

	// Save collection
	err = saveUserCardCollection(ctx, userId, userCardCollection)
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
	ctx.EmitTopics([]byte(data), createUserEventTopic(userId))

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

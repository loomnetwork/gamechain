package battleground

import (
	"encoding/hex"
	"fmt"
	battleground_proto "github.com/loomnetwork/gamechain/battleground/proto"
	orctype "github.com/loomnetwork/gamechain/types/oracle"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	loom "github.com/loomnetwork/go-loom"
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

	if configuration.CardSyncDataVersion == "" {
		return fmt.Errorf("configuration.CardSyncDataVersion is not set")
	}

	userId, err := loadUserIdByAddress(ctx, userAddress)
	if err != nil {
		return err
	}

	// take default collection and add owned cards
	cardCollection, err := loadDefaultCardCollection(ctx, configuration.CardSyncDataVersion)
	if err != nil {
		return errors.Wrap(err, "error initializing user default card collection")
	}

	for _, rawCard := range ownedCards {
		cardCollection.Cards = z.syncCardToCollection(ctx, userId, rawCard.CardTokenId.Value.Int64(), rawCard.Amount.Value.Int64(), cardCollection.Cards)
	}

	err = saveUserCardCollection(ctx, userId, cardCollection)
	fmt.Println("-----------------")
	for _, card := range cardCollection.Cards {
		fmt.Printf("card: %s, amount: %d\n", card.CardKey.String(), card.Amount)
	}

	if err != nil {
		return err
	}

	persistentData, err := loadUserPersistentData(ctx, userId)
	if err != nil {
		return err
	}

	persistentData.LastFullCardSyncPlasmachainBlockHeight = blockHeight
	err = saveUserPersistentData(ctx, userId, persistentData)
	if err != nil {
		return err
	}

	return nil
}

func (z *ZombieBattleground) loadUserCardCollection(ctx contract.StaticContext, version string, userID string) ([]*zb_data.CardCollectionCard, error) {
	configuration, err := loadContractConfiguration(ctx)
	if err != nil {
		return nil, err
	}

	var cards []*zb_data.CardCollectionCard
	if !configuration.UseCardLibraryAsUserCollection {
		userCollectionCards, err := loadUserCardCollectionRaw(ctx, userID)
		if err != nil {
			return nil, err
		}

		cards = userCollectionCards.Cards
	} else {
		// Construct fake collection with max count of every card
		cardLibrary, err := loadCardLibrary(ctx, version)
		if err != nil {
			return nil, err
		}

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

			cards = append(cards, &zb_data.CardCollectionCard{
				CardKey: card.CardKey,
				Amount: int64(maxAmountOfCardInDeck),
			})
		}
	}

	return cards, nil
}

func (z *ZombieBattleground) syncCardToCollection(ctx contract.Context, userID string, cardTokenId int64, amount int64, cardCollection []*zb_data.CardCollectionCard) []*zb_data.CardCollectionCard {
	// Map from cardTokenId to mouldID
	// the formular is cardTokenId = mouldID + x
	// for example cardTokenId 250 = 25 + 0
	//   or 161 = 16 + 1
	cardKey := cardKeyFromCardTokenId(cardTokenId)

	// We are allowing unknown cards to be added.
	// This is to handle the case user buying a card not existing on gamechain yet.

	// add to collection
	found := false
	for i := range cardCollection {
		if cardCollection[i].CardKey == cardKey {
			cardCollection[i].Amount += amount
			found = true
			break
		}
	}
	if !found {
		cardCollection = append(cardCollection, &zb_data.CardCollectionCard{
			CardKey: cardKey,
			Amount:  amount,
		})
	}

	return cardCollection
}

/*func (z *ZombieBattleground) syncCardToCollection(ctx contract.Context, userID string, cardTokenId int64, amount int64, version string) error {
	if err := z.validateOracle(ctx); err != nil {
		return err
	}

	cardCollection, err := loadUserCardCollectionRaw(ctx, userID)
	if err != nil {
		return err
	}

	// Map from cardTokenId to mouldID
	// the formular is cardTokenId = mouldID + x
	// for example cardTokenId 250 = 25 + 0
	//   or 161 = 16 + 1
	cardKey := cardKeyFromCardTokenId(cardTokenId)

	// We are allowing unknown cards to be added.
	// This is to handle the case user buying a card not existing on gamechain yet.

	// add to collection
	found := false
	for i := range cardCollection.Cards {
		if cardCollection.Cards[i].CardKey == cardKey {
			cardCollection.Cards[i].Amount += amount
			found = true
			break
		}
	}
	if !found {
		cardCollection.Cards = append(cardCollection.Cards, &zb_data.CardCollectionCard{
			CardKey: cardKey,
			Amount:  amount,
		})
	}
	return saveUserCardCollection(ctx, userID, cardCollection)
}*/

// loads the list of card amount changes, merges with changes, saves the list back
func (z *ZombieBattleground) saveAddressToCardAmountsChangeMapDelta(ctx contract.Context, deltaAddressToCardKeyToAmountChangesMap addressToCardKeyToAmountChangesMap) error {
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
		for _, cardAmountChangeItem := range cardAmountChangesContainer.CardAmountChange {
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

		cardAmountChangesContainer.CardAmountChange = []*zb_data.CardAmountChangeItem{}
		for _, cardKey := range cardKeys {
			amountChange := cardKeyToAmountChangeMap[cardKey]
			cardAmountChangeItem := &zb_data.CardAmountChangeItem{
				CardKey: cardKey,
				AmountChange: amountChange,
			}
			cardAmountChangesContainer.CardAmountChange = append(cardAmountChangesContainer.CardAmountChange, cardAmountChangeItem)
		}

		err = savePendingCardAmountChangeItemsContainerByAddress(ctx, deltaAddress, cardAmountChangesContainer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (z *ZombieBattleground) handleCardAmountChange(
	addressToCardAmountsChangeMap addressToCardKeyToAmountChangesMap,
	from common.LocalAddress,
	to common.LocalAddress,
	cardKey battleground_proto.CardKey,
	amount uint,
) error {
	fromHex := hex.EncodeToString(from)
	toHex := hex.EncodeToString(to)

	fromCardAmountsChangeMap, exists := addressToCardAmountsChangeMap[fromHex]
	if !exists {
		fromCardAmountsChangeMap = cardKeyToAmountChangeMap{}
		addressToCardAmountsChangeMap[fromHex] = fromCardAmountsChangeMap
	}

	toCardAmountsChangeMap, exists := addressToCardAmountsChangeMap[toHex]
	if !exists {
		toCardAmountsChangeMap = cardKeyToAmountChangeMap{}
		addressToCardAmountsChangeMap[toHex] = toCardAmountsChangeMap
	}

	fromAmountValue, _ := fromCardAmountsChangeMap[cardKey]
	toAmountValue, _ := toCardAmountsChangeMap[cardKey]

	fromAmountValue = fromAmountValue - int64(amount)
	toAmountValue = toAmountValue + int64(amount)

	fromCardAmountsChangeMap[cardKey] = fromAmountValue
	toCardAmountsChangeMap[cardKey] = toAmountValue

	return nil
}

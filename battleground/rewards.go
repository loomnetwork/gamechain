package battleground

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	battleground_proto "github.com/loomnetwork/gamechain/battleground/proto"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/common"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/pkg/errors"
	"math/big"
	"sort"
	"strings"
)

type MintingContext struct {
	context               contract.Context
	contractConfiguration *zb_data.ContractConfiguration
	contractState         *zb_data.ContractState
	generator             *MintingReceiptGenerator
}

func NewMintingContext(ctx contract.Context) (MintingContext, error) {
	mintingContext := MintingContext{
		context: ctx,
	}
	gatewayPrivateKey, err := crypto.HexToECDSA(purchaseGatewayPrivateKeyHexString)
	if err != nil {
		err = errors.Wrap(err, "error getting private key")
		return MintingContext{}, err
	}

	mintingContext.contractConfiguration, err = loadContractConfiguration(ctx)
	if err != nil {
		return MintingContext{}, err
	}

	generator, err := NewMintingReceiptGenerator(gatewayPrivateKey, uint(mintingContext.contractConfiguration.FiatPurchaseContractVersion))
	if err != nil {
		return MintingContext{}, err
	}
	mintingContext.generator = &generator

	mintingContext.contractState, err = loadContractState(ctx)
	if err != nil {
		return MintingContext{}, err
	}

	return mintingContext, nil
}

func (m *MintingContext) MintBoosterPacksReceipt(userId *big.Int, amount uint) (*MintingReceipt, error) {
	txId := m.contractState.CurrentFiatPurchaseTxId.Value.Int
	receipt, err := m.generator.CreateBoosterReceipt(userId, amount, new(big.Int).SetBytes(txId.Bytes()))
	if err != nil {
		return nil, err
	}

	m.contractState.CurrentFiatPurchaseTxId.Value.Int.Add(m.contractState.CurrentFiatPurchaseTxId.Value.Int, big.NewInt(1))
	return receipt, nil
}

func (m *MintingContext) SaveState() error {
	err := saveContractState(m.context, m.contractState)
	if err != nil {
		return err
	}

	return nil
}

func mintBoosterPacksAndSave(ctx contract.Context, userId string, userIdInt *big.Int, amount uint) (*MintingReceipt, error) {
	// Create minting receipt
	mintingContext, err := NewMintingContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create minting context")
	}

	if userIdInt == nil {
		userIdInt = parseUserIdToNumber(userId)
	}
	receipt, err := mintingContext.MintBoosterPacksReceipt(userIdInt, amount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to mint")
	}

	err = mintingContext.SaveState()
	if err != nil {
		return nil, err
	}

	// Add to the list of pending receipts
	receiptCollection, err := loadPendingMintingTransactionReceipts(ctx, userId)
	if err != nil {
		return nil, err
	}

	receiptCollection.Receipts = append(receiptCollection.Receipts, receipt.MarshalPB())
	err = savePendingMintingTransactionReceipts(ctx, userId, receiptCollection)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

func (z *ZombieBattleground) syncCardToCollection(ctx contract.Context, userID string, cardTokenId int64, amount int64, version string) error {
	if err := z.validateOracle(ctx); err != nil {
		return err
	}

	cardCollection, err := loadCardCollectionRaw(ctx, userID)
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
	return saveCardCollection(ctx, userID, cardCollection)
}

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

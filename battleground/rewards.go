package battleground

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/pkg/errors"
	"math/big"
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
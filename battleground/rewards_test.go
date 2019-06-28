package battleground

import (
	"github.com/eosspark/eos-go/common/hexutil"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/loomnetwork/gamechain/tools/battleground_utility"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	assert "github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestRewardMinting(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	userId := defaultUserIdPrefix + "373"
	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  userId,
		Image:   "PathToImage",
		Version: "v1",
	}, t)

	t.Run("UpdateContractConfiguration should succeed", func(t *testing.T) {
		request := zb_calls.UpdateContractConfigurationRequest{
			SetFiatPurchaseContractVersion: true,
			FiatPurchaseContractVersion:    3,
			SetInitialFiatPurchaseTxId:     true,
			InitialFiatPurchaseTxId:        battleground_utility.MarshalBigIntProto(big.NewInt(100)),
		}

		err := c.UpdateContractConfiguration(ctx, &request)
		assert.Nil(t, err)
	})

	var prevTxId *big.Int
	t.Run("DebugMintBoosterPackReceipt should create receipt", func(t *testing.T) {
		req := &zb_calls.DebugMintBoosterPackReceiptRequest{
			UserId:        battleground_utility.MarshalBigIntProto(big.NewInt(373)),
			BoosterAmount: 3,
		}
		receipt, err := c.DebugMintBoosterPackReceipt(ctx, req)
		assert.Nil(t, err)
		assert.NotNil(t, receipt)
		assert.Equal(t, int64(373), receipt.Receipt.UserId.Value.Int.Int64())
		assert.Equal(t, uint64(3), receipt.Receipt.Booster)
		assert.Equal(t, 65, len(receipt.Receipt.VerifyHash.Signature))

		prevTxId = receipt.Receipt.TxId.Value.Int

		// Validate Solidity-like
		uint256Type, _ := abi.NewType("uint256")
		uint256Arr11Type, _ := abi.NewType("uint256[11]")
		arguments := abi.Arguments{
			{
				Type: uint256Type,
			},
			{
				Type: uint256Arr11Type,
			},
			{
				Type: uint256Type,
			},
			{
				Type: uint256Type,
			},
		}

		bytes, err := arguments.Pack(
			receipt.Receipt.UserId.Value.Int,
			[11]*big.Int{
				big.NewInt(int64(receipt.Receipt.Booster)),
				big.NewInt(int64(receipt.Receipt.Super)),
				big.NewInt(int64(receipt.Receipt.Air)),
				big.NewInt(int64(receipt.Receipt.Earth)),
				big.NewInt(int64(receipt.Receipt.Fire)),
				big.NewInt(int64(receipt.Receipt.Life)),
				big.NewInt(int64(receipt.Receipt.Toxic)),
				big.NewInt(int64(receipt.Receipt.Water)),
				big.NewInt(int64(receipt.Receipt.Small)),
				big.NewInt(int64(receipt.Receipt.Minion)),
				big.NewInt(int64(receipt.Receipt.Binance)),
			},
			receipt.Receipt.TxId.Value.Int,
			big.NewInt(3),
		)

		assert.Nil(t, err)

		var hashBytes []byte
		hash := sha3.NewKeccak256()
		hash.Write(bytes)
		hashBytes = hash.Sum(hashBytes)

		assert.Equal(t, hexutil.Encode(receipt.Receipt.VerifyHash.Hash), hexutil.Encode(hashBytes))
	})

	t.Run("GetPendingMintingTransactionReceipts should contain the receipt", func(t *testing.T) {
		req := &zb_calls.GetPendingMintingTransactionReceiptsRequest{
			UserId: userId,
		}
		receipts, err := c.GetPendingMintingTransactionReceipts(ctx, req)
		assert.Nil(t, err)
		assert.NotNil(t, receipts)
		assert.Equal(t, 1, len(receipts.ReceiptCollection.Receipts))
		assert.Equal(t, int64(373), receipts.ReceiptCollection.Receipts[0].UserId.Value.Int.Int64())
		assert.Equal(t, uint64(3), receipts.ReceiptCollection.Receipts[0].Booster)
		assert.True(t, receipts.ReceiptCollection.Receipts[0].TxId.Value.Int.Cmp(prevTxId) == 0)
	})

	t.Run("ConfirmPendingMintingTransactionReceipt should succeed", func(t *testing.T) {
		req := &zb_calls.ConfirmPendingMintingTransactionReceiptRequest{
			UserId: userId,
			TxId:   battleground_utility.MarshalBigIntProto(prevTxId),
		}
		err := c.ConfirmPendingMintingTransactionReceipt(ctx, req)
		assert.Nil(t, err)
	})

	t.Run("Receipt should be absent from pending receipts list", func(t *testing.T) {
		req := &zb_calls.GetPendingMintingTransactionReceiptsRequest{
			UserId: userId,
		}
		receipts, err := c.GetPendingMintingTransactionReceipts(ctx, req)
		assert.Nil(t, err)
		assert.NotNil(t, receipts)
		assert.Equal(t, 0, len(receipts.ReceiptCollection.Receipts))
	})

	t.Run("Confirming already confirmed receipt should fail", func(t *testing.T) {
		req := &zb_calls.ConfirmPendingMintingTransactionReceiptRequest{
			UserId: userId,
			TxId:   battleground_utility.MarshalBigIntProto(prevTxId),
		}
		err := c.ConfirmPendingMintingTransactionReceipt(ctx, req)
		assert.NotNil(t, err)
	})
}

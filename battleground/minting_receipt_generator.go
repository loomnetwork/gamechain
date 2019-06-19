package battleground

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/go-loom/common"
	"github.com/loomnetwork/go-loom/types"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"github.com/pkg/errors"
	"math/big"
	"strconv"
)

type MintingReceiptGenerator struct {
	gatewayPrivateKey *ecdsa.PrivateKey
	contractVersion   uint
}

type VerifySignResult struct {
	Hash      string `json:"hash"`
	Signature string `json:"signature"`
}

type MintingReceipt struct {
	VerifyHash *VerifySignResult
	UserId     *big.Int
	Booster    uint
	Air        uint
	Earth      uint
	Fire       uint
	Life       uint
	Toxic      uint
	Water      uint
	Super      uint
	Small      uint
	Minion     uint
	Binance    uint
	TxID       *big.Int
}

func NewMintingReceiptGenerator(gatewayPrivateKey *ecdsa.PrivateKey, contractVersion uint) (MintingReceiptGenerator, error) {
	if gatewayPrivateKey == nil {
		return MintingReceiptGenerator{}, fmt.Errorf("private key is nil")
	}

	return MintingReceiptGenerator{
		gatewayPrivateKey: gatewayPrivateKey,
		contractVersion:   contractVersion,
	}, nil
}

func (generator *MintingReceiptGenerator) CreateBoosterReceipt(userId *big.Int, boosterAmount uint, txId *big.Int) (*MintingReceipt, error) {
	verifyHash, err := generator.generateEosVerifySignResult(
		userId,
		generator.gatewayPrivateKey,
		boosterAmount,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		txId,
		generator.contractVersion)

	if err != nil {
		err = errors.Wrap(err, "error generating reward hash")
		return nil, err
	}

	response := MintingReceipt{
		VerifyHash: verifyHash,
		UserId:     userId,
		Booster:    boosterAmount,
		Air:        0,
		Earth:      0,
		Fire:       0,
		Life:       0,
		Toxic:      0,
		Water:      0,
		Super:      0,
		Small:      0,
		Minion:     0,
		Binance:    0,
		TxID:       txId}

	return &response, nil
}

func (generator *MintingReceiptGenerator) generateEosVerifySignResult(
	userId *big.Int,
	privateKey *ecdsa.PrivateKey,
	booster uint,
	air uint,
	earth uint,
	fire uint,
	life uint,
	toxic uint,
	water uint,
	super uint,
	small uint,
	minion uint,
	binance uint,
	txID *big.Int,
	contractVersion uint) (*VerifySignResult, error) {

	hash, err := generator.createEosHash(userId, booster, air, earth, fire, life, toxic, water, super, small, minion, binance, txID, contractVersion)

	if err != nil {
		return nil, err
	}

	sig, err := generator.soliditySign(hash, privateKey)

	if err != nil {
		return nil, err
	}

	return &VerifySignResult{
		Hash:      "0x" + hex.EncodeToString(hash),
		Signature: "0x" + hex.EncodeToString(sig),
	}, nil
}

func (generator *MintingReceiptGenerator) createEosHash(
	userID *big.Int,
	booster uint,
	air uint,
	earth uint,
	fire uint,
	life uint,
	toxic uint,
	water uint,
	super uint,
	small uint,
	minion uint,
	binance uint,
	txID *big.Int,
	contractVersion uint) ([]byte, error) {

	hash := solsha3.SoliditySHA3(
		solsha3.Uint256(userID.String()),
		solsha3.Uint256(strconv.FormatUint(uint64(booster), 10)),
		solsha3.Uint256(strconv.FormatUint(uint64(super), 10)),
		solsha3.Uint256(strconv.FormatUint(uint64(air), 10)),
		solsha3.Uint256(strconv.FormatUint(uint64(earth), 10)),
		solsha3.Uint256(strconv.FormatUint(uint64(fire), 10)),
		solsha3.Uint256(strconv.FormatUint(uint64(life), 10)),
		solsha3.Uint256(strconv.FormatUint(uint64(toxic), 10)),
		solsha3.Uint256(strconv.FormatUint(uint64(water), 10)),
		solsha3.Uint256(strconv.FormatUint(uint64(small), 10)),
		solsha3.Uint256(strconv.FormatUint(uint64(minion), 10)),
		solsha3.Uint256(strconv.FormatUint(uint64(binance), 10)),
		solsha3.Uint256(txID.String()),
		solsha3.Uint256(strconv.FormatUint(uint64(contractVersion), 10)),
	)

	if len(hash) == 0 {
		return nil, errors.New("failed to generate hash")
	}

	return hash, nil
}

func (generator *MintingReceiptGenerator) soliditySign(data []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := crypto.Sign(data, privateKey)
	if err != nil {
		return nil, err
	}

	v := sig[len(sig)-1]
	sig[len(sig)-1] = v + 27
	return sig, nil
}

func (t *MintingReceipt) MarshalPB() *zb_data.MintingTransactionReceipt {
	return &zb_data.MintingTransactionReceipt{
		VerifyHash: &zb_data.MintingTransactionReceipt_VerifySignResult{
			Hash:      hexutil.MustDecode(t.VerifyHash.Hash),
			Signature: hexutil.MustDecode(t.VerifyHash.Signature),
		},
		UserId:  &types.BigUInt{Value: common.BigUInt{Int: t.UserId}},
		Booster: uint64(t.Booster),
		Air:     uint64(t.Air),
		Earth:   uint64(t.Earth),
		Fire:    uint64(t.Fire),
		Life:    uint64(t.Life),
		Toxic:   uint64(t.Toxic),
		Water:   uint64(t.Water),
		Super:   uint64(t.Super),
		Small:   uint64(t.Small),
		Minion:  uint64(t.Minion),
		Binance: uint64(t.Binance),
		TxId:    &types.BigUInt{Value: common.BigUInt{Int: t.TxID}},
	}
}

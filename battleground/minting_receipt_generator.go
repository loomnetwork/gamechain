package battleground

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/loomnetwork/gamechain/tools/battleground_utility"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
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
	Super      uint
	Air        uint
	Earth      uint
	Fire       uint
	Life       uint
	Toxic      uint
	Water      uint
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

func (generator *MintingReceiptGenerator) CreateGenericPackReceipt(
	userId *big.Int,
	booster uint,
	super uint,
	air uint,
	earth uint,
	fire uint,
	life uint,
	toxic uint,
	water uint,
	small uint,
	minion uint,
	binance uint,
	txId *big.Int,
) (*MintingReceipt, error) {
	verifyHash, err := generator.generateVerifySignResult(
		userId,
		generator.gatewayPrivateKey,
		booster,
		super,
		air,
		earth,
		fire,
		life,
		toxic,
		water,
		small,
		minion,
		binance,
		txId,
		generator.contractVersion)

	if err != nil {
		err = errors.Wrap(err, "error generating reward hash")
		return nil, err
	}

	response := MintingReceipt{
		VerifyHash: verifyHash,
		UserId:     userId,
		Booster:    booster,
		Super:      super,
		Air:        air,
		Earth:      earth,
		Fire:       fire,
		Life:       life,
		Toxic:      toxic,
		Water:      water,
		Small:      small,
		Minion:     minion,
		Binance:    binance,
		TxID:       txId}

	return &response, nil
}

func (generator *MintingReceiptGenerator) generateVerifySignResult(
	userId *big.Int,
	privateKey *ecdsa.PrivateKey,
	booster uint,
	super uint,
	air uint,
	earth uint,
	fire uint,
	life uint,
	toxic uint,
	water uint,
	small uint,
	minion uint,
	binance uint,
	txID *big.Int,
	contractVersion uint) (*VerifySignResult, error) {

	hash, err := generator.createHash(userId, booster, super, air, earth, fire, life, toxic, water, small, minion, binance, txID, contractVersion)

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

func (generator *MintingReceiptGenerator) createHash(
	userID *big.Int,
	booster uint,
	super uint,
	air uint,
	earth uint,
	fire uint,
	life uint,
	toxic uint,
	water uint,
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
		UserId:  battleground_utility.MarshalBigIntProto(t.UserId),
		Booster: uint64(t.Booster),
		Super:   uint64(t.Super),
		Air:     uint64(t.Air),
		Earth:   uint64(t.Earth),
		Fire:    uint64(t.Fire),
		Life:    uint64(t.Life),
		Toxic:   uint64(t.Toxic),
		Water:   uint64(t.Water),
		Small:   uint64(t.Small),
		Minion:  uint64(t.Minion),
		Binance: uint64(t.Binance),
		TxId:    battleground_utility.MarshalBigIntProto(t.TxID),
	}
}

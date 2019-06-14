package battleground

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"github.com/pkg/errors"
	"math/big"
	"strconv"
)

const contractVersion = 3

type ReceiptGenerator struct {
	gatewayPrivateKey *ecdsa.PrivateKey
}

type VerifySignResult struct {
	Hash      string `json:"hash"`
	Signature string `json:"signature"`
}

type TransactionResponse struct {
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

func NewReceiptGenerator(gatewayPrivateKey *ecdsa.PrivateKey) (ReceiptGenerator, error) {
	if gatewayPrivateKey == nil {
		return ReceiptGenerator{}, fmt.Errorf("private key is nil")
	}

	return ReceiptGenerator{
		gatewayPrivateKey: gatewayPrivateKey,
	}, nil
}

func (generator *ReceiptGenerator) CalculateAbsoluteTxId(relativeTxId *big.Int) *big.Int {
	txId := big.NewInt(1)
	txId = txId.Lsh(txId, 127)
	return txId.Add(txId, relativeTxId)
}

func (generator *ReceiptGenerator) CreateBoosterReceipt(userId *big.Int, boosterAmount uint, txId *big.Int) (*TransactionResponse, error) {
	verifyHash, err := generator.generateVerifyEosHash(
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
		contractVersion)

	if err != nil {
		err = errors.Wrap(err, "error generating reward hash")
		return nil, err
	}

	response := TransactionResponse{
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

func (generator *ReceiptGenerator) generateVerifyEosHash(
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

func (generator *ReceiptGenerator) createEosHash(
	UserID *big.Int,
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
		solsha3.Uint256(UserID.String()),
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

func (generator *ReceiptGenerator) soliditySign(data []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := crypto.Sign(data, privateKey)
	if err != nil {
		return nil, err
	}

	v := sig[len(sig)-1]
	sig[len(sig)-1] = v + 27
	return sig, nil
}
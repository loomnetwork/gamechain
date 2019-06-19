package battleground

import (
	"encoding/hex"
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	battleground_proto "github.com/loomnetwork/gamechain/battleground/proto"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	"github.com/loomnetwork/go-loom/common"
	"github.com/loomnetwork/go-loom/types"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	assert "github.com/stretchr/testify/require"
)

var initRequest = zb_calls.InitRequest {
}

var updateInitRequest = zb_calls.UpdateInitRequest {
}

func readJsonFileToProtobuf(filename string, message proto.Message) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	json := string(bytes)
	if err := jsonpb.UnmarshalString(json, message); err != nil {
		return errors.Wrap(err, "error parsing JSON file " + filename)
	}

	return nil
}

func setup(c *ZombieBattleground, pubKeyHex string, addr *loom.Address, ctx *contract.Context, t *testing.T) {
	debugEnabled = true

	// random key
	purchaseGatewayPrivateKeyHexString = "527969b4754fca7c3c6146c7c2a12ce1d0dda4a7e75cfb8e3465e0393d531176"

	updateInitRequest.InitData = &zb_data.InitData{}
	err := readJsonFileToProtobuf("../test_data/simple-init.json", updateInitRequest.InitData)
	assert.Nil(t, err)

	initRequest = zb_calls.InitRequest{
		DefaultDecks:         updateInitRequest.InitData.DefaultDecks,
		DefaultCollection:    updateInitRequest.InitData.DefaultCollection,
		Cards:                updateInitRequest.InitData.Cards,
		Overlords:            updateInitRequest.InitData.Overlords,
		AiDecks:              updateInitRequest.InitData.AiDecks,
		Version:              updateInitRequest.InitData.Version,
		Oracle:               updateInitRequest.InitData.Oracle,
		OverlordLeveling:     updateInitRequest.InitData.OverlordLeveling,
	}

	c = &ZombieBattleground{}
	pubKey, _ := hex.DecodeString(pubKeyHex)

	addr = &loom.Address{
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}

	*ctx = contract.WrapPluginContext(
		plugin.CreateFakeContext(*addr, *addr),
	)

	err = c.Init(*ctx, &initRequest)
	assert.Nil(t, err)

	request := zb_calls.UpdateContractConfigurationRequest{
		SetFiatPurchaseContractVersion: true,
		FiatPurchaseContractVersion:    3,
		SetInitialFiatPurchaseTxId:     true,
		InitialFiatPurchaseTxId:        &types.BigUInt{Value: common.BigUInt{Int: big.NewInt(100)}},
	}

	err = c.UpdateContractConfiguration(*ctx, &request)
	assert.Nil(t, err)
}

func setupAccount(c *ZombieBattleground, ctx contract.Context, upsertAccountRequest *zb_calls.UpsertAccountRequest, t *testing.T) {
	err := c.CreateAccount(ctx, upsertAccountRequest)
	assert.Nil(t, err)
}

func TestContractConfigurationAndState(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	t.Run("UpdateContractConfiguration should succeed", func(t *testing.T) {
		request := zb_calls.UpdateContractConfigurationRequest{
			SetFiatPurchaseContractVersion: true,
			FiatPurchaseContractVersion:    373,
			SetInitialFiatPurchaseTxId:     true,
			InitialFiatPurchaseTxId:        &types.BigUInt{Value: common.BigUInt{Int: big.NewInt(100)}},
		}

		err := c.UpdateContractConfiguration(ctx, &request)
		assert.Nil(t, err)
	})

	t.Run("GetContractState should succeed", func(t *testing.T) {
		response, err := c.GetContractState(ctx, &zb_calls.EmptyRequest{})
		assert.Nil(t, err)
		assert.Equal(t, int64(100), response.State.CurrentFiatPurchaseTxId.Value.Int.Int64())
		assert.Equal(t, uint64(0), response.State.LastPlasmachainBlockNumber)
	})

	t.Run("GetContractConfiguration should succeed", func(t *testing.T) {
		response, err := c.GetContractConfiguration(ctx, &zb_calls.EmptyRequest{})
		assert.Nil(t, err)
		assert.Equal(t, int64(100), response.Configuration.InitialFiatPurchaseTxId.Value.Int.Int64())
		assert.Equal(t, uint64(373), response.Configuration.FiatPurchaseContractVersion)
	})
}

func TestAccountOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "AccountUser",
		Image:   "PathToImage",
		Version: "v1",
	}, t)

	t.Run("UpdateAccount", func(t *testing.T) {
		account, err := c.UpdateAccount(ctx, &zb_calls.UpsertAccountRequest{
			UserId:      "AccountUser",
			Image:       "PathToImage2",
			CurrentTier: 5,
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(5), account.CurrentTier)
		assert.Equal(t, "PathToImage2", account.Image)
	})

	t.Run("GetAccount", func(t *testing.T) {
		account, err := c.GetAccount(ctx, &zb_calls.GetAccountRequest{
			UserId: "AccountUser",
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(5), account.CurrentTier)
		assert.Equal(t, "PathToImage2", account.Image)
	})
}

func TestCardCollectionCardOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "8996b813617b283f81ea1747fbddbe73fe4b5fce0eac0728e47de51d8e506701"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "CardUser",
		Image:   "PathToImage",
		Version: "v1",
	}, t)

	CardCollectionCard, err := c.GetCollection(ctx, &zb_calls.GetCollectionRequest{
		UserId: "CardUser",
		Version: "v1",
	})
	assert.Nil(t, err)
	assert.Equal(t, 12, len(CardCollectionCard.Cards))

}

func TestDeckOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "7796b813617b283f81ea1747fbddbe73fe4b5fce0eac0728e47de51d8e506701"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "DeckUser",
		Image:   "PathToImage",
		Version: "v1",
	}, t)

	t.Run("ListDecks", func(t *testing.T) {
		deckResponse, err := c.ListDecks(ctx, &zb_calls.ListDecksRequest{
			UserId: "DeckUser",
			Version: "v1",
		})
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(deckResponse.Decks))
		assert.Equal(t, int64(1), deckResponse.Decks[0].Id)
		assert.Equal(t, "Default", deckResponse.Decks[0].Name)
	})

	t.Run("GetDeck (Not Exists)", func(t *testing.T) {
		deckResponse, err := c.GetDeck(ctx, &zb_calls.GetDeckRequest{
			UserId: "DeckUser",
			DeckId: 0xDEADBEEF,
			Version: "v1",
		})
		assert.Equal(t, (*zb_calls.GetDeckResponse)(nil), deckResponse)
		assert.Equal(t, contract.ErrNotFound, err)
	})

	t.Run("GetDeck", func(t *testing.T) {
		deckResponse, err := c.GetDeck(ctx, &zb_calls.GetDeckRequest{
			UserId: "DeckUser",
			DeckId: 1,
			Version: "v1",
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), deckResponse.Deck.Id) // id should start from 1
		assert.Equal(t, "Default", deckResponse.Deck.Name)
	})

	var createDeckResponse *zb_calls.CreateDeckResponse
	t.Run("CreateDeck", func(t *testing.T) {
		var err error
		createDeckResponse, err = c.CreateDeck(ctx, &zb_calls.CreateDeckRequest{
			UserId: "DeckUser",
			Deck: &zb_data.Deck{
				Name:       "NewDeck",
				OverlordId: 1,
				Cards: []*zb_data.DeckCard{
					{
						Amount:   1,
						CardKey: battleground_proto.CardKey{MouldId: 43},
					},
					{
						Amount:   1,
						CardKey: battleground_proto.CardKey{MouldId: 48},
					},
				},
			},
			Version: "v1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, createDeckResponse)

		deckResponse, err := c.ListDecks(ctx, &zb_calls.ListDecksRequest{
			UserId: "DeckUser",
			Version: "v1",
		})

		assert.Equal(t, nil, err)
		assert.Equal(t, 2, len(deckResponse.Decks))
	})

	t.Run("CreateDeck (Invalid Requested Amount)", func(t *testing.T) {
		_, err := c.CreateDeck(ctx, &zb_calls.CreateDeckRequest{
			UserId: "DeckUser",
			Deck: &zb_data.Deck{
				Name:       "NewDeck",
				OverlordId: 1,
				Cards: []*zb_data.DeckCard{
					{
						Amount:   200,
						CardKey: battleground_proto.CardKey{MouldId: 43},
					},
					{
						Amount:   100,
						CardKey: battleground_proto.CardKey{MouldId: 48},
					},
				},
			},
			Version: "v1",
		})

		assert.NotNil(t, err)
	})

	t.Run("CreateDeck (Invalid Requested CardName)", func(t *testing.T) {
		_, err := c.CreateDeck(ctx, &zb_calls.CreateDeckRequest{
			UserId: "DeckUser",
			Deck: &zb_data.Deck{
				Name:       "NewDeck",
				OverlordId: 1,
				Cards: []*zb_data.DeckCard{
					{
						Amount:   2,
						CardKey: battleground_proto.CardKey{MouldId: -100},
					},
					{
						Amount:   1,
						CardKey: battleground_proto.CardKey{MouldId: -101},
					},
				},
			},
			Version: "v1",
		})

		assert.NotNil(t, err)
	})

	t.Run("CreateDeck (Same name not allowed)", func(t *testing.T) {
		_, err := c.CreateDeck(ctx, &zb_calls.CreateDeckRequest{
			UserId: "DeckUser",
			Deck: &zb_data.Deck{
				Name:       "Default",
				OverlordId: 1,
				Cards: []*zb_data.DeckCard{
					{
						Amount:   1,
						CardKey: battleground_proto.CardKey{MouldId: 43},
					},
					{
						Amount:   1,
						CardKey: battleground_proto.CardKey{MouldId: 48},
					},
				},
			},
			Version: "v1",
		})

		assert.NotNil(t, err)
	})

	t.Run("CreateDeck (Same name with different case not allowed)", func(t *testing.T) {
		_, err := c.CreateDeck(ctx, &zb_calls.CreateDeckRequest{
			UserId: "DeckUser",
			Deck: &zb_data.Deck{
				Name:       "nEWdECK",
				OverlordId: 1,
				Cards: []*zb_data.DeckCard{
					{
						Amount:   1,
						CardKey: battleground_proto.CardKey{MouldId: 43},
					},
					{
						Amount:   1,
						CardKey: battleground_proto.CardKey{MouldId: 48},
					},
				},
			},
			Version: "v1",
		})

		assert.NotNil(t, err)
	})

	t.Run("EditDeck", func(t *testing.T) {
		err := c.EditDeck(ctx, &zb_calls.EditDeckRequest{
			UserId: "DeckUser",
			Deck: &zb_data.Deck{
				Id:         2,
				Name:       "Edited",
				OverlordId: 1,
				Cards: []*zb_data.DeckCard{
					{
						Amount:   1,
						CardKey: battleground_proto.CardKey{MouldId: 43},
					},
					{
						Amount:   1,
						CardKey: battleground_proto.CardKey{MouldId: 48},
					},
				},
			},
			Version: "v1",
		})
		assert.Nil(t, err)

		getDeckResponse, err := c.GetDeck(ctx, &zb_calls.GetDeckRequest{
			UserId: "DeckUser",
			DeckId: 2,
			Version: "v1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, getDeckResponse)
		assert.Equal(t, "Edited", getDeckResponse.Deck.Name)
	})

	t.Run("EditDeck (attempt to set more number of cards)", func(t *testing.T) {
		t.Skip("Edit deck skips checking the number of cards")
		err := c.EditDeck(ctx, &zb_calls.EditDeckRequest{
			UserId: "DeckUser",
			Deck: &zb_data.Deck{
				Id:         2,
				Name:       "Edited",
				OverlordId: 1,
				Cards: []*zb_data.DeckCard{
					{
						Amount:  100,
						CardKey: battleground_proto.CardKey{MouldId: 43},
					},
					{
						Amount:  1,
						CardKey: battleground_proto.CardKey{MouldId: 48},
					},
				},
			},
			Version: "v1",
		})
		assert.NotNil(t, err)
	})

	t.Run("EditDeck (same name while editing is allowed)", func(t *testing.T) {
		err := c.EditDeck(ctx, &zb_calls.EditDeckRequest{
			UserId: "DeckUser",
			Deck: &zb_data.Deck{
				Id:         2,
				Name:       "Edited",
				OverlordId: 1,
				Cards: []*zb_data.DeckCard{
					{
						Amount:   1,
						CardKey: battleground_proto.CardKey{MouldId: 43},
					},
					{
						Amount:   1,
						CardKey: battleground_proto.CardKey{MouldId: 48},
					},
				},
			},
			Version: "v1",
		})

		assert.Nil(t, err)
	})

	t.Run("EditDeck (attempt to set duplicate name with different case)", func(t *testing.T) {
		err := c.EditDeck(ctx, &zb_calls.EditDeckRequest{
			UserId: "DeckUser",
			Deck: &zb_data.Deck{
				Id:         2,
				Name:       "dEFAULT",
				OverlordId: 1,
				Cards: []*zb_data.DeckCard{
					{
						Amount:   1,
						CardKey: battleground_proto.CardKey{MouldId: 43},
					},
					{
						Amount:   1,
						CardKey: battleground_proto.CardKey{MouldId: 48},
					},
				},
			},
			Version: "v1",
		})

		assert.NotNil(t, err)
	})

	t.Run("DeleteDeck", func(t *testing.T) {
		assert.NotNil(t, createDeckResponse)
		err := c.DeleteDeck(ctx, &zb_calls.DeleteDeckRequest{
			UserId: "DeckUser",
			DeckId: createDeckResponse.DeckId,
			Version: "v1",
		})

		assert.Nil(t, err)
	})

	t.Run("DeleteDeck (Non existant)", func(t *testing.T) {
		err := c.DeleteDeck(ctx, &zb_calls.DeleteDeckRequest{
			UserId: "DeckUser",
			DeckId: 0xDEADBEEF,
			Version: "v1",
		})

		assert.NotNil(t, err)
	})
}

func TestCardOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	t.Run("ListCardLibrary", func(t *testing.T) {
		cardResponse, err := c.ListCardLibrary(ctx, &zb_calls.ListCardLibraryRequest{
			Version: "v1",
		})

		assert.Nil(t, err)
		// we expect all the cards in InitRequest.Cards
		assert.Equal(t, len(initRequest.Cards), len(cardResponse.Cards))
	})

	t.Run("ListOverlordLibrary", func(t *testing.T) {
		overlordsResponse, err := c.ListOverlordLibrary(ctx, &zb_calls.ListOverlordLibraryRequest{
			Version: "v1",
		})

		assert.Nil(t, err)
		assert.Equal(t, 2, len(overlordsResponse.Overlords))
	})
}

func TestCardDataUpgradeAndValidation(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "7796b813617b283f81ea1747fbddbe73fe4b5fce0eac0728e47de51d8e506701"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "DeckUser",
		Image:   "PathToImage",
		Version: "v1",
	}, t)

	t.Run("Should remove unknown cards from decks", func(t *testing.T) {
		deckList := &zb_data.DeckList{
			Decks: []*zb_data.Deck{
				{
					Name:       "NewDeck",
					OverlordId: 1,
					Cards: []*zb_data.DeckCard{
						{
							Amount:  1,
							CardKey: battleground_proto.CardKey{MouldId: -1},
						},
						{
							Amount:  2,
							CardKey: battleground_proto.CardKey{MouldId: 1},
						},
					},
				},
			},
		}

		err := saveDecks(ctx, "v1", "DeckUser", deckList)
		assert.Nil(t, err)

		deckResponse, err := c.ListDecks(ctx, &zb_calls.ListDecksRequest{
			UserId: "DeckUser",
			Version: "v1",
		})

		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(deckResponse.Decks))
		assert.Equal(t, 1, len(deckResponse.Decks[0].Cards))
		assert.Equal(t, int64(1), deckResponse.Decks[0].Cards[0].CardKey.MouldId)
	})
}

func TestOverlordsOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "7696b824516b283f81ea1747fbddbe73fe4b5fce0eac0728e47de41d8e306701"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "OverlordUser",
		Image:   "PathToImage",
		Version: "v1",
	}, t)

	t.Run("ListOverlordUserInstances", func(t *testing.T) {
		overlordsResponse, err := c.ListOverlordUserInstances(ctx, &zb_calls.ListOverlordUserInstancesRequest{
			UserId:  "OverlordUser",
			Version: "v1",
		})

		assert.Nil(t, err)
		assert.Equal(t, 2, len(overlordsResponse.Overlords))
	})

	t.Run("GetOverlordUserInstance", func(t *testing.T) {
		overlordResponse, err := c.GetOverlordUserInstance(ctx, &zb_calls.GetOverlordUserInstanceRequest{
			UserId:     "OverlordUser",
			OverlordId: 1,
			Version:    "v1",
		})

		assert.Nil(t, err)
		assert.NotNil(t, overlordResponse.Overlord)
	})

	t.Run("GetOverlordUserInstance (Overlord not exists)", func(t *testing.T) {
		_, err := c.GetOverlordUserInstance(ctx, &zb_calls.GetOverlordUserInstanceRequest{
			UserId:     "OverlordUser",
			OverlordId: 10,
			Version:    "v1",
		})

		assert.NotNil(t, err)
	})
}

func TestUpdateInitDataOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	t.Run("UpdateInit", func(t *testing.T) {
		err := c.UpdateInit(ctx, &updateInitRequest)

		assert.Nil(t, err)
	})

	t.Run("UpdateInit with missing data", func(t *testing.T) {
		updateInitRequestWithMissingCards := proto.Clone(&updateInitRequest).(*zb_calls.UpdateInitRequest)
		updateInitRequestWithMissingCards.InitData.Cards = nil
		err := c.UpdateInit(ctx, updateInitRequestWithMissingCards)

		assert.NotNil(t, err)
		assert.Equal(t, "'cards' key missing", err.Error())
	})
}

func TestFindMatchOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)

	// make users have decks
	t.Run("ListDecksPlayer1", func(t *testing.T) {
		deckResponse, err := c.ListDecks(ctx, &zb_calls.ListDecksRequest{
			UserId: "player-1",
			Version: "v1",
		})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(deckResponse.Decks))
		assert.Equal(t, int64(1), deckResponse.Decks[0].Id)
		assert.Equal(t, "Default", deckResponse.Decks[0].Name)
	})
	t.Run("ListDecksPlayer2", func(t *testing.T) {
		deckResponse, err := c.ListDecks(ctx, &zb_calls.ListDecksRequest{
			UserId: "player-2",
			Version: "v1",
		})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(deckResponse.Decks))
		assert.Equal(t, int64(1), deckResponse.Decks[0].Id)
		assert.Equal(t, "Default", deckResponse.Decks[0].Name)
	})

	var matchID int64

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "two players should be matching")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "two players should be matching")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "two players should be matching")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-2",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "two players should be matching")
		assert.Equal(t, zb_data.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("GetMatch", func(t *testing.T) {
		response, err := c.GetMatch(ctx, &zb_calls.GetMatchRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the second player should 2 player states")
		assert.Equal(t, zb_data.Match_Started, response.Match.Status, "match status should be 'started'")
	})

	t.Run("EndMatch", func(t *testing.T) {
		_, err := c.EndMatch(ctx, &zb_calls.EndMatchRequest{
			MatchId:          matchID,
			UserId:           "player-1",
			WinnerId:         "player-2",
			MatchExperiences: []int64{123, 350},
		})
		assert.Nil(t, err)
	})

	t.Run("Check level and experience after match", func(t *testing.T) {
		getOverlordResponse1, err := c.GetOverlordUserInstance(ctx, &zb_calls.GetOverlordUserInstanceRequest{
			UserId:     "player-1",
			OverlordId: 1,
			Version:    "v1",
		})

		assert.Nil(t, err)
		assert.Equal(t, int64(123), getOverlordResponse1.Overlord.UserData.Experience)
		assert.True(t, len(getOverlordResponse1.Overlord.UserData.UnlockedSkillIds) == 0)
		assert.Equal(t, int64(1), getOverlordResponse1.Overlord.UserData.Level)

		getOverlordResponse2, err := c.GetOverlordUserInstance(ctx, &zb_calls.GetOverlordUserInstanceRequest{
			UserId:     "player-2",
			OverlordId: 1,
			Version:    "v1",
		})

		assert.Nil(t, err)
		assert.Equal(t, int64(350), getOverlordResponse2.Overlord.UserData.Experience)
		assert.Equal(t, true, len(getOverlordResponse2.Overlord.UserData.UnlockedSkillIds) > 0)
		assert.Equal(t, int64(7), getOverlordResponse2.Overlord.UserData.Level)
	})

	t.Run("Check level/experience notifications after match", func(t *testing.T) {
		getNotificationsResponse1, err := c.GetNotifications(ctx, &zb_calls.GetNotificationsRequest{
			UserId: "player-1",
		})

		assert.Nil(t, err)
		notificationEndMatch1 := getNotificationsResponse1.Notifications[0].Notification.(*zb_data.Notification_EndMatch).EndMatch
		assert.Equal(t, int64(1), notificationEndMatch1.OverlordId)
		assert.Equal(t, int32(1), notificationEndMatch1.OldLevel)
		assert.Equal(t, int64(0), notificationEndMatch1.OldExperience)
		assert.Equal(t, int32(1), notificationEndMatch1.NewLevel)
		assert.Equal(t, int64(123), notificationEndMatch1.NewExperience)
		assert.Equal(t, false, notificationEndMatch1.IsWin)

		getNotificationsResponse2, err := c.GetNotifications(ctx, &zb_calls.GetNotificationsRequest{
			UserId: "player-2",
		})

		assert.Nil(t, err)
		assert.Equal(t, int(1), len(getNotificationsResponse2.Notifications))
		notificationEndMatch2 := getNotificationsResponse2.Notifications[0].Notification.(*zb_data.Notification_EndMatch).EndMatch
		assert.Equal(t, int64(1), notificationEndMatch2.OverlordId)
		assert.Equal(t, int32(1), notificationEndMatch2.OldLevel)
		assert.Equal(t, int64(0), notificationEndMatch2.OldExperience)
		assert.Equal(t, int32(7), notificationEndMatch2.NewLevel)
		assert.Equal(t, int64(350), notificationEndMatch2.NewExperience)
		assert.Equal(t, true, notificationEndMatch2.IsWin)
	})

	t.Run("GetMatchAfterLeaving", func(t *testing.T) {
		response, err := c.GetMatch(ctx, &zb_calls.GetMatchRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the second player should 2 player states")
		assert.Equal(t, zb_data.Match_Ended, response.Match.Status, "match status should be 'ended'")
	})
}

func TestOverlordLeveling(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	userId := defaultUserIdPrefix + "373"
	version := "v1"
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  userId,
		Version: version,
	}, t)

	levelingData, err := loadOverlordLevelingData(ctx, version)
	assert.Nil(t, err)

	t.Run("Level 2 is reached and one reward is given", func(t *testing.T) {
		addedExperience := getRequiredExperienceForLevel(levelingData, 2)
		err = applyExperience(ctx, version, levelingData, userId, big.NewInt(373), 1, int64(addedExperience), 1, true)
		assert.Nil(t, err)

		getNotificationsResponse1, err := c.GetNotifications(ctx, &zb_calls.GetNotificationsRequest{UserId: userId})
		assert.Nil(t, err)

		notificationEndMatch := getNotificationsResponse1.Notifications[0].Notification.(*zb_data.Notification_EndMatch).EndMatch
		assert.Equal(t, int64(1), notificationEndMatch.OverlordId)
		assert.Equal(t, int32(1), notificationEndMatch.OldLevel)
		assert.Equal(t, int64(0), notificationEndMatch.OldExperience)
		assert.Equal(t, int64(addedExperience), notificationEndMatch.NewExperience)
		assert.Equal(t, int32(2), notificationEndMatch.NewLevel)

		assert.Equal(t, 1, len(notificationEndMatch.Rewards))
		assert.Equal(t, int32(2), notificationEndMatch.Rewards[0].Level)
		boosterPackReward, ok := notificationEndMatch.Rewards[0].GetReward().(*zb_data.LevelReward_BoosterPackReward)
		assert.True(t, ok)
		assert.Equal(t, int32(1), boosterPackReward.BoosterPackReward.Amount)
	})

	t.Run("Level 3 is reached and no rewards are given", func(t *testing.T) {
		addedExperience := levelingData.ExperienceStep
		err = applyExperience(ctx, version, levelingData, userId, big.NewInt(373), 1, int64(addedExperience), 1, true)
		assert.Nil(t, err)

		getNotificationsResponse1, err := c.GetNotifications(ctx, &zb_calls.GetNotificationsRequest{UserId: userId})
		assert.Nil(t, err)

		notificationEndMatch := getNotificationsResponse1.Notifications[0].Notification.(*zb_data.Notification_EndMatch).EndMatch
		assert.Equal(t, int64(1), notificationEndMatch.OverlordId)
		assert.Equal(t, int32(2), notificationEndMatch.OldLevel)
		assert.Equal(t, int64(levelingData.Fixed + levelingData.ExperienceStep), notificationEndMatch.OldExperience)
		assert.Equal(t, int32(3), notificationEndMatch.NewLevel)
		assert.Equal(t, getRequiredExperienceForLevel(levelingData, 2) + int64(addedExperience), notificationEndMatch.NewExperience)

		assert.Equal(t, 0, len(notificationEndMatch.Rewards))
	})
}

func TestCancelFindMatchOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)

	var matchID int64

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("CancelFindmatch", func(t *testing.T) {
		_, err := c.CancelFindMatch(ctx, &zb_calls.CancelFindMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
	})

	t.Run("GetMatch", func(t *testing.T) {
		response, err := c.GetMatch(ctx, &zb_calls.GetMatchRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.Equal(t, zb_data.Match_Canceled, response.Match.Status)
	})
}

func TestCancelFindMatchOnEndedMatchOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)

	var matchID int64

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-2",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("EndMatch", func(t *testing.T) {
		response, err := c.EndMatch(ctx, &zb_calls.EndMatchRequest{
			UserId:   "player-2",
			MatchId:  matchID,
			WinnerId: "player-2",
			MatchExperiences: []int64{0, 0},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})

	t.Run("GetMatch", func(t *testing.T) {
		response, err := c.GetMatch(ctx, &zb_calls.GetMatchRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.Equal(t, zb_data.Match_Ended, response.Match.Status)
	})

	t.Run("CancelFindmatch", func(t *testing.T) {
		_, err := c.CancelFindMatch(ctx, &zb_calls.CancelFindMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
	})

	t.Run("GetMatch", func(t *testing.T) {
		response, err := c.GetMatch(ctx, &zb_calls.GetMatchRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.Equal(t, zb_data.Match_Ended, response.Match.Status)
	})
}

func TestFindMatchWithTagOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-3",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-1-tag",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-2-tag",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-3-tag",
		Version: "v1",
	}, t)

	var matchID, matchIDTag int64

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-2",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	tags := []string{"tag1"}

	t.Run("RegisterPlayerPoolTag", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1-tag",
				Version: "v1",
				Tags:    tags,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPoolTag", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2-tag",
				Version: "v1",
				Tags:    tags,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("FindmatchTag", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-1-tag",
			Tags:   tags,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchIDTag = response.Match.Id
	})

	t.Run("FindmatchTag", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-2-tag",
			Tags:   tags,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchIDTag, response.Match.Id)
	})

	t.Run("AcceptMatchTag", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-1-tag",
			MatchId: matchIDTag,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchIDTag, response.Match.Id)
		assert.NotEqual(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatchTag", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-2-tag",
			MatchId: matchIDTag,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchIDTag, response.Match.Id)
		assert.NotEqual(t, matchID, response.Match.Id)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the second player should 2 player states")
		assert.Equal(t, zb_data.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	// check tag and non-tag players don't get matched
	tags = []string{"tag3"}

	t.Run("RegisterPlayerPoolTag", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-3",
				Version: "v1",
				Tags:    tags,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPoolTag", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-3-tag",
				Version: "v1",
				Tags:    tags,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		_, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-3",
		})
		assert.NotNil(t, err)
	})

	t.Run("FindmatchTag", func(t *testing.T) {
		_, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-3-tag",
		})
		assert.NotNil(t, err)
	})
}

func TestFindMatchWithTagGroupOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-1-tag",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-2-tag",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-3-tag",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-4-tag",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-5-tag",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-6-tag",
		Version: "v1",
	}, t)

	tags1 := []string{"tags1"}
	tags2 := []string{"tags2"}
	tags3 := []string{"tags3", "othertag"}

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1-tag",
				Version: "v1",
				Tags:    tags1,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-3-tag",
				Version: "v1",
				Tags:    tags2,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-5-tag",
				Version: "v1",
				Tags:    tags3,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-1-tag",
			Tags:   tags1,
		})
		assert.Nil(t, err)
		assert.Equal(t, false, response.MatchFound)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-3-tag",
			Tags:   tags2,
		})
		assert.Nil(t, err)
		assert.Equal(t, false, response.MatchFound)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-5-tag",
			Tags:   tags3,
		})
		assert.Nil(t, err)
		assert.Equal(t, false, response.MatchFound)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-6-tag",
				Version: "v1",
				Tags:    tags3,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2-tag",
				Version: "v1",
				Tags:    tags1,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-4-tag",
				Version: "v1",
				Tags:    tags2,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-1-tag",
			Tags:   tags1,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, "player-1-tag", response.Match.PlayerStates[0].Id, "Player should be player-1-tag")
		assert.Equal(t, "player-2-tag", response.Match.PlayerStates[1].Id, "Player should be player-2-tag")
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-3-tag",
			Tags:   tags2,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, "player-3-tag", response.Match.PlayerStates[0].Id, "Player should be player-3-tag")
		assert.Equal(t, "player-4-tag", response.Match.PlayerStates[1].Id, "Player should be player-4-tag")
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-5-tag",
			Tags:   tags3,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, "player-5-tag", response.Match.PlayerStates[0].Id, "Player should be player-5-tag")
		assert.Equal(t, "player-6-tag", response.Match.PlayerStates[1].Id, "Player should be player-6-tag")
	})
}

func TestMatchMakingPlayerPool(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	numPlayers := 10
	for i := 0; i < numPlayers; i++ {
		setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
			UserId:  fmt.Sprintf("player-%d", i+1),
			Version: "v1",
		}, t)
	}

	for i := 0; i < numPlayers; i++ {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  fmt.Sprintf("player-%d", i+1),
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	}

	for i := 0; i < numPlayers; i++ {
		func(i int) {
			response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
				UserId: fmt.Sprintf("player-%d", i+1),
			})
			assert.Nil(t, err)
			assert.NotNil(t, response)
		}(i)
	}
}

func TestMatchMakingTimeout(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr *loom.Address
	var ctx contract.Context

	// setup ctx
	c = &ZombieBattleground{}
	pubKey, _ := hex.DecodeString(pubKeyHexString)
	addr = &loom.Address{
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}
	now := time.Now()
	fc := plugin.CreateFakeContext(*addr, *addr)
	fc.SetTime(now)
	ctx = contract.WrapPluginContext(fc)
	err := c.Init(ctx, &initRequest)
	assert.Nil(t, err)

	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-3",
		Version: "v1",
	}, t)

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-3",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	var matchID int64
	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	// move time forward to expire the matchmaking
	fc.SetTime(now.Add(2 * MMTimeout))

	t.Run("FindMatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.Equal(t, zb_data.Match_Timedout, response.Match.Status)
	})

	t.Run("GetMatch", func(t *testing.T) {
		response, err := c.GetMatch(ctx, &zb_calls.GetMatchRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.Equal(t, zb_data.Match_Timedout, response.Match.Status)
	})
}

func TestGameStateOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)

	// make users have decks
	t.Run("ListDecksPlayer1", func(t *testing.T) {
		deckResponse, err := c.ListDecks(ctx, &zb_calls.ListDecksRequest{
			UserId: "player-1",
			Version: "v1",
		})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(deckResponse.Decks))
		assert.Equal(t, int64(1), deckResponse.Decks[0].Id)
		assert.Equal(t, "Default", deckResponse.Decks[0].Name)
	})
	t.Run("ListDecksPlayer2", func(t *testing.T) {
		deckResponse, err := c.ListDecks(ctx, &zb_calls.ListDecksRequest{
			UserId: "player-2",
			Version: "v1",
		})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(deckResponse.Decks))
		assert.Equal(t, int64(1), deckResponse.Decks[0].Id)
		assert.Equal(t, "Default", deckResponse.Decks[0].Name)
	})

	var matchID int64

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-2",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("GetMatch", func(t *testing.T) {
		response, err := c.GetMatch(ctx, &zb_calls.GetMatchRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the second player should 2 player states")
		assert.Equal(t, zb_data.Match_Started, response.Match.Status, "match status should be 'started'")
	})

	// Note: since the toss coin seed is always 0 for testing, we always get 0 as the first player
	t.Run("SendEndturnPlayer2_Failed", func(t *testing.T) {
		_, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_EndTurn,
				PlayerId:   "player-2",
			},
		})
		assert.NotNil(t, err)
		assert.Equal(t, errInvalidPlayer, err)
	})

	t.Run("SendEndturnPlayer1_Success", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_EndTurn,
				PlayerId:   "player-1",
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)

		gameState, err := loadGameState(ctx, matchID)
		assert.EqualValues(t, 1, gameState.CurrentPlayerIndex, "player-2 should be active")
	})
	t.Run("SendEndturnPlayer2_Success", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_EndTurn,
				PlayerId:   "player-2",
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)

		gameState, err := loadGameState(ctx, matchID)
		assert.EqualValues(t, 0, gameState.CurrentPlayerIndex, "player-1 should be active")
	})
	t.Run("SendCardPlayPlayer1", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_CardPlay,
				PlayerId:   "player-1",
				Action: &zb_data.PlayerAction_CardPlay{
					CardPlay: &zb_data.PlayerActionCardPlay{
						Card: &zb_data.InstanceId{Id: 8},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("SendCardAbilityPlayer1", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_CardAbilityUsed,
				PlayerId:   "player-1",
				Action: &zb_data.PlayerAction_CardAbilityUsed{
					CardAbilityUsed: &zb_data.PlayerActionCardAbilityUsed{
						Card: &zb_data.InstanceId{Id: 1},
						Targets: []*zb_data.Unit{
							&zb_data.Unit{
								InstanceId: &zb_data.InstanceId{Id: 2},
							},
						},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("SendOverlordSkillUsedPlayer1", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_OverlordSkillUsed,
				PlayerId:   "player-1",
				Action: &zb_data.PlayerAction_OverlordSkillUsed{
					OverlordSkillUsed: &zb_data.PlayerActionOverlordSkillUsed{
						SkillId: 1,
						Target: &zb_data.Unit{
							InstanceId: &zb_data.InstanceId{Id: 2},
						},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("SendRankBuff", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_RankBuff,
				PlayerId:   "player-1",
				Action: &zb_data.PlayerAction_RankBuff{
					RankBuff: &zb_data.PlayerActionRankBuff{
						Card: &zb_data.InstanceId{Id: 1},
						Targets: []*zb_data.Unit{
							&zb_data.Unit{
								InstanceId: &zb_data.InstanceId{Id: 2},
							},
						},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("SendEndturnPlayer1_Success2", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_EndTurn,
				PlayerId:   "player-1",
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)

		gameState, err := loadGameState(ctx, matchID)
		assert.EqualValues(t, 1, gameState.CurrentPlayerIndex, "player-2 should be active")
	})
	t.Run("SendCardPlayPlayer2", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_CardPlay,
				PlayerId:   "player-2",
				Action: &zb_data.PlayerAction_CardPlay{
					CardPlay: &zb_data.PlayerActionCardPlay{
						Card: &zb_data.InstanceId{Id: 13},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("SendCardAttackPlayer2", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_CardAttack,
				PlayerId:   "player-2",
				Action: &zb_data.PlayerAction_CardAttack{
					CardAttack: &zb_data.PlayerActionCardAttack{
						Attacker: &zb_data.InstanceId{Id: 13},
						Target: &zb_data.Unit{
							InstanceId: &zb_data.InstanceId{Id: 8},
						},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("SendCardAbilityPlayer2", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_CardAbilityUsed,
				PlayerId:   "player-2",
				Action: &zb_data.PlayerAction_CardAbilityUsed{
					CardAbilityUsed: &zb_data.PlayerActionCardAbilityUsed{
						Card: &zb_data.InstanceId{Id: 1},
						Targets: []*zb_data.Unit{
							&zb_data.Unit{
								InstanceId: &zb_data.InstanceId{Id: 2},
							},
						},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("SendOverlordSkillUsedPlayer2", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_OverlordSkillUsed,
				PlayerId:   "player-2",
				Action: &zb_data.PlayerAction_OverlordSkillUsed{
					OverlordSkillUsed: &zb_data.PlayerActionOverlordSkillUsed{
						SkillId: 1,
						Target: &zb_data.Unit{
							InstanceId: &zb_data.InstanceId{Id: 2},
						},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("LeaveMatch", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_LeaveMatch,
				PlayerId:   "player-2",
				Action: &zb_data.PlayerAction_LeaveMatch{
					LeaveMatch: &zb_data.PlayerActionLeaveMatch{},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb_calls.GetGameStateRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.GameState.IsEnded, "game state should be ended after use leaves the match")
	})
}

func TestGameModeOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)

	var ID string

	t.Run("Add Game Mode", func(t *testing.T) {
		gameMode, err := c.AddGameMode(ctx, &zb_calls.GameModeRequest{
			Name:        "Test game mode",
			Description: "Just a test",
			Version:     "0.1",
			Address:     "0xf16a25a1b4e6434bacf9d037d69d675dcf852691",
		})
		assert.Nil(t, err)
		ID = gameMode.ID
		assert.Equal(t, "Test game mode", gameMode.Name)
		assert.Equal(t, "Just a test", gameMode.Description)
		assert.Equal(t, "0.1", gameMode.Version)
		assert.Equal(t, zb_data.GameModeType_Community, gameMode.GameModeType)
	})

	t.Run("Get Game Mode", func(t *testing.T) {
		gameMode, err := c.GetGameMode(ctx, &zb_calls.GetGameModeRequest{
			ID: ID,
		})
		assert.Nil(t, err)
		assert.Equal(t, "Test game mode", gameMode.Name)
		assert.Equal(t, "Just a test", gameMode.Description)
		assert.Equal(t, "0.1", gameMode.Version)
		assert.Equal(t, zb_data.GameModeType_Community, gameMode.GameModeType)
	})

	t.Run("Add another Game Mode", func(t *testing.T) {
		gameMode, err := c.AddGameMode(ctx, &zb_calls.GameModeRequest{
			Name:        "Test game mode 2",
			Description: "Just another test",
			Version:     "0.1",
			Address:     "0xf16a25a1b4e6434bacf9d037d69d675dcf852692",
		})
		assert.Nil(t, err)
		assert.Equal(t, "Test game mode 2", gameMode.Name)
		assert.Equal(t, "Just another test", gameMode.Description)
		assert.Equal(t, "0.1", gameMode.Version)
		assert.Equal(t, zb_data.GameModeType_Community, gameMode.GameModeType)
	})

	t.Run("Update a Game Mode", func(t *testing.T) {
		gameMode, err := c.UpdateGameMode(ctx, &zb_calls.UpdateGameModeRequest{
			ID:          ID,
			Name:        "Test game mode",
			Description: "Changed description",
			Version:     "0.2",
		})
		assert.Nil(t, err)
		assert.Equal(t, "Test game mode", gameMode.Name)
		assert.Equal(t, "Changed description", gameMode.Description)
		assert.Equal(t, "0.2", gameMode.Version)
		assert.Equal(t, zb_data.GameModeType_Community, gameMode.GameModeType)
	})

	t.Run("List Game Modes", func(t *testing.T) {
		gameModeList, err := c.ListGameModes(ctx, &zb_calls.ListGameModesRequest{})
		assert.Nil(t, err)
		assert.Equal(t, 2, len(gameModeList.GameModes))
		assert.Equal(t, ID, gameModeList.GameModes[0].ID)
		assert.Equal(t, "0.2", gameModeList.GameModes[0].Version)
		assert.Equal(t, "Test game mode 2", gameModeList.GameModes[1].Name)
	})

	t.Run("Delete Game Mode", func(t *testing.T) {
		err := c.DeleteGameMode(ctx, &zb_calls.DeleteGameModeRequest{
			ID: ID,
		})
		assert.Nil(t, err)
	})

	t.Run("GameModeList should not contain deleted GameMode", func(t *testing.T) {
		gameModeList, err := c.ListGameModes(ctx, &zb_calls.ListGameModesRequest{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(gameModeList.GameModes))
		assert.NotEqual(t, ID, gameModeList.GameModes[0].ID)
		assert.Equal(t, "Test game mode 2", gameModeList.GameModes[0].Name)
	})
}

func TestCardPlayOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)

	var matchID int64

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("Acceptmatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the second player should 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'Matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("Acceptmatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-2",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the second player should 2 player states")
		assert.Equal(t, zb_data.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("SendCardPlayPlayer1", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_CardPlay,
				PlayerId:   "player-1",
				Action: &zb_data.PlayerAction_CardPlay{
					CardPlay: &zb_data.PlayerActionCardPlay{
						Card: &zb_data.InstanceId{Id: 8},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
}

func TestCheckGameStatusWithTimeout(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr *loom.Address
	var ctx contract.Context

	// setup ctx
	c = &ZombieBattleground{}
	pubKey, _ := hex.DecodeString(pubKeyHexString)
	addr = &loom.Address{
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}
	now := time.Now()
	fc := plugin.CreateFakeContext(*addr, *addr)
	fc.SetTime(now)
	ctx = contract.WrapPluginContext(fc)
	err := c.Init(ctx, &initRequest)
	assert.Nil(t, err)

	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)

	var matchID int64

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
				DebugCheats: zb_data.DebugCheatsConfiguration{
					Enabled:             true,
					UseCustomRandomSeed: true,
					CustomRandomSeed:    2,
				},
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
				DebugCheats: zb_data.DebugCheatsConfiguration{
					Enabled:             true,
					UseCustomRandomSeed: true,
					CustomRandomSeed:    2,
				},
			},
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-2",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb_calls.GetGameStateRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
		assert.EqualValues(t, 0, response.GameState.CurrentPlayerIndex)
	})

	t.Run("SendEndturnPlayer1_Success", func(t *testing.T) {
		_, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: 1,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_EndTurn,
				PlayerId:   "player-1",
			},
		})
		assert.Nil(t, err)
	})
	t.Run("SendEndturnPlayer2_Success", func(t *testing.T) {
		_, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: 1,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_EndTurn,
				PlayerId:   "player-2",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb_calls.GetGameStateRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
		assert.EqualValues(t, 0, response.GameState.CurrentPlayerIndex)
	})

	t.Run("Player1_Leave", func(t *testing.T) {
		_, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: 1,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_LeaveMatch,
				PlayerId:   "player-1",
				Action: &zb_data.PlayerAction_LeaveMatch{
					LeaveMatch: &zb_data.PlayerActionLeaveMatch{
						Reason: zb_data.PlayerActionLeaveMatch_PlayerLeave,
					},
				},
			},
		})
		assert.Nil(t, err)
	})

	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb_calls.GetGameStateRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.GameState.IsEnded, "game state should be ended after use leaves the match")
		latestAction := response.GameState.PlayerActions[len(response.GameState.PlayerActions)-1]
		assert.Equal(t, zb_enums.PlayerActionType_LeaveMatch, latestAction.ActionType)
		assert.Equal(t, "player-2", response.GameState.Winner)
	})

	t.Run("SendAnyActionShould_Failed", func(t *testing.T) {
		_, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: 1,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_EndTurn,
				PlayerId:   "player-2",
			},
		})
		assert.NotNil(t, err)
	})
}

func TestCheckGameStatusNoPlayerAction(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr *loom.Address
	var ctx contract.Context

	// setup ctx
	c = &ZombieBattleground{}
	pubKey, _ := hex.DecodeString(pubKeyHexString)
	addr = &loom.Address{
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}
	now := time.Now()
	fc := plugin.CreateFakeContext(*addr, *addr)
	fc.SetTime(now)
	ctx = contract.WrapPluginContext(fc)
	err := c.Init(ctx, &initRequest)
	assert.Nil(t, err)

	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)

	var matchID int64

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
				DebugCheats: zb_data.DebugCheatsConfiguration{
					Enabled:             true,
					UseCustomRandomSeed: true,
					CustomRandomSeed:    2,
				},
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
				DebugCheats: zb_data.DebugCheatsConfiguration{
					Enabled:             true,
					UseCustomRandomSeed: true,
					CustomRandomSeed:    2,
				},
			},
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-2",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb_calls.GetGameStateRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
		assert.EqualValues(t, 0, response.GameState.CurrentPlayerIndex)
	})

	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb_calls.GetGameStateRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
		assert.EqualValues(t, 0, response.GameState.CurrentPlayerIndex)
	})

	// player1 don't sent enturn within TurnTimeout
	// move time forward to expire the player's turn
	fc.SetTime(now.Add(TurnTimeout + (time.Second * 10)))

	t.Run("Player2_CheckStatus", func(t *testing.T) {
		_, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: 1,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_LeaveMatch,
				PlayerId:   "player-1",
				Action: &zb_data.PlayerAction_LeaveMatch{
					LeaveMatch: &zb_data.PlayerActionLeaveMatch{
						Reason: zb_data.PlayerActionLeaveMatch_PlayerLeave,
					},
				},
			},
		})
		assert.Nil(t, err)
	})

	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb_calls.GetGameStateRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.GameState.IsEnded, "game state should be ended after use leaves the match")
		latestAction := response.GameState.PlayerActions[len(response.GameState.PlayerActions)-1]
		assert.Equal(t, zb_enums.PlayerActionType_LeaveMatch, latestAction.ActionType)
		assert.Equal(t, "player-2", response.GameState.Winner)
	})
}

func TestAIDeckOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	aiDecks := []*zb_data.AIDeck{
		{
			Deck: &zb_data.Deck{
				Id:         1,
				OverlordId: 1,
				Name:       "AI Decks",
				Cards: []*zb_data.DeckCard{
					{CardKey: battleground_proto.CardKey{MouldId: 1}, Amount: 2},
					{CardKey: battleground_proto.CardKey{MouldId: 2}, Amount: 2},
					{CardKey: battleground_proto.CardKey{MouldId: 3}, Amount: 2},
					{CardKey: battleground_proto.CardKey{MouldId: 4}, Amount: 2},
					{CardKey: battleground_proto.CardKey{MouldId: 5}, Amount: 2},
					{CardKey: battleground_proto.CardKey{MouldId: 6}, Amount: 2},
					{CardKey: battleground_proto.CardKey{MouldId: 7}, Amount: 1},
					{CardKey: battleground_proto.CardKey{MouldId: 8}, Amount: 1},
					{CardKey: battleground_proto.CardKey{MouldId: 9}, Amount: 1},
					{CardKey: battleground_proto.CardKey{MouldId: 10}, Amount: 1},
					{CardKey: battleground_proto.CardKey{MouldId: 11}, Amount: 1},
					{CardKey: battleground_proto.CardKey{MouldId: 12}, Amount: 1},
				},
			},
			Type: zb_enums.AIType_MIXED_AI,
		},
	}

	t.Run("Get AI Decks", func(t *testing.T) {
		err := saveAIDecks(ctx, "v1", &zb_data.AIDeckList{Decks: aiDecks})
		_, err = c.GetAIDecks(ctx, &zb_calls.GetAIDecksRequest{
			Version: "v1",
		})
		assert.Nil(t, err)
	})

	aiDecks = []*zb_data.AIDeck{
		{
			Deck: &zb_data.Deck{
				Id:         1,
				OverlordId: 1,
				Name:       "AI Decks",
				Cards: []*zb_data.DeckCard{
					{CardKey: battleground_proto.CardKey{MouldId: -1}, Amount: 2},
				},
			},
			Type: zb_enums.AIType_MIXED_AI,
		},
	}
}

func TestKeepAlive(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr *loom.Address
	var ctx contract.Context

	// setup ctx
	c = &ZombieBattleground{}
	pubKey, _ := hex.DecodeString(pubKeyHexString)
	addr = &loom.Address{
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}
	now := time.Now()
	fc := plugin.CreateFakeContext(*addr, *addr)
	fc.SetTime(now)
	ctx = contract.WrapPluginContext(fc)
	err := c.Init(ctx, &initRequest)
	assert.Nil(t, err)

	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)

	var matchID int64

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
				DebugCheats: zb_data.DebugCheatsConfiguration{
					Enabled:             true,
					UseCustomRandomSeed: true,
					CustomRandomSeed:    2,
				},
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb_calls.RegisterPlayerPoolRequest{
			RegistrationData: &zb_data.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
				DebugCheats: zb_data.DebugCheatsConfiguration{
					Enabled:             true,
					UseCustomRandomSeed: true,
					CustomRandomSeed:    2,
				},
			},
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb_calls.FindMatchRequest{
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb_calls.AcceptMatchRequest{
			UserId:  "player-2",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb_data.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb_calls.GetGameStateRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
		assert.EqualValues(t, 0, response.GameState.CurrentPlayerIndex)
	})

	t.Run("KeepAlivePlayer1_Success", func(t *testing.T) {
		_, err := c.KeepAlive(ctx, &zb_calls.KeepAliveRequest{
			MatchId: 1,
			UserId:  "player-1",
		})
		assert.Nil(t, err)
	})

	t.Run("KeepAlivePlayer2_Success", func(t *testing.T) {
		_, err := c.KeepAlive(ctx, &zb_calls.KeepAliveRequest{
			MatchId: 1,
			UserId:  "player-2",
		})
		assert.Nil(t, err)
	})

	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb_calls.GetGameStateRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
		assert.EqualValues(t, 0, response.GameState.CurrentPlayerIndex)
	})

	// keep player2 alive
	fc.SetTime(now.Add(time.Second * 5))
	t.Run("KeepAlivePlayer2_Success", func(t *testing.T) {
		_, err := c.KeepAlive(ctx, &zb_calls.KeepAliveRequest{
			MatchId: 1,
			UserId:  "player-2",
		})
		assert.Nil(t, err)
	})

	// player1 fails to keep alive
	// move time forward to expire the player's turn
	fc.SetTime(now.Add(KeepAliveTimeout + time.Second*5))

	t.Run("KeepAlivePlayer2_Winner", func(t *testing.T) {
		_, err := c.KeepAlive(ctx, &zb_calls.KeepAliveRequest{
			MatchId: 1,
			UserId:  "player-2",
		})
		assert.Nil(t, err)
	})

	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb_calls.GetGameStateRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.GameState.IsEnded, "game state should be ended after use leaves the match")
		latestAction := response.GameState.PlayerActions[len(response.GameState.PlayerActions)-1]
		assert.Equal(t, zb_enums.PlayerActionType_LeaveMatch, latestAction.ActionType)
		assert.Equal(t, "player-2", response.GameState.Winner)
	})

	t.Run("KeepAliveAfterGameEndedShouldNot_Fail", func(t *testing.T) {
		_, err := c.KeepAlive(ctx, &zb_calls.KeepAliveRequest{
			MatchId: 1,
			UserId:  "player-1",
		})
		assert.Nil(t, err)
	})

	t.Run("KeepAliveAfterGameEndedShouldNot_Fail", func(t *testing.T) {
		_, err := c.KeepAlive(ctx, &zb_calls.KeepAliveRequest{
			MatchId: 1,
			UserId:  "player-2",
		})
		assert.Nil(t, err)
	})

	t.Run("SendAnyActionShould_Fail", func(t *testing.T) {
		_, err := c.SendPlayerAction(ctx, &zb_calls.PlayerActionRequest{
			MatchId: 1,
			PlayerAction: &zb_data.PlayerAction{
				ActionType: zb_enums.PlayerActionType_EndTurn,
				PlayerId:   "player-2",
			},
		})
		assert.NotNil(t, err)
	})
}

func TestHashSignature(t *testing.T) {
	privateKeyStr := "921660bf3e5c7a404beed663f00462645fd8d50751d21e262f6f1a3b7e5b5da3"
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	assert.Nil(t, err)

	generator, err := NewMintingReceiptGenerator(privateKey, 1)
	assert.Nil(t, err)
	verifySignResult, err := generator.generateEosVerifySignResult(
		big.NewInt(5),
		generator.gatewayPrivateKey,
		3,
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
		big.NewInt(373),
		generator.contractVersion)

	assert.Nil(t, err)
	assert.Equal(t, "0x4cc40300d439c52ca6bb577ac42d883fff58df81309c8a41d7a9f54a1641c110", verifySignResult.Hash)
	assert.Equal(t, "0x6f27beb227b13842d994360d549f16242cfe474a147da9a391b80460cb3ec291047ee645a1a5c634287368b2d582a8f935ac3f0348840a71d9a015e8e6e7a1d11c", verifySignResult.Signature)
}
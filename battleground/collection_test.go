package battleground

import (
	battleground_proto "github.com/loomnetwork/gamechain/battleground/proto"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/go-loom"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	assert "github.com/stretchr/testify/require"
	"testing"
)

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

	t.Run("GetCollection should return default card collection", func(t *testing.T) {
		CardCollectionCard, err := c.GetCollection(ctx, &zb_calls.GetCollectionRequest{
			UserId:  "CardUser",
			Version: "v1",
		})
		assert.Nil(t, err)
		assert.Equal(t, 8, len(CardCollectionCard.Cards))
	})

	t.Run("GetCollection should not return unknown card stored in card collection", func(t *testing.T) {
		const unknownCardMouldId = 999999

		// Add unknown card to collection
		collection, err := loadUserCardCollectionRaw(ctx, "CardUser")
		assert.Nil(t, err)
		collection.Cards = append(collection.Cards, &zb_data.CardCollectionCard{
			CardKey: battleground_proto.CardKey{
				MouldId: unknownCardMouldId,
			},
			Amount: 3,
		})
		err = saveUserCardCollection(ctx, "CardUser", collection)
		assert.Nil(t, err)

		// Check
		getCollectionResponse, err := c.GetCollection(ctx, &zb_calls.GetCollectionRequest{
			UserId:  "CardUser",
			Version: "v1",
		})
		assert.Nil(t, err)

		for _, collectionCard := range getCollectionResponse.Cards {
			if collectionCard.CardKey.MouldId == unknownCardMouldId {
				assert.Fail(t, "unknown card returned in user collection")
			}
		}
	})

	t.Run("GetCollection should return fake full collection", func(t *testing.T) {
		configuration, err := loadContractConfiguration(ctx)
		assert.Nil(t, err)
		configuration.UseCardLibraryAsUserCollection = true
		err = saveContractConfiguration(ctx, configuration)
		assert.Nil(t, err)

		CardCollectionCard, err := c.GetCollection(ctx, &zb_calls.GetCollectionRequest{
			UserId:  "CardUser",
			Version: "v1",
		})
		assert.Nil(t, err)
		assert.Equal(t, 100, len(CardCollectionCard.Cards))
	})
}

func TestSyncCardToCollection(t *testing.T) {
	c := &ZombieBattleground{}
	collection := c.syncCardAmountChangesToCollection(
		[]*zb_data.CardCollectionCard{
			{
				CardKey: battleground_proto.CardKey{MouldId: 2},
				Amount:  3,
			},
			{
				CardKey: battleground_proto.CardKey{MouldId: 5},
				Amount:  2,
			},
		},
		[]CardAmountChangeItem{
			{
				CardKey:      battleground_proto.CardKey{MouldId: 2},
				AmountChange: 1,
			},
			{
				CardKey:      battleground_proto.CardKey{MouldId: 5},
				AmountChange: -2,
			},
		},
	)

	assert.NotNil(t, collection)
	assert.Equal(t, 1, len(collection))
	assert.Equal(t, int64(2), collection[0].CardKey.MouldId)
	assert.Equal(t, int64(4), collection[0].Amount)
}

func TestLimitDeckByCardCollection(t *testing.T) {
	deck := &zb_data.Deck{
		Cards: []*zb_data.DeckCard{
			{
				CardKey: battleground_proto.CardKey{MouldId: 2},
				Amount:  3,
			},
			{
				CardKey: battleground_proto.CardKey{MouldId: 5},
				Amount:  4,
			},
			{
				CardKey: battleground_proto.CardKey{MouldId: 7},
				Amount:  4,
			},
		},
	}

	t.Run("Collection shouldn't change", func(t *testing.T) {
		collectionCards := []*zb_data.CardCollectionCard{
			{
				CardKey: battleground_proto.CardKey{MouldId: 2},
				Amount:  3,
			},
			{
				CardKey: battleground_proto.CardKey{MouldId: 5},
				Amount:  4,
			},
			{
				CardKey: battleground_proto.CardKey{MouldId: 7},
				Amount:  4,
			},
		}

		limitDeckByCardCollection(deck, collectionCards)

		assert.Equal(t, 3, len(deck.Cards))
		assert.Equal(t, int64(2), deck.Cards[0].CardKey.MouldId)
		assert.Equal(t, int64(3), deck.Cards[0].Amount)
		assert.Equal(t, int64(5), deck.Cards[1].CardKey.MouldId)
		assert.Equal(t, int64(4), deck.Cards[1].Amount)
		assert.Equal(t, int64(7), deck.Cards[2].CardKey.MouldId)
		assert.Equal(t, int64(4), deck.Cards[2].Amount)
	})

	t.Run("Deck shouldn't contain more cards than exists in collection", func(t *testing.T) {
		collectionCards := []*zb_data.CardCollectionCard{
			{
				CardKey: battleground_proto.CardKey{MouldId: 2},
				Amount:  3,
			},
			{
				CardKey: battleground_proto.CardKey{MouldId: 5},
				Amount:  2,
			},
		}

		limitDeckByCardCollection(deck, collectionCards)

		assert.Equal(t, 2, len(deck.Cards))
		assert.Equal(t, int64(2), deck.Cards[0].CardKey.MouldId)
		assert.Equal(t, int64(3), deck.Cards[0].Amount)
		assert.Equal(t, int64(5), deck.Cards[1].CardKey.MouldId)
		assert.Equal(t, int64(2), deck.Cards[1].Amount)
	})
}

func TestDebugCheatSetFullCardCollection(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "CardUser",
		Image:   "PathToImage",
		Version: "v1",
	}, t)

	collection, err := loadUserCardCollectionRaw(ctx, "CardUser")
	assert.Nil(t, err)
	assert.Equal(t, 8, len(collection.Cards))

	_, err = c.DebugCheatSetFullCardCollection(ctx, &zb_calls.DebugCheatSetFullCardCollectionRequest{
		UserId:  "CardUser",
		Version: "v1",
	})
	assert.Nil(t, err)

	collection, err = loadUserCardCollectionRaw(ctx, "CardUser")
	assert.Nil(t, err)
	assert.Equal(t, 100, len(collection.Cards))
}

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
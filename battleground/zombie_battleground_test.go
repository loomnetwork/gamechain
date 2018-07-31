package battleground

import (
	"encoding/hex"
	"testing"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/stretchr/testify/assert"
)

var initRequest = zb.InitRequest{
	DefaultCollection: []*zb.CardCollection{
		{
			CardId: 1,
			Amount: 4,
		},
		{
			CardId: 2,
			Amount: 3,
		},
		{
			CardId: 3,
			Amount: 5,
		},
		{
			CardId: 4,
			Amount: 4,
		},
		{
			CardId: 5,
			Amount: 3,
		},
		{
			CardId: 6,
			Amount: 5,
		},
		{
			CardId: 7,
			Amount: 4,
		},
		{
			CardId: 8,
			Amount: 3,
		},
		{
			CardId: 9,
			Amount: 5,
		},
		{
			CardId: 10,
			Amount: 4,
		},
		{
			CardId: 11,
			Amount: 3,
		},
		{
			CardId: 12,
			Amount: 5,
		},
	},
	Heroes: []*zb.Hero{
		{
			HeroId:  0,
			Icon:    "asdasd",
			Name:    "Golem Hero",
			Element: 2,
			Skills: []*zb.Skill{{
				Title:           "Deffence",
				SkillType:       4,
				SkillTargetType: 0,
				Cost:            2,
				Value:           2,
			}},
		},
		{
			HeroId:  1,
			Icon:    "asdasd",
			Name:    "Pyro Hero",
			Element: 0,
			Skills: []*zb.Skill{{
				Title:           "Fireball",
				SkillType:       2,
				SkillTargetType: 0,
				Cost:            2,
				Value:           1,
			}},
		},
	},
	Cards: []*zb.Card{
		{
			Id:      1,
			Element: "Air",
			Name:    "Banshee",
			Rank:    "Minion",
			Type:    "Feral",
			Damage:  2,
			Health:  1,
			Cost:    2,
			Ability: "Feral",
			Effects: []*zb.Effect{
				{
					Trigger:  "entry",
					Effect:   "feral",
					Duration: "permanent",
					Target:   "self",
				},
			},
			CardViewInfo: &zb.CardViewInfo{
				Position: &zb.Coordinates{
					X: 1.5,
					Y: 2.5,
					Z: 3.5,
				},
				Scale: &zb.Coordinates{
					X: 0.5,
					Y: 0.5,
					Z: 0.5,
				},
			},
		},
		{
			Id:      2,
			Element: "Air",
			Name:    "Breezee",
			Rank:    "Minion",
			Type:    "Walker",
			Damage:  1,
			Health:  1,
			Cost:    1,
			Ability: "-",
			Effects: []*zb.Effect{
				{
					Trigger: "death",
					Effect:  "attack_strength_buff",
					Target:  "friendly_selectable",
				},
			},
		},
	},
	DefaultDecks: []*zb.Deck{
		{
			HeroId: 1,
			Name:   "Default",
			Cards: []*zb.CardCollection{
				{
					CardId: 1,
					Amount: 2,
				},
				{
					CardId: 2,
					Amount: 2,
				},
				{
					CardId: 3,
					Amount: 2,
				},
				{
					CardId: 4,
					Amount: 2,
				},
				{
					CardId: 5,
					Amount: 2,
				},
				{
					CardId: 6,
					Amount: 2,
				},
				{
					CardId: 7,
					Amount: 1,
				},
				{
					CardId: 8,
					Amount: 1,
				},
				{
					CardId: 9,
					Amount: 1,
				},
				{
					CardId: 10,
					Amount: 1,
				},
				{
					CardId: 11,
					Amount: 1,
				},
				{
					CardId: 12,
					Amount: 1,
				},
			},
		},
	},
}

func setup(c *ZombieBattleground, pubKeyHex string, addr *loom.Address, ctx *contract.Context, t *testing.T) {

	c = &ZombieBattleground{}
	pubKey, _ := hex.DecodeString(pubKeyHex)

	addr = &loom.Address{
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}

	*ctx = contract.WrapPluginContext(
		plugin.CreateFakeContext(*addr, *addr),
	)

	err := c.Init(*ctx, &initRequest)
	assert.Nil(t, err)
}

func setupAccount(c *ZombieBattleground, ctx contract.Context, upsertAccountRequest *zb.UpsertAccountRequest, t *testing.T) {
	err := c.CreateAccount(ctx, upsertAccountRequest)
	assert.Nil(t, err)
}

func TestAccountOperations(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId: "AccountUser",
		Image:  "PathToImage",
	}, t)

	t.Run("UpdateAccount", func(t *testing.T) {
		account, err := c.UpdateAccount(ctx, &zb.UpsertAccountRequest{
			UserId:      "AccountUser",
			Image:       "PathToImage2",
			CurrentTier: 5,
		})
		assert.Equal(t, nil, err)
		assert.Equal(t, int32(5), account.CurrentTier)
		assert.Equal(t, "PathToImage2", account.Image)
	})

	t.Run("GetAccount", func(t *testing.T) {
		account, err := c.GetAccount(ctx, &zb.GetAccountRequest{
			UserId: "AccountUser",
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(5), account.CurrentTier)
		assert.Equal(t, "PathToImage2", account.Image)
	})
}

func TestCardCollectionOperations(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "8996b813617b283f81ea1747fbddbe73fe4b5fce0eac0728e47de51d8e506701"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId: "CardUser",
		Image:  "PathToImage",
	}, t)

	cardCollection, err := c.GetCollection(ctx, &zb.GetCollectionRequest{
		UserId: "CardUser",
	})
	assert.Nil(t, err)
	assert.Equal(t, 12, len(cardCollection.Cards))

}

func TestDeckOperations(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "7796b813617b283f81ea1747fbddbe73fe4b5fce0eac0728e47de51d8e506701"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId: "DeckUser",
		Image:  "PathToImage",
	}, t)

	t.Run("ListDecks", func(t *testing.T) {
		deckResponse, err := c.ListDecks(ctx, &zb.ListDecksRequest{
			UserId: "DeckUser",
		})
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(deckResponse.Decks))
		assert.Equal(t, "Default", deckResponse.Decks[0].Name)
	})

	t.Run("GetDeck (Not Exists)", func(t *testing.T) {
		deckResponse, err := c.GetDeck(ctx, &zb.GetDeckRequest{
			UserId:   "DeckUser",
			DeckName: "NotExists",
		})
		assert.Equal(t, (*zb.GetDeckResponse)(nil), deckResponse)
		assert.Equal(t, contract.ErrNotFound, err)
	})

	t.Run("GetDeck", func(t *testing.T) {
		deckResponse, err := c.GetDeck(ctx, &zb.GetDeckRequest{
			UserId:   "DeckUser",
			DeckName: "Default",
		})
		assert.Nil(t, err)
		assert.Equal(t, "Default", deckResponse.Deck.Name)
	})

	t.Run("AddDeck", func(t *testing.T) {
		err := c.CreateDeck(ctx, &zb.CreateDeckRequest{
			UserId: "DeckUser",
			Deck: &zb.Deck{
				Name:   "NewDeck",
				HeroId: 1,
				Cards: []*zb.CardCollection{
					{
						Amount: 1,
						CardId: 2,
					},
					{
						Amount: 1,
						CardId: 3,
					},
				},
			},
		})

		assert.Nil(t, err)

		deckResponse, err := c.ListDecks(ctx, &zb.ListDecksRequest{
			UserId: "DeckUser",
		})

		assert.Equal(t, nil, err)
		assert.Equal(t, 2, len(deckResponse.Decks))
	})

	t.Run("AddDeck (Invalid Requested Amount)", func(t *testing.T) {
		err := c.CreateDeck(ctx, &zb.CreateDeckRequest{
			UserId: "DeckUser",
			Deck: &zb.Deck{
				Name:   "NewDeck",
				HeroId: 1,
				Cards: []*zb.CardCollection{
					{
						Amount: 200,
						CardId: 2,
					},
					{
						Amount: 100,
						CardId: 3,
					},
				},
			},
		})

		assert.NotNil(t, err)
	})

	t.Run("AddDeck (Invalid Requested CardId)", func(t *testing.T) {
		err := c.CreateDeck(ctx, &zb.CreateDeckRequest{
			UserId: "DeckUser",
			Deck: &zb.Deck{
				Name:   "NewDeck",
				HeroId: 1,
				Cards: []*zb.CardCollection{
					{
						Amount: 2,
						CardId: 234,
					},
					{
						Amount: 1,
						CardId: 345,
					},
				},
			},
		})

		assert.NotNil(t, err)
	})

	t.Run("DeleteDeck", func(t *testing.T) {
		err := c.DeleteDeck(ctx, &zb.DeleteDeckRequest{
			UserId:   "DeckUser",
			DeckName: "NewDeck",
		})

		assert.Nil(t, err)
	})

	t.Run("DeleteDeck (Non existant)", func(t *testing.T) {
		err := c.DeleteDeck(ctx, &zb.DeleteDeckRequest{
			UserId:   "DeckUser",
			DeckName: "NotExists",
		})

		assert.NotNil(t, err)
	})
}

func TestCardOperations(t *testing.T) {

	var c *ZombieBattleground
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	t.Run("ListCard", func(t *testing.T) {
		cardResponse, err := c.ListCardLibrary(ctx, &zb.ListCardLibraryRequest{})

		assert.Nil(t, err)
		assert.Equal(t, 2, len(cardResponse.Sets[0].Cards))
	})

	t.Run("ListHero", func(t *testing.T) {
		heroResponse, err := c.ListHero(ctx, &zb.ListHeroRequest{})

		assert.Nil(t, err)
		assert.Equal(t, 2, len(heroResponse.Heroes))
	})
}

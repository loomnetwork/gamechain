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
			CardName: "Banshee",
			Amount:   4,
		},
		{
			CardName: "Breezee",
			Amount:   3,
		},
		{
			CardName: "Buffer",
			Amount:   5,
		},
		{
			CardName: "Soothsayer",
			Amount:   4,
		},
		{
			CardName: "Wheezy",
			Amount:   3,
		},
		{
			CardName: "Whiffer",
			Amount:   5,
		},
		{
			CardName: "Whizpar",
			Amount:   4,
		},
		{
			CardName: "Zhocker",
			Amount:   3,
		},
		{
			CardName: "Bouncer",
			Amount:   5,
		},
		{
			CardName: "Dragger",
			Amount:   4,
		},
		{
			CardName: "Guzt",
			Amount:   3,
		},
		{
			CardName: "Pushhh",
			Amount:   5,
		},
	},
	DefaultHeroes: []*zb.Hero{
		{
			HeroId:     0,
			Experience: 0,
			Level:      1,
			Skills: []*zb.Skill{{
				Title:           "Attack",
				Skill:           3,
				SkillTargetType: 0,
				Value:           1,
			}},
		},
		{
			HeroId:     1,
			Experience: 0,
			Level:      2,
			Skills: []*zb.Skill{{
				Title:           "Deffence",
				Skill:           4,
				SkillTargetType: 0,
				Value:           2,
			}},
		},
	},
	Heroes: []*zb.Hero{
		{
			HeroId:  0,
			Icon:    "asdasd",
			Name:    "Golem Hero",
			Element: 2,
		},
		{
			HeroId:  1,
			Icon:    "asdasd",
			Name:    "Pyro Hero",
			Element: 0,
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
			HeroId: 2,
			Name:   "Default",
			Cards: []*zb.CardCollection{
				{
					CardName: "Banshee",
					Amount:   2,
				},
				{
					CardName: "Breezee",
					Amount:   2,
				},
				{
					CardName: "Buffer",
					Amount:   2,
				},
				{
					CardName: "Soothsayer",
					Amount:   2,
				},
				{
					CardName: "Wheezy",
					Amount:   2,
				},
				{
					CardName: "Whiffer",
					Amount:   2,
				},
				{
					CardName: "Whizpar",
					Amount:   1,
				},
				{
					CardName: "Zhocker",
					Amount:   1,
				},
				{
					CardName: "Bouncer",
					Amount:   1,
				},
				{
					CardName: "Dragger",
					Amount:   1,
				},
				{
					CardName: "Guzt",
					Amount:   1,
				},
				{
					CardName: "Pushhh",
					Amount:   1,
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
						Amount:   1,
						CardName: "Breezee",
					},
					{
						Amount:   1,
						CardName: "Buffer",
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
						Amount:   200,
						CardName: "Breezee",
					},
					{
						Amount:   100,
						CardName: "Buffer",
					},
				},
			},
		})

		assert.NotNil(t, err)
	})

	t.Run("AddDeck (Invalid Requested CardName)", func(t *testing.T) {
		err := c.CreateDeck(ctx, &zb.CreateDeckRequest{
			UserId: "DeckUser",
			Deck: &zb.Deck{
				Name:   "NewDeck",
				HeroId: 1,
				Cards: []*zb.CardCollection{
					{
						Amount:   2,
						CardName: "InvalidName1",
					},
					{
						Amount:   1,
						CardName: "InvalidName2",
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
		heroResponse, err := c.ListHeroes(ctx, &zb.ListHeroRequest{})

		assert.Nil(t, err)
		assert.Equal(t, 2, len(heroResponse.Heroes))
	})
}

func TestHeroOperations(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "7696b824516b283f81ea1747fbddbe73fe4b5fce0eac0728e47de41d8e306701"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId: "HeroUser",
		Image:  "PathToImage",
	}, t)

	t.Run("ListHeroesForUser", func(t *testing.T) {
		heroesResponse, err := c.ListHeroesForUser(ctx, &zb.ListHeroesForUserRequest{
			UserId: "HeroUser",
		})

		assert.Nil(t, err)
		assert.Equal(t, 2, len(heroesResponse.Heroes))
	})

	t.Run("GetHeroForUser", func(t *testing.T) {
		heroResponse, err := c.GetHeroForUser(ctx, &zb.GetHeroForUserRequest{
			UserId: "HeroUser",
			HeroId: 1,
		})

		assert.Nil(t, err)
		assert.NotNil(t, heroResponse.Hero)
	})

	t.Run("GetHeroForUser (Hero not exists)", func(t *testing.T) {
		_, err := c.GetHeroForUser(ctx, &zb.GetHeroForUserRequest{
			UserId: "HeroUser",
			HeroId: 10,
		})

		assert.NotNil(t, err)
	})

	t.Run("AddHeroExperience", func(t *testing.T) {
		resp, err := c.AddHeroExperience(ctx, &zb.AddHeroExperienceRequest{
			UserId:     "HeroUser",
			HeroId:     0,
			Experience: 2,
		})

		assert.Nil(t, err)
		assert.Equal(t, int64(2), resp.Experience)
	})

	t.Run("AddHeroExperience (Negative experience)", func(t *testing.T) {
		_, err := c.AddHeroExperience(ctx, &zb.AddHeroExperienceRequest{
			UserId:     "HeroUser",
			HeroId:     0,
			Experience: -2,
		})

		assert.NotNil(t, err)
	})

	t.Run("AddHeroExperience (Non existant hero)", func(t *testing.T) {
		_, err := c.AddHeroExperience(ctx, &zb.AddHeroExperienceRequest{
			UserId:     "HeroUser",
			HeroId:     100,
			Experience: 2,
		})

		assert.NotNil(t, err)
	})

	t.Run("GetHeroSkills", func(t *testing.T) {
		skillResponse, err := c.GetHeroSkills(ctx, &zb.GetHeroSkillsRequest{
			UserId: "HeroUser",
			HeroId: 0,
		})

		assert.Nil(t, err)
		assert.Equal(t, 1, len(skillResponse.Skills))
	})

	t.Run("GetHeroSkills (Non existant hero)", func(t *testing.T) {
		_, err := c.GetHeroSkills(ctx, &zb.GetHeroSkillsRequest{
			UserId: "HeroUser",
			HeroId: 100,
		})

		assert.NotNil(t, err)
	})

}

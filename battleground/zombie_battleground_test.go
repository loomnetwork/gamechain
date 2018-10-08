package battleground

import (
	"encoding/hex"
	"testing"

	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/stretchr/testify/assert"
)

var initRequest = zb.InitRequest{
	Version: "v1",
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
	Heroes: []*zb.Hero{
		{
			HeroId:     0,
			Experience: 0,
			Level:      1,
			Skills: []*zb.Skill{{
				Title:        "Attack",
				Skill:        "Skill0",
				SkillTargets: "zb.Skill_ALL_CARDS|zb.Skill_PLAYER_CARD",
				Value:        1,
			}},
		},
		{
			HeroId:     1,
			Experience: 0,
			Level:      2,
			Skills: []*zb.Skill{{
				Title:        "Deffence",
				Skill:        "Skill1",
				SkillTargets: "zb.Skill_PLAYER|zb.Skill_OPPONENT_CARD",
				Value:        2,
			}},
		},
	},
	Cards: []*zb.Card{
		{
			Id:      1,
			Set:     "Air",
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
			Set:     "Air",
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
			Id:     0,
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

var updateInitRequest = zb.UpdateInitRequest{
	Version: "v2",
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
	Heroes: []*zb.Hero{
		{
			HeroId:     0,
			Experience: 0,
			Level:      1,
			Skills: []*zb.Skill{{
				Title:        "Attack",
				Skill:        "Skill0",
				SkillTargets: "zb.Skill_ALL_CARDS|zb.Skill_PLAYER_CARD",
				Value:        1,
			}},
		},
		{
			HeroId:     1,
			Experience: 0,
			Level:      2,
			Skills: []*zb.Skill{{
				Title:        "Deffence",
				Skill:        "Skill1",
				SkillTargets: "zb.Skill_PLAYER|zb.Skill_OPPONENT_CARD",
				Value:        2,
			}},
		},
	},
	Cards: []*zb.Card{
		{
			Id:      1,
			Set:     "Air",
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
			Set:     "Air",
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
		{
			Id:      3,
			Set:     "Air",
			Name:    "NewCard",
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
			Id:     0,
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

var updateCardListRequest = zb.UpdateCardListRequest{
	Version: "v2",
	Cards: []*zb.Card{
		{
			Id:      1,
			Set:     "Air",
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
			Set:     "Air",
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
		{
			Id:      3,
			Set:     "Air",
			Name:    "NewCard",
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
		UserId:  "AccountUser",
		Image:   "PathToImage",
		Version: "v1",
	}, t)

	t.Run("UpdateAccount", func(t *testing.T) {
		account, err := c.UpdateAccount(ctx, &zb.UpsertAccountRequest{
			UserId:      "AccountUser",
			Image:       "PathToImage2",
			CurrentTier: 5,
		})
		assert.Nil(t, err)
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
		UserId:  "CardUser",
		Image:   "PathToImage",
		Version: "v1",
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
		UserId:  "DeckUser",
		Image:   "PathToImage",
		Version: "v1",
	}, t)

	t.Run("ListDecks", func(t *testing.T) {
		deckResponse, err := c.ListDecks(ctx, &zb.ListDecksRequest{
			UserId: "DeckUser",
		})
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(deckResponse.Decks))
		assert.Equal(t, int64(1), deckResponse.Decks[0].Id)
		assert.Equal(t, "Default", deckResponse.Decks[0].Name)
	})

	t.Run("GetDeck (Not Exists)", func(t *testing.T) {
		deckResponse, err := c.GetDeck(ctx, &zb.GetDeckRequest{
			UserId: "DeckUser",
			DeckId: 0xDEADBEEF,
		})
		assert.Equal(t, (*zb.GetDeckResponse)(nil), deckResponse)
		assert.Equal(t, contract.ErrNotFound, err)
	})

	t.Run("GetDeck", func(t *testing.T) {
		deckResponse, err := c.GetDeck(ctx, &zb.GetDeckRequest{
			UserId: "DeckUser",
			DeckId: 1,
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), deckResponse.Deck.Id) // id should start from 1
		assert.Equal(t, "Default", deckResponse.Deck.Name)
	})

	var createDeckResponse *zb.CreateDeckResponse
	t.Run("CreateDeck", func(t *testing.T) {
		var err error
		createDeckResponse, err = c.CreateDeck(ctx, &zb.CreateDeckRequest{
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
			Version: "v1",
		})

		assert.NotNil(t, createDeckResponse)
		assert.Nil(t, err)

		deckResponse, err := c.ListDecks(ctx, &zb.ListDecksRequest{
			UserId: "DeckUser",
		})

		assert.Equal(t, nil, err)
		assert.Equal(t, 2, len(deckResponse.Decks))
	})

	t.Run("CreateDeck (Invalid Requested Amount)", func(t *testing.T) {
		_, err := c.CreateDeck(ctx, &zb.CreateDeckRequest{
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
			Version: "v1",
		})

		assert.NotNil(t, err)
	})

	t.Run("CreateDeck (Invalid Requested CardName)", func(t *testing.T) {
		_, err := c.CreateDeck(ctx, &zb.CreateDeckRequest{
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
			Version: "v1",
		})

		assert.NotNil(t, err)
	})

	t.Run("CreateDeck (Same name not allowed)", func(t *testing.T) {
		_, err := c.CreateDeck(ctx, &zb.CreateDeckRequest{
			UserId: "DeckUser",
			Deck: &zb.Deck{
				Name:   "Default",
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
			Version: "v1",
		})

		assert.NotNil(t, err)
	})

	t.Run("CreateDeck (Same name with different case not allowed)", func(t *testing.T) {
		_, err := c.CreateDeck(ctx, &zb.CreateDeckRequest{
			UserId: "DeckUser",
			Deck: &zb.Deck{
				Name:   "nEWdECK",
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
			Version: "v1",
		})

		assert.NotNil(t, err)
	})

	t.Run("EditDeck", func(t *testing.T) {
		err := c.EditDeck(ctx, &zb.EditDeckRequest{
			UserId: "DeckUser",
			Deck: &zb.Deck{
				Id:     2,
				Name:   "Edited",
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
			Version: "v1",
		})

		assert.Nil(t, err)

		getDeckResponse, err := c.GetDeck(ctx, &zb.GetDeckRequest{
			UserId: "DeckUser",
			DeckId: 2,
		})

		assert.Nil(t, err)

		assert.Equal(t, "Edited", getDeckResponse.Deck.Name)
	})

	t.Run("EditDeck (attempt to set more number of cards)", func(t *testing.T) {
		err := c.EditDeck(ctx, &zb.EditDeckRequest{
			UserId: "DeckUser",
			Deck: &zb.Deck{
				Id:     2,
				Name:   "Edited",
				HeroId: 1,
				Cards: []*zb.CardCollection{
					{
						Amount:   100,
						CardName: "Breezee",
					},
					{
						Amount:   1,
						CardName: "Buffer",
					},
				},
			},
			Version: "v1",
		})

		assert.NotNil(t, err)
	})

	t.Run("EditDeck (same name while editing is allowed)", func(t *testing.T) {
		err := c.EditDeck(ctx, &zb.EditDeckRequest{
			UserId: "DeckUser",
			Deck: &zb.Deck{
				Id:     2,
				Name:   "Edited",
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
			Version: "v1",
		})

		assert.Nil(t, err)
	})

	t.Run("EditDeck (attempt to set duplicate name with different case)", func(t *testing.T) {
		err := c.EditDeck(ctx, &zb.EditDeckRequest{
			UserId: "DeckUser",
			Deck: &zb.Deck{
				Id:     2,
				Name:   "dEFAULT",
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
			Version: "v1",
		})

		assert.NotNil(t, err)
	})

	t.Run("DeleteDeck", func(t *testing.T) {
		assert.NotNil(t, createDeckResponse)
		err := c.DeleteDeck(ctx, &zb.DeleteDeckRequest{
			UserId: "DeckUser",
			DeckId: createDeckResponse.DeckId,
		})

		assert.Nil(t, err)
	})

	t.Run("DeleteDeck (Non existant)", func(t *testing.T) {
		err := c.DeleteDeck(ctx, &zb.DeleteDeckRequest{
			UserId: "DeckUser",
			DeckId: 0xDEADBEEF,
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
		cardResponse, err := c.ListCardLibrary(ctx, &zb.ListCardLibraryRequest{
			Version: "v1",
		})

		assert.Nil(t, err)
		assert.Equal(t, 2, len(cardResponse.Sets[0].Cards))
	})

	t.Run("ListHeroLibrary", func(t *testing.T) {
		heroResponse, err := c.ListHeroLibrary(ctx, &zb.ListHeroLibraryRequest{
			Version: "v1",
		})

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
		UserId:  "HeroUser",
		Image:   "PathToImage",
		Version: "v1",
	}, t)

	t.Run("ListHeroes", func(t *testing.T) {
		heroesResponse, err := c.ListHeroes(ctx, &zb.ListHeroesRequest{
			UserId: "HeroUser",
		})

		assert.Nil(t, err)
		assert.Equal(t, 2, len(heroesResponse.Heroes))
	})

	t.Run("GetHero", func(t *testing.T) {
		heroResponse, err := c.GetHero(ctx, &zb.GetHeroRequest{
			UserId: "HeroUser",
			HeroId: 1,
		})

		assert.Nil(t, err)
		assert.NotNil(t, heroResponse.Hero)
	})

	t.Run("GetHero (Hero not exists)", func(t *testing.T) {
		_, err := c.GetHero(ctx, &zb.GetHeroRequest{
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

func TestUpdateInitDataOperations(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	t.Run("UpdateInit", func(t *testing.T) {
		err := c.UpdateInit(ctx, &updateInitRequest)

		assert.Nil(t, err)
	})
}

func TestUpdateCardListOperations(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	t.Run("UpdateCardList", func(t *testing.T) {
		err := c.UpdateCardList(ctx, &updateCardListRequest)

		assert.Nil(t, err)
	})
}

func TestFindMatchOperations(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)

	// make users have decks
	t.Run("ListDecksPlayer1", func(t *testing.T) {
		deckResponse, err := c.ListDecks(ctx, &zb.ListDecksRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(deckResponse.Decks))
		assert.Equal(t, int64(1), deckResponse.Decks[0].Id)
		assert.Equal(t, "Default", deckResponse.Decks[0].Name)
	})
	t.Run("ListDecksPlayer2", func(t *testing.T) {
		deckResponse, err := c.ListDecks(ctx, &zb.ListDecksRequest{
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(deckResponse.Decks))
		assert.Equal(t, int64(1), deckResponse.Decks[0].Id)
		assert.Equal(t, "Default", deckResponse.Decks[0].Name)
	})

	var matchID int64

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			DeckId: 1,
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, len(response.Match.PlayerStates), "the first player should see only 1 player state")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			DeckId: 1,
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the second player should 2 player states")
		assert.Equal(t, zb.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("GetMatch", func(t *testing.T) {
		response, err := c.GetMatch(ctx, &zb.GetMatchRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the second player should 2 player states")
		assert.Equal(t, zb.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.NotNil(t, response.GameState)
	})

	t.Run("EndMatch", func(t *testing.T) {
		_, err := c.EndMatch(ctx, &zb.EndMatchRequest{
			MatchId:  matchID,
			UserId:   "player-1",
			WinnerId: "player-2",
		})
		assert.Nil(t, err)
	})

	t.Run("GetMatchAfterLeaving", func(t *testing.T) {
		response, err := c.GetMatch(ctx, &zb.GetMatchRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the second player should 2 player states")
		assert.Equal(t, zb.Match_Ended, response.Match.Status, "match status should be 'ended'")
		assert.NotNil(t, response.GameState)
	})
}

func TestGameStateOperations(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)

	// make users have decks
	t.Run("ListDecksPlayer1", func(t *testing.T) {
		deckResponse, err := c.ListDecks(ctx, &zb.ListDecksRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(deckResponse.Decks))
		assert.Equal(t, int64(1), deckResponse.Decks[0].Id)
		assert.Equal(t, "Default", deckResponse.Decks[0].Name)
	})
	t.Run("ListDecksPlayer2", func(t *testing.T) {
		deckResponse, err := c.ListDecks(ctx, &zb.ListDecksRequest{
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(deckResponse.Decks))
		assert.Equal(t, int64(1), deckResponse.Decks[0].Id)
		assert.Equal(t, "Default", deckResponse.Decks[0].Name)
	})

	var matchID int64

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			DeckId: 1,
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, len(response.Match.PlayerStates), "the first player should see only 1 player state")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			DeckId: 1,
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the second player should 2 player states")
		assert.Equal(t, zb.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("GetMatch", func(t *testing.T) {
		response, err := c.GetMatch(ctx, &zb.GetMatchRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the second player should 2 player states")
		assert.Equal(t, zb.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.NotNil(t, response.GameState)
	})

	// Note: since the toss coin seed is always 0 for testing, we always get 0 as the first player
	t.Run("SendEndturnPlayer2_Failed", func(t *testing.T) {
		_, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_EndTurn,
				PlayerId:   "player-2",
			},
		})
		assert.NotNil(t, err)
		assert.Equal(t, errInvalidPlayer, err)
	})

	t.Run("SendEndturnPlayer1_Success", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_EndTurn,
				PlayerId:   "player-1",
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.EqualValues(t, 1, response.GameState.CurrentPlayerIndex, "player-2 should be active")
	})
	t.Run("SendEndturnPlayer2_Success", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_EndTurn,
				PlayerId:   "player-2",
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.EqualValues(t, 0, response.GameState.CurrentPlayerIndex, "player-1 should be active")
	})
	t.Run("SendCardPlayPlayer1", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_CardPlay,
				PlayerId:   "player-1",
				Action: &zb.PlayerAction_CardPlay{
					CardPlay: &zb.PlayerActionCardPlay{
						Card: &zb.CardInstance{
							InstanceId: 1,
						},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("SendCardAttackPlayer1", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_CardAttack,
				PlayerId:   "player-1",
				Action: &zb.PlayerAction_CardAttack{
					CardAttack: &zb.PlayerActionCardAttack{
						Attacker: &zb.CardInstance{
							InstanceId: 1,
						},
						AffectObjectType: zb.AffectObjectType_CARD,
						Target: &zb.Unit{
							InstanceId: 2,
						},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("SendCardAbilityPlayer1", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_CardAbilityUsed,
				PlayerId:   "player-1",
				Action: &zb.PlayerAction_CardAbilityUsed{
					CardAbilityUsed: &zb.PlayerActionCardAbilityUsed{
						Card: &zb.CardInstance{
							InstanceId: 1,
						},
						Targets: []*zb.Unit{
							&zb.Unit{
								InstanceId:       2,
								AffectObjectType: zb.AffectObjectType_CARD,
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
		response, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_OverlordSkillUsed,
				PlayerId:   "player-1",
				Action: &zb.PlayerAction_OverlordSkillUsed{
					OverlordSkillUsed: &zb.PlayerActionOverlordSkillUsed{
						SkillId:          1,
						AffectObjectType: zb.AffectObjectType_CARD,
						Target: &zb.Unit{
							InstanceId: 2,
						},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("SendEndturnPlayer1_Success2", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_EndTurn,
				PlayerId:   "player-1",
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.EqualValues(t, 1, response.GameState.CurrentPlayerIndex, "player-2 should be active")
	})
	t.Run("SendCardPlayPlayer2", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_CardPlay,
				PlayerId:   "player-2",
				Action: &zb.PlayerAction_CardPlay{
					CardPlay: &zb.PlayerActionCardPlay{
						Card: &zb.CardInstance{
							InstanceId: 1,
						},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("SendCardAttackPlayer2", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_CardAttack,
				PlayerId:   "player-2",
				Action: &zb.PlayerAction_CardAttack{
					CardAttack: &zb.PlayerActionCardAttack{
						Attacker: &zb.CardInstance{
							InstanceId: 1,
						},
						AffectObjectType: zb.AffectObjectType_CARD,
						Target: &zb.Unit{
							InstanceId: 2,
						},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("SendCardAbilityPlayer2", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_CardAbilityUsed,
				PlayerId:   "player-2",
				Action: &zb.PlayerAction_CardAbilityUsed{
					CardAbilityUsed: &zb.PlayerActionCardAbilityUsed{
						Card: &zb.CardInstance{
							InstanceId: 1,
						},
						Targets: []*zb.Unit{
							&zb.Unit{
								InstanceId:       2,
								AffectObjectType: zb.AffectObjectType_CARD,
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
		response, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_OverlordSkillUsed,
				PlayerId:   "player-2",
				Action: &zb.PlayerAction_OverlordSkillUsed{
					OverlordSkillUsed: &zb.PlayerActionOverlordSkillUsed{
						SkillId: 1,
						Target: &zb.Unit{
							InstanceId:       2,
							AffectObjectType: zb.AffectObjectType_CARD,
						},
					},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("LeaveMatch", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_LeaveMatch,
				PlayerId:   "player-2",
				Action: &zb.PlayerAction_LeaveMatch{
					LeaveMatch: &zb.PlayerActionLeaveMatch{},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})
	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb.GetGameStateRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.GameState.IsEnded, "game state should be ended after use leaves the match")
	})
}

func TestGameModeOperations(t *testing.T) {
	var c *ZombieBattleground
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)

	var ID string

	t.Run("Add Game Mode", func(t *testing.T) {
		gameMode, err := c.AddGameMode(ctx, &zb.GameModeRequest{
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
		assert.Equal(t, zb.GameModeType_Community, gameMode.GameModeType)
	})

	t.Run("Get Game Mode", func(t *testing.T) {
		gameMode, err := c.GetGameMode(ctx, &zb.GetGameModeRequest{
			ID: ID,
		})
		assert.Nil(t, err)
		assert.Equal(t, "Test game mode", gameMode.Name)
		assert.Equal(t, "Just a test", gameMode.Description)
		assert.Equal(t, "0.1", gameMode.Version)
		assert.Equal(t, zb.GameModeType_Community, gameMode.GameModeType)
	})

	t.Run("Add another Game Mode", func(t *testing.T) {
		gameMode, err := c.AddGameMode(ctx, &zb.GameModeRequest{
			Name:        "Test game mode 2",
			Description: "Just another test",
			Version:     "0.1",
			Address:     "0xf16a25a1b4e6434bacf9d037d69d675dcf852692",
		})
		assert.Nil(t, err)
		assert.Equal(t, "Test game mode 2", gameMode.Name)
		assert.Equal(t, "Just another test", gameMode.Description)
		assert.Equal(t, "0.1", gameMode.Version)
		assert.Equal(t, zb.GameModeType_Community, gameMode.GameModeType)
	})

	t.Run("Update a Game Mode", func(t *testing.T) {
		gameMode, err := c.UpdateGameMode(ctx, &zb.UpdateGameModeRequest{
			ID:          ID,
			Name:        "Test game mode",
			Description: "Changed description",
			Version:     "0.2",
		})
		assert.Nil(t, err)
		assert.Equal(t, "Test game mode", gameMode.Name)
		assert.Equal(t, "Changed description", gameMode.Description)
		assert.Equal(t, "0.2", gameMode.Version)
		assert.Equal(t, zb.GameModeType_Community, gameMode.GameModeType)
	})

	t.Run("List Game Modes", func(t *testing.T) {
		gameModeList, err := c.ListGameModes(ctx, &zb.ListGameModesRequest{})
		assert.Nil(t, err)
		assert.Equal(t, 2, len(gameModeList.GameModes))
		assert.Equal(t, ID, gameModeList.GameModes[0].ID)
		assert.Equal(t, "0.2", gameModeList.GameModes[0].Version)
		assert.Equal(t, "Test game mode 2", gameModeList.GameModes[1].Name)
	})

	t.Run("Delete Game Mode", func(t *testing.T) {
		err := c.DeleteGameMode(ctx, &zb.DeleteGameModeRequest{
			ID: ID,
		})
		assert.Nil(t, err)
	})

	t.Run("GameModeList should not contain deleted GameMode", func(t *testing.T) {
		gameModeList, err := c.ListGameModes(ctx, &zb.ListGameModesRequest{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(gameModeList.GameModes))
		assert.NotEqual(t, ID, gameModeList.GameModes[0].ID)
		assert.Equal(t, "Test game mode 2", gameModeList.GameModes[0].Name)
	})
}

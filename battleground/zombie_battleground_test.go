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
			Skill: &zb.Skill{
				Title:           "Deffence",
				SkillType:       4,
				SkillTargetType: 0,
				Cost:            2,
				Value:           2,
			},
		},
		{
			HeroId:  1,
			Icon:    "asdasd",
			Name:    "Pyro Hero",
			Element: 0,
			Skill: &zb.Skill{
				Title:           "Fireball",
				SkillType:       2,
				SkillTargetType: 0,
				Cost:            2,
				Value:           1,
			},
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

const publicKeyHex = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"

func TestAccountAndDeckOperations(t *testing.T) {
	var err error

	c := &ZombieBattleground{}
	pubKey, _ := hex.DecodeString(publicKeyHex)

	addr := loom.Address{
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}

	ctx := contract.WrapPluginContext(
		plugin.CreateFakeContext(addr, addr),
	)

	err = c.Init(ctx, &initRequest)
	assert.Equal(t, err, nil)

	t.Run("AccountOperations", func(t *testing.T) {
		t.Run("CreateAccount", func(t *testing.T) {
			err = c.CreateAccount(ctx, &zb.UpsertAccountRequest{
				UserId:      "TestUser",
				Image:       "PathToImage",
				CurrentTier: 4,
			})
			assert.Equal(t, nil, err)
		})

		t.Run("UpdateAccount", func(t *testing.T) {
			account, err := c.UpdateAccount(ctx, &zb.UpsertAccountRequest{
				UserId:      "TestUser",
				Image:       "PathToImage2",
				CurrentTier: 5,
			})
			assert.Equal(t, nil, err)
			if err != nil {
				return
			}

			assert.Equal(t, int32(5), account.CurrentTier)
			assert.Equal(t, "PathToImage2", account.Image)
		})

		t.Run("GetAccount", func(t *testing.T) {
			account, err := c.GetAccount(ctx, &zb.GetAccountRequest{
				UserId: "TestUser",
			})

			assert.Equal(t, nil, err)
			if err != nil {
				return
			}

			assert.Equal(t, int32(5), account.CurrentTier)
			assert.Equal(t, "PathToImage2", account.Image)
		})
	})

	t.Run("DeckOperations", func(t *testing.T) {
		t.Run("ListDecks", func(t *testing.T) {
			deckResponse, err := c.ListDecks(ctx, &zb.ListDecksRequest{
				UserId: "TestUser",
			})

			assert.Equal(t, nil, err)
			if err != nil {
				return
			}

			assert.Equal(t, 1, len(deckResponse.Decks))
			assert.Equal(t, "Default", deckResponse.Decks[0].Name)
		})

		t.Run("GetDeck (Not Exists)", func(t *testing.T) {
			deckResponse, err := c.GetDeck(ctx, &zb.GetDeckRequest{
				UserId:   "TestUser",
				DeckName: "NotExists",
			})

			assert.Equal(t, (*zb.GetDeckResponse)(nil), deckResponse)
			assert.Equal(t, contract.ErrNotFound, err)
		})

		t.Run("GetDeck", func(t *testing.T) {
			deckResponse, err := c.GetDeck(ctx, &zb.GetDeckRequest{
				UserId:   "TestUser",
				DeckName: "Default",
			})

			assert.Equal(t, nil, err)
			assert.Equal(t, "Default", deckResponse.Deck.Name)
		})

		t.Run("AddDeck", func(t *testing.T) {
			err := c.CreateDeck(ctx, &zb.CreateDeckRequest{
				UserId: "TestUser",
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

			assert.Equal(t, nil, err)

			deckResponse, err := c.ListDecks(ctx, &zb.ListDecksRequest{
				UserId: "TestUser",
			})

			assert.Equal(t, nil, err)
			if err != nil {
				return
			}

			assert.Equal(t, 2, len(deckResponse.Decks))
		})

		t.Run("DeleteDeck", func(t *testing.T) {
			err := c.DeleteDeck(ctx, &zb.DeleteDeckRequest{
				UserId:   "TestUser",
				DeckName: "NewDeck",
			})

			assert.Equal(t, nil, err)
		})

		t.Run("DeleteDeck (Non existant)", func(t *testing.T) {
			err := c.DeleteDeck(ctx, &zb.DeleteDeckRequest{
				UserId:   "TestUser",
				DeckName: "NotExists",
			})

			assert.NotEqual(t, nil, err)
		})

		t.Run("ListCard", func(t *testing.T) {
			cardResponse, err := c.ListCardLibrary(ctx, &zb.ListCardLibraryRequest{})

			assert.Equal(t, nil, err)
			assert.Equal(t, 2, len(cardResponse.Sets[0].Cards))
		})

		t.Run("ListHero", func(t *testing.T) {
			heroResponse, err := c.ListHero(ctx, &zb.ListHeroRequest{})

			assert.Equal(t, nil, err)
			assert.Equal(t, 2, len(heroResponse.Heroes))
		})
	})
}

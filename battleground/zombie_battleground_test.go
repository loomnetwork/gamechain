package battleground

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/loomnetwork/gamechain/types/zb"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/stretchr/testify/assert"
)

var initRequest = zb.InitRequest{
	Version: "v1",
	DefaultCollection: []*zb.CardCollection{
		{CardName: "Banshee", Amount: 4},
		{CardName: "Breezee", Amount: 3},
		{CardName: "Buffer", Amount: 5},
		{CardName: "Soothsayer", Amount: 4},
		{CardName: "Wheezy", Amount: 3},
		{CardName: "Whiffer", Amount: 5},
		{CardName: "Whizpar", Amount: 4},
		{CardName: "Zhocker", Amount: 3},
		{CardName: "Bouncer", Amount: 5},
		{CardName: "Dragger", Amount: 4},
		{CardName: "Guzt", Amount: 3},
		{CardName: "Pushhh", Amount: 5},
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
		{Id: 1, Name: "Whizpar", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Air"},
		{Id: 34, Name: "Wheezy", Damage: 1, Health: 2, Kind: "CREATURE", Set: "Air"},
		{Id: 50, Name: "Soothsayer", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Air"},
		{Id: 142, Name: "Fumez", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Air"},
		{Id: 2, Name: "Pushhh", Damage: 3, Health: 3, Kind: "CREATURE", Set: "Air"},
		{Id: 3, Name: "Ztormmcaller", Damage: 3, Health: 3, Kind: "CREATURE", Set: "Air"},
		{Id: 32, Name: "Bouncer", Damage: 2, Health: 3, Kind: "CREATURE", Set: "Air"},
		{Id: 143, Name: "Gaz", Damage: 2, Health: 2, Kind: "CREATURE", Set: "Air"},
		{Id: 96, Name: "Draft", Damage: 4, Health: 5, Kind: "CREATURE", Set: "Air"},
		{Id: 97, Name: "MonZoon", Damage: 6, Health: 6, Kind: "CREATURE", Set: "Air"},
		{Id: 4, Name: "Zeuz", Damage: 5, Health: 6, Kind: "CREATURE", Set: "Air"},
		{Id: 94, Name: "Ztorm Shield", Damage: 4, Health: 4, Kind: "CREATURE", Set: "Air"},
		{Id: 5, Name: "Rockky", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Earth"},
		{Id: 7, Name: "Bolderr", Damage: 1, Health: 2, Kind: "CREATURE", Set: "Earth"},
		{Id: 98, Name: "Blocker", Damage: 0, Health: 3, Kind: "CREATURE", Set: "Earth"},
		{Id: 101, Name: "Slab", Damage: 3, Health: 4, Kind: "CREATURE", Set: "Earth"},
		{Id: 144, Name: "Pit", Damage: 0, Health: 2, Kind: "CREATURE", Set: "Earth"},
		{Id: 6, Name: "Golem", Damage: 2, Health: 6, Kind: "CREATURE", Set: "Earth"},
		{Id: 9, Name: "Walley", Damage: 2, Health: 2, Kind: "CREATURE", Set: "Earth"},
		{Id: 35, Name: "Tiny", Damage: 0, Health: 7, Kind: "CREATURE", Set: "Earth"},
		{Id: 59, Name: "Spiker", Damage: 2, Health: 3, Kind: "CREATURE", Set: "Earth"},
		{Id: 145, Name: "Crater", Damage: 1, Health: 4, Kind: "CREATURE", Set: "Earth"},
		{Id: 36, Name: "Earthshaker", Damage: 4, Health: 4, Kind: "CREATURE", Set: "Earth"},
		{Id: 61, Name: "IgneouZ", Damage: 3, Health: 3, Kind: "CREATURE", Set: "Earth"},
		{Id: 103, Name: "Pyrite", Damage: 0, Health: 8, Kind: "CREATURE", Set: "Earth"},
		{Id: 8, Name: "Mountain", Damage: 6, Health: 8, Kind: "CREATURE", Set: "Earth"},
		{Id: 104, Name: "Gaea", Damage: 4, Health: 7, Kind: "CREATURE", Set: "Earth"},
		{Id: 10, Name: "Pyromaz", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Fire"},
		{Id: 146, Name: "Quazi", Damage: 0, Health: 1, Kind: "CREATURE", Set: "Fire"},
		{Id: 11, Name: "Burrrnn", Damage: 2, Health: 2, Kind: "CREATURE", Set: "Fire"},
		{Id: 12, Name: "Cynderman", Damage: 2, Health: 3, Kind: "CREATURE", Set: "Fire"},
		{Id: 38, Name: "Werezomb", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Fire"},
		{Id: 147, Name: "Modo", Damage: 2, Health: 2, Kind: "CREATURE", Set: "Fire"},
		{Id: 14, Name: "Fire-Maw", Damage: 3, Health: 3, Kind: "CREATURE", Set: "Fire"},
		{Id: 67, Name: "Zhampion", Damage: 5, Health: 2, Kind: "CREATURE", Set: "Fire"},
		{Id: 13, Name: "Gargantua", Damage: 6, Health: 8, Kind: "CREATURE", Set: "Fire"},
		{Id: 37, Name: "Cerberus", Damage: 7, Health: 8, Kind: "CREATURE", Set: "Fire"},
		{Id: 18, Name: "Chainsaw", Damage: 0, Health: 0, Kind: "SPELL", Set: "Item"},
		{Id: 70, Name: "Goo Beaker", Damage: 0, Health: 0, Kind: "SPELL", Set: "Item"},
		{Id: 15, Name: "Stapler", Damage: 0, Health: 0, Kind: "SPELL", Set: "Item"},
		{Id: 16, Name: "Nail Bomb", Damage: 0, Health: 0, Kind: "SPELL", Set: "Item"},
		{Id: 17, Name: "Goo Bottles", Damage: 0, Health: 0, Kind: "SPELL", Set: "Item"},
		{Id: 117, Name: "Fresh Meat", Damage: 0, Health: 0, Kind: "SPELL", Set: "Item"},
		{Id: 19, Name: "Azuraz", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Life"},
		{Id: 75, Name: "Bloomer", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Life"},
		{Id: 148, Name: "Zap", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Life"},
		{Id: 20, Name: "Shroom", Damage: 4, Health: 2, Kind: "CREATURE", Set: "Life"},
		{Id: 21, Name: "Vindrom", Damage: 2, Health: 3, Kind: "CREATURE", Set: "Life"},
		{Id: 23, Name: "Puffer", Damage: 2, Health: 2, Kind: "CREATURE", Set: "Life"},
		{Id: 44, Name: "Sapper", Damage: 2, Health: 4, Kind: "CREATURE", Set: "Life"},
		{Id: 45, Name: "Keeper", Damage: 1, Health: 3, Kind: "CREATURE", Set: "Life"},
		{Id: 149, Name: "Cactuz", Damage: 2, Health: 2, Kind: "CREATURE", Set: "Life"},
		{Id: 22, Name: "Shammann", Damage: 5, Health: 6, Kind: "CREATURE", Set: "Life"},
		{Id: 42, Name: "Z-Virus", Damage: 0, Health: 0, Kind: "CREATURE", Set: "Life"},
		{Id: 43, Name: "Yggdrazil", Damage: 4, Health: 5, Kind: "CREATURE", Set: "Life"},
		{Id: 100, Name: "Zombie 1/1", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Others"},
		{Id: 101, Name: "Zombie 2/2", Damage: 2, Health: 2, Kind: "CREATURE", Set: "Others"},
		{Id: 102, Name: "Zombie Feral", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Others"},
		{Id: 155, Name: "Tainted Goo", Damage: 0, Health: 0, Kind: "SPELL", Set: "Others"},
		{Id: 156, Name: "Corrupted Goo", Damage: 0, Health: 0, Kind: "SPELL", Set: "Others"},
		{Id: 78, Name: "Rainz", Damage: 3, Health: 4, Kind: "CREATURE", Set: "Others"},
		{Id: 125, Name: "Blight", Damage: 5, Health: 4, Kind: "CREATURE", Set: "Others"},
		{Id: 131, Name: "Zteroid", Damage: 5, Health: 4, Kind: "CREATURE", Set: "Others"},
		{Id: 108, Name: "BurZt", Damage: 4, Health: 1, Kind: "CREATURE", Set: "Others"},
		{Id: 135, Name: "Vortex", Damage: 6, Health: 7, Kind: "CREATURE", Set: "Others"},
		{Id: 60, Name: "Defender", Damage: 4, Health: 6, Kind: "CREATURE", Set: "Others"},
		{Id: 24, Name: "Poizom", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Toxic"},
		{Id: 26, Name: "Hazmaz", Damage: 1, Health: 2, Kind: "CREATURE", Set: "Toxic"},
		{Id: 46, Name: "Zpitter", Damage: 2, Health: 2, Kind: "CREATURE", Set: "Toxic"},
		{Id: 150, Name: "Zeptic", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Toxic"},
		{Id: 25, Name: "Ghoul", Damage: 3, Health: 2, Kind: "CREATURE", Set: "Toxic"},
		{Id: 47, Name: "Zeeter", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Toxic"},
		{Id: 151, Name: "Hazzard", Damage: 4, Health: 4, Kind: "CREATURE", Set: "Toxic"},
		{Id: 85, Name: "Zludge", Damage: 4, Health: 4, Kind: "CREATURE", Set: "Toxic"},
		{Id: 127, Name: "Ectoplasm", Damage: 2, Health: 2, Kind: "CREATURE", Set: "Toxic"},
		{Id: 27, Name: "Cherno-bill", Damage: 7, Health: 9, Kind: "CREATURE", Set: "Toxic"},
		{Id: 132, Name: "GooZilla", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Toxic"},
		{Id: 129, Name: "Zlopper", Damage: 3, Health: 5, Kind: "CREATURE", Set: "Toxic"},
		{Id: 28, Name: "Izze", Damage: 1, Health: 1, Kind: "CREATURE", Set: "Water"},
		{Id: 49, Name: "Znowman", Damage: 0, Health: 5, Kind: "CREATURE", Set: "Water"},
		{Id: 152, Name: "Ozmoziz", Damage: 1, Health: 2, Kind: "CREATURE", Set: "Water"},
		{Id: 29, Name: "Jetter", Damage: 3, Health: 4, Kind: "CREATURE", Set: "Water"},
		{Id: 30, Name: "Freezzee", Damage: 2, Health: 2, Kind: "CREATURE", Set: "Water"},
		{Id: 153, Name: "Geyzer", Damage: 2, Health: 3, Kind: "CREATURE", Set: "Water"},
		{Id: 90, Name: "Blizzard", Damage: 3, Health: 4, Kind: "CREATURE", Set: "Water"},
		{Id: 139, Name: "Froztbite", Damage: 0, Health: 6, Kind: "CREATURE", Set: "Water"},
		{Id: 48, Name: "Zhatterer", Damage: 1, Health: 2, Kind: "CREATURE", Set: "Water"},
		{Id: 141, Name: "Maelstrom", Damage: 5, Health: 5, Kind: "CREATURE", Set: "Water"},
		{Id: 31, Name: "Tzunamy", Damage: 6, Health: 6, Kind: "CREATURE", Set: "Water"},
		{Id: 999999, Name: "Germs", Damage: 6, Health: 6, Kind: "CREATURE", Set: "Water"}, // added this card for TestDeserializeGameStateChangeActions2 test
	},
	DefaultDecks: []*zb.Deck{
		{
			Id:     0,
			HeroId: 2,
			Name:   "Default",
			Cards: []*zb.CardCollection{
				{CardName: "Azuraz", Amount: 2},
				{CardName: "Puffer", Amount: 2},
				{CardName: "Soothsayer", Amount: 2},
				{CardName: "Wheezy", Amount: 2},
				{CardName: "Whizpar", Amount: 1},
				{CardName: "Bouncer", Amount: 1},
				{CardName: "Pushhh", Amount: 1},
			},
		},
	},
}

var updateInitRequest = zb.UpdateInitRequest{
	Version: "v2",
	DefaultCollection: []*zb.CardCollection{
		{CardName: "Banshee", Amount: 4},
		{CardName: "Breezee", Amount: 3},
		{CardName: "Buffer", Amount: 5},
		{CardName: "Soothsayer", Amount: 4},
		{CardName: "Wheezy", Amount: 3},
		{CardName: "Whiffer", Amount: 5},
		{CardName: "Whizpar", Amount: 4},
		{CardName: "Zhocker", Amount: 3},
		{CardName: "Bouncer", Amount: 5},
		{CardName: "Dragger", Amount: 4},
		{CardName: "Guzt", Amount: 3},
		{CardName: "Pushhh", Amount: 5},
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
			Name:    "Azuraz",
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
				{CardName: "Banshee", Amount: 2},
				{CardName: "Breezee", Amount: 2},
				{CardName: "Buffer", Amount: 2},
				{CardName: "Soothsayer", Amount: 2},
				{CardName: "Wheezy", Amount: 2},
				{CardName: "Whiffer", Amount: 2},
				{CardName: "Whizpar", Amount: 1},
				{CardName: "Zhocker", Amount: 1},
				{CardName: "Bouncer", Amount: 1},
				{CardName: "Dragger", Amount: 1},
				{CardName: "Guzt", Amount: 1},
				{CardName: "Pushhh", Amount: 1},
			},
		},
	},
}

var updateInitRequestWithoutHeroes = zb.UpdateInitRequest{
	Version:    "v2",
	OldVersion: "v1",
	DefaultCollection: []*zb.CardCollection{
		{CardName: "Banshee", Amount: 4},
		{CardName: "Breezee", Amount: 3},
		{CardName: "Buffer", Amount: 5},
		{CardName: "Soothsayer", Amount: 4},
		{CardName: "Wheezy", Amount: 3},
		{CardName: "Whiffer", Amount: 5},
		{CardName: "Whizpar", Amount: 4},
		{CardName: "Zhocker", Amount: 3},
		{CardName: "Bouncer", Amount: 5},
		{CardName: "Dragger", Amount: 4},
		{CardName: "Guzt", Amount: 3},
		{CardName: "Pushhh", Amount: 5},
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
			Name:    "Azuraz",
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
				{CardName: "Banshee", Amount: 2},
				{CardName: "Breezee", Amount: 2},
				{CardName: "Buffer", Amount: 2},
				{CardName: "Soothsayer", Amount: 2},
				{CardName: "Wheezy", Amount: 2},
				{CardName: "Whiffer", Amount: 2},
				{CardName: "Whizpar", Amount: 1},
				{CardName: "Zhocker", Amount: 1},
				{CardName: "Bouncer", Amount: 1},
				{CardName: "Dragger", Amount: 1},
				{CardName: "Guzt", Amount: 1},
				{CardName: "Pushhh", Amount: 1},
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
	c := &ZombieBattleground{}
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
	c := &ZombieBattleground{}
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
	c := &ZombieBattleground{}
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
						CardName: "Azuraz",
					},
					{
						Amount:   1,
						CardName: "Puffer",
					},
				},
			},
			Version: "v1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, createDeckResponse)

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
						CardName: "Azuraz",
					},
					{
						Amount:   100,
						CardName: "Puffer",
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
						CardName: "Azuraz",
					},
					{
						Amount:   1,
						CardName: "Puffer",
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
						CardName: "Azuraz",
					},
					{
						Amount:   1,
						CardName: "Puffer",
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
						CardName: "Azuraz",
					},
					{
						Amount:   1,
						CardName: "Puffer",
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
		assert.NotNil(t, getDeckResponse)
		assert.Equal(t, "Edited", getDeckResponse.Deck.Name)
	})

	t.Run("EditDeck (attempt to set more number of cards)", func(t *testing.T) {
		t.Skip("Edit deck skips checking the number of cards")
		err := c.EditDeck(ctx, &zb.EditDeckRequest{
			UserId: "DeckUser",
			Deck: &zb.Deck{
				Id:     2,
				Name:   "Edited",
				HeroId: 1,
				Cards: []*zb.CardCollection{
					{
						Amount:   100,
						CardName: "Azuraz",
					},
					{
						Amount:   1,
						CardName: "Puffer",
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
						CardName: "Azuraz",
					},
					{
						Amount:   1,
						CardName: "Puffer",
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
						CardName: "Azuraz",
					},
					{
						Amount:   1,
						CardName: "Puffer",
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
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	t.Run("ListCardLibrary", func(t *testing.T) {
		cardResponse, err := c.ListCardLibrary(ctx, &zb.ListCardLibraryRequest{
			Version: "v1",
		})

		assert.Nil(t, err)
		// we expect Air, Earth, Fire, Item, Life, Others, Toxic, Water
		assert.Equal(t, 8, len(cardResponse.Sets))
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
	c := &ZombieBattleground{}
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
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	t.Run("UpdateInit", func(t *testing.T) {
		err := c.UpdateInit(ctx, &updateInitRequest)

		assert.Nil(t, err)
	})

	t.Run("UpdateInit with old card data", func(t *testing.T) {
		err := c.UpdateInit(ctx, &updateInitRequestWithoutHeroes)

		assert.Nil(t, err)
	})

	t.Run("UpdateInit with old card data but without old version (failing)", func(t *testing.T) {
		updateInitRequestWithoutHeroes.OldVersion = ""
		err := c.UpdateInit(ctx, &updateInitRequestWithoutHeroes)

		assert.NotNil(t, err)
	})

}

func TestUpdateCardListOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

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
			},
			{
				Id:      2,
				Set:     "Air",
				Name:    "Azuraz",
				Rank:    "Minion",
				Type:    "Walker",
				Damage:  1,
				Health:  1,
				Cost:    1,
				Ability: "-",
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
			},
		},
	}

	t.Run("UpdateCardList", func(t *testing.T) {
		err := c.UpdateCardList(ctx, &updateCardListRequest)
		assert.Nil(t, err)
	})
	t.Run("Check card v2", func(t *testing.T) {
		req := zb.GetCardListRequest{Version: "v2"}
		resp, err := c.GetCardList(ctx, &req)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.EqualValues(t, updateCardListRequest.Cards, resp.Cards)
	})
	t.Run("Check not exsiting version v3", func(t *testing.T) {
		req := zb.GetCardListRequest{Version: "v3"}
		_, err := c.GetCardList(ctx, &req)
		assert.NotNil(t, err)
	})
	// create deck with new card version
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "user1",
		Image:   "PathToImage",
		Version: "v1",
	}, t)
	t.Run("Create deck", func(t *testing.T) {
		_, err := c.CreateDeck(ctx, &zb.CreateDeckRequest{
			UserId: "user1",
			Deck: &zb.Deck{
				Name:   "deck1",
				HeroId: 1,
				Cards: []*zb.CardCollection{
					{
						Amount:   1,
						CardName: "Banshee",
					},
					{
						Amount:   3,
						CardName: "NewCard",
					},
				},
			},
			Version: "v2",
		})
		assert.Nil(t, err)
	})
	t.Run("Create deck with wrong card version", func(t *testing.T) {
		_, err := c.CreateDeck(ctx, &zb.CreateDeckRequest{
			UserId: "user1",
			Deck: &zb.Deck{
				Name:   "deck2",
				HeroId: 1,
				Cards: []*zb.CardCollection{
					{
						Amount:   1,
						CardName: "Azuraz",
					},
					{
						Amount:   1,
						CardName: "Puffer",
					},
				},
			},
			Version: "v2",
		})
		assert.NotNil(t, err)
	})

}

func TestUpdateHeroLibraryOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	var updateHeroLibraryRequest = zb.UpdateHeroLibraryRequest{
		Version: "v2",
		Heroes: []*zb.Hero{
			{
				HeroId:           1,
				Name:             "Hero1v2",
				ShortDescription: "Hero1v2",
			},
			{
				HeroId:           2,
				Name:             "Hero2v2",
				ShortDescription: "Hero2v2",
			},
			{
				HeroId:           3,
				Name:             "Hero3v2",
				ShortDescription: "Hero2v2",
			},
		},
	}

	t.Run("Update hero library v2", func(t *testing.T) {
		_, err := c.UpdateHeroLibrary(ctx, &updateHeroLibraryRequest)
		assert.Nil(t, err)
	})
	t.Run("Check hero library v2", func(t *testing.T) {
		req := zb.ListHeroLibraryRequest{Version: "v2"}
		resp, err := c.ListHeroLibrary(ctx, &req)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.EqualValues(t, updateHeroLibraryRequest.Heroes, resp.Heroes)
	})
	t.Run("Check not exsiting version v3", func(t *testing.T) {
		req := zb.ListHeroLibraryRequest{Version: "v3"}
		_, err := c.ListHeroLibrary(ctx, &req)
		assert.NotNil(t, err)
	})
}

func TestFindMatchOperations(t *testing.T) {
	c := &ZombieBattleground{}
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

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-1",
			Version: "v1",
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-2",
			Version: "v1",
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "two players should be matching")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "two players should be matching")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb.AcceptMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "two players should be matching")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb.AcceptMatchRequest{
			UserId:  "player-2",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "two players should be matching")
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

func TestCancelFindMatchOperations(t *testing.T) {
	c := &ZombieBattleground{}
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

	var matchID int64

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-1",
			Version: "v1",
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-2",
			Version: "v1",
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("CancelFindmatch", func(t *testing.T) {
		_, err := c.CancelFindMatch(ctx, &zb.CancelFindMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
	})

	t.Run("GetMatch", func(t *testing.T) {
		_, err := c.GetMatch(ctx, &zb.GetMatchRequest{
			MatchId: matchID,
		})
		assert.Equal(t, err, contract.ErrNotFound)
	})
}

func TestDebugFindMatchOperations(t *testing.T) {
	c := &ZombieBattleground{}
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

	var matchID int64

	t.Run("DebugFindmatch", func(t *testing.T) {
		response, err := c.DebugFindMatch(ctx, &zb.DebugFindMatchRequest{
			UserId:  "player-1",
			Version: "v1",
			Deck: &zb.Deck{
				HeroId: 2,
				Name:   "DebugDeck1",
				Cards: []*zb.CardCollection{
					{CardName: "Azuraz", Amount: 2},
					{CardName: "Puffer", Amount: 2},
					{CardName: "Soothsayer", Amount: 2},
					{CardName: "Wheezy", Amount: 2},
					{CardName: "Whizpar", Amount: 1},
					{CardName: "Bouncer", Amount: 1},
					{CardName: "Pushhh", Amount: 1},
				},
			},
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, len(response.Match.PlayerStates), "the first player should see only 1 player state")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("DebugFindmatch", func(t *testing.T) {
		response, err := c.DebugFindMatch(ctx, &zb.DebugFindMatchRequest{
			UserId:  "player-2",
			Version: "v1",
			Deck: &zb.Deck{
				HeroId: 2,
				Name:   "DebugDeck1",
				Cards: []*zb.CardCollection{
					{CardName: "Azuraz", Amount: 2},
					{CardName: "Puffer", Amount: 2},
					{CardName: "Soothsayer", Amount: 2},
					{CardName: "Wheezy", Amount: 2},
					{CardName: "Whizpar", Amount: 1},
					{CardName: "Bouncer", Amount: 1},
					{CardName: "Pushhh", Amount: 1},
				},
			},
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
}

func TestFindMatchWithTagOperations(t *testing.T) {
	c := &ZombieBattleground{}
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
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-3",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-1-tag",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-2-tag",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-3-tag",
		Version: "v1",
	}, t)

	var matchID, matchIDTag int64

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-1",
			Version: "v1",
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-2",
			Version: "v1",
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb.AcceptMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb.AcceptMatchRequest{
			UserId:  "player-2",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	tags := []string{"tag1"}

	t.Run("RegisterPlayerPoolTag", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-1-tag",
			Version: "v1",
			Tags:    tags,
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPoolTag", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-2-tag",
			Version: "v1",
			Tags:    tags,
		})
		assert.Nil(t, err)
	})

	t.Run("FindmatchTag", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-1-tag",
			Tags:   tags,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchIDTag = response.Match.Id
	})

	t.Run("FindmatchTag", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-2-tag",
			Tags:   tags,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchIDTag, response.Match.Id)
	})

	t.Run("AcceptMatchTag", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb.AcceptMatchRequest{
			UserId:  "player-1-tag",
			MatchId: matchIDTag,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchIDTag, response.Match.Id)
		assert.NotEqual(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatchTag", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb.AcceptMatchRequest{
			UserId:  "player-2-tag",
			MatchId: matchIDTag,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchIDTag, response.Match.Id)
		assert.NotEqual(t, matchID, response.Match.Id)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the second player should 2 player states")
		assert.Equal(t, zb.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	// check tag and non-tag players don't get matched
	tags = []string{"tag3"}

	t.Run("RegisterPlayerPoolTag", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-3",
			Version: "v1",
			Tags:    tags,
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPoolTag", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-3-tag",
			Version: "v1",
			Tags:    tags,
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		_, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-3",
		})
		assert.NotNil(t, err)
	})

	t.Run("FindmatchTag", func(t *testing.T) {
		_, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-3-tag",
		})
		assert.NotNil(t, err)
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
		setupAccount(c, ctx, &zb.UpsertAccountRequest{
			UserId:  fmt.Sprintf("player-%d", i+1),
			Version: "v1",
		}, t)
	}

	for i := 0; i < numPlayers; i++ {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  fmt.Sprintf("player-%d", i+1),
			Version: "v1",
		})
		assert.Nil(t, err)
	}

	for i := 0; i < numPlayers; i++ {
		func(i int) {
			response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
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

	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-3",
		Version: "v1",
	}, t)

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-1",
			Version: "v1",
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-2",
			Version: "v1",
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-3",
			Version: "v1",
		})
		assert.Nil(t, err)
	})

	var matchID int64
	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	// move time forward to expire the matchmaking
	fc.SetTime(now.Add(2 * MMTimeout))

	t.Run("FindMatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.Equal(t, zb.Match_Timedout, response.Match.Status)
	})

	t.Run("GetMatch", func(t *testing.T) {
		_, err := c.GetMatch(ctx, &zb.GetMatchRequest{
			MatchId: matchID,
		})
		assert.NotNil(t, err)
	})
}

func TestGameStateOperations(t *testing.T) {
	c := &ZombieBattleground{}
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

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-1",
			Version: "v1",
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-2",
			Version: "v1",
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb.AcceptMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb.AcceptMatchRequest{
			UserId:  "player-2",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
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
							InstanceId: 8,
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
	t.Run("SendRankBuff", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_RankBuff,
				PlayerId:   "player-1",
				Action: &zb.PlayerAction_RankBuff{
					RankBuff: &zb.PlayerActionRankBuff{
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
							InstanceId: 13,
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
							InstanceId: 13,
						},
						AffectObjectType: zb.AffectObjectType_CHARACTER,
						Target: &zb.Unit{
							InstanceId: 8,
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
	c := &ZombieBattleground{}
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

func TestCardPlayOperations(t *testing.T) {
	c := &ZombieBattleground{}
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

	var matchID int64

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-1",
			Version: "v1",
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:  1,
			UserId:  "player-2",
			Version: "v1",
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("Acceptmatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb.AcceptMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the second player should 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'Matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("Acceptmatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb.AcceptMatchRequest{
			UserId:  "player-2",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the second player should 2 player states")
		assert.Equal(t, zb.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
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
							InstanceId: 8,
						},
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

	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)

	var matchID int64

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:     1,
			UserId:     "player-1",
			Version:    "v1",
			RandomSeed: 2,
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:     1,
			UserId:     "player-2",
			Version:    "v1",
			RandomSeed: 2,
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb.AcceptMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb.AcceptMatchRequest{
			UserId:  "player-2",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb.GetGameStateRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
		assert.EqualValues(t, 0, response.GameState.CurrentPlayerIndex)
	})

	t.Run("SendEndturnPlayer1_Success", func(t *testing.T) {
		_, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: 1,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_EndTurn,
				PlayerId:   "player-1",
			},
		})
		assert.Nil(t, err)
	})
	t.Run("SendEndturnPlayer2_Success", func(t *testing.T) {
		_, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: 1,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_EndTurn,
				PlayerId:   "player-2",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb.GetGameStateRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
		assert.EqualValues(t, 0, response.GameState.CurrentPlayerIndex)
	})

	// player1 don't sent enturn within TurnTimeout
	// move time forward to expire the player's turn
	fc.SetTime(now.Add(TurnTimeout + (time.Second * 10)))

	t.Run("Player2_CheckStatus", func(t *testing.T) {
		_, err := c.CheckGameStatus(ctx, &zb.CheckGameStatusRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
	})

	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb.GetGameStateRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.GameState.IsEnded, "game state should be ended after use leaves the match")
		latestAction := response.GameState.PlayerActions[len(response.GameState.PlayerActions)-1]
		assert.Equal(t, zb.PlayerActionType_LeaveMatch, latestAction.ActionType)
		assert.Equal(t, "player-2", response.GameState.Winner)
	})

	t.Run("SendAnyActionShould_Failed", func(t *testing.T) {
		_, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: 1,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_EndTurn,
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

	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-1",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-2",
		Version: "v1",
	}, t)

	var matchID int64

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:     1,
			UserId:     "player-1",
			Version:    "v1",
			RandomSeed: 2,
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			DeckId:     1,
			UserId:     "player-2",
			Version:    "v1",
			RandomSeed: 2,
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-1",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		matchID = response.Match.Id
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb.AcceptMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("AcceptMatch", func(t *testing.T) {
		response, err := c.AcceptMatch(ctx, &zb.AcceptMatchRequest{
			UserId:  "player-2",
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Started, response.Match.Status, "match status should be 'started'")
		assert.Equal(t, matchID, response.Match.Id)
	})

	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb.GetGameStateRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
		assert.EqualValues(t, 0, response.GameState.CurrentPlayerIndex)
	})

	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb.GetGameStateRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
		assert.EqualValues(t, 0, response.GameState.CurrentPlayerIndex)
	})

	// player1 don't sent enturn within TurnTimeout
	// move time forward to expire the player's turn
	fc.SetTime(now.Add(TurnTimeout + (time.Second * 10)))

	t.Run("Player2_CheckStatus", func(t *testing.T) {
		_, err := c.CheckGameStatus(ctx, &zb.CheckGameStatusRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
	})

	t.Run("GetGameState", func(t *testing.T) {
		response, err := c.GetGameState(ctx, &zb.GetGameStateRequest{
			MatchId: 1,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.GameState.IsEnded, "game state should be ended after use leaves the match")
		latestAction := response.GameState.PlayerActions[len(response.GameState.PlayerActions)-1]
		assert.Equal(t, zb.PlayerActionType_LeaveMatch, latestAction.ActionType)
		assert.Equal(t, "player-2", response.GameState.Winner)
	})
}

func TestAIDeckOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	aiDecks := []*zb.Deck{
		{
			Id:     1,
			HeroId: 2,
			Name:   "AI Decks",
			Cards: []*zb.CardCollection{
				{CardName: "Banshee", Amount: 2},
				{CardName: "Breezee", Amount: 2},
				{CardName: "Buffer", Amount: 2},
				{CardName: "Soothsayer", Amount: 2},
				{CardName: "Wheezy", Amount: 2},
				{CardName: "Whiffer", Amount: 2},
				{CardName: "Whizpar", Amount: 1},
				{CardName: "Zhocker", Amount: 1},
				{CardName: "Bouncer", Amount: 1},
				{CardName: "Dragger", Amount: 1},
				{CardName: "Guzt", Amount: 1},
				{CardName: "Pushhh", Amount: 1},
			},
		},
	}

	t.Run("Set AI Decks", func(t *testing.T) {
		req := &zb.SetAIDecksRequest{
			Decks:   aiDecks,
			Version: "v1",
		}
		err := c.SetAIDecks(ctx, req)
		assert.Nil(t, err)
	})

	t.Run("Get AI Decks", func(t *testing.T) {
		_, err := c.GetAIDecks(ctx, &zb.GetAIDecksRequest{
			Version: "v1",
		})
		assert.Nil(t, err)
	})
}

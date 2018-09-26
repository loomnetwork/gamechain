package main

import (
	"encoding/hex"
	"testing"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/zombie_battleground/battleground"
	"github.com/loomnetwork/zombie_battleground/types/zb"
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

func setup(c *battleground.ZombieBattleground, pubKeyHex string, addr *loom.Address, ctx *contract.Context, t *testing.T) {

	c = &battleground.ZombieBattleground{}
	pubKey, _ := hex.DecodeString(pubKeyHex)

	addr = &loom.Address{
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}

	*ctx = contract.WrapPluginContext(
		plugin.CreateFakeContext(*addr, *addr),
	)

	err := c.Init(*ctx, &battleground.initRequest)
	assert.Nil(t, err)
}

func setupAccount(c *battleground.ZombieBattleground, ctx contract.Context, upsertAccountRequest *zb.UpsertAccountRequest, t *testing.T) {
	err := c.CreateAccount(ctx, upsertAccountRequest)
	assert.Nil(t, err)
}
func setupZBContract() {

	var c *battleground.ZombieBattleground
	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "AccountUser",
		Image:   "PathToImage",
		Version: "v1",
	}, t)

}

func main() {

	runGocui()
	return
}

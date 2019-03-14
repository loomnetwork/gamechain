package main

import (
	"encoding/hex"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/gamechain/battleground"
	"github.com/loomnetwork/gamechain/types/zb"
)

var initRequest = zb.InitRequest{
	Version: "v1",
	DefaultCollection: []*zb.CardCollectionCard{
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
				Title: "Attack",
				Skill: zb.OverlordSkillKind_IceBolt,
				SkillTargets: []zb.OverlordAbilityTarget_Enum{
					zb.OverlordAbilityTarget_AllCards,
					zb.OverlordAbilityTarget_PlayerCard,
				},
				Value: 1,
			}},
		},
		{
			HeroId:     1,
			Experience: 0,
			Level:      2,
			Skills: []*zb.Skill{{
				Title: "Deffence",
				Skill: zb.OverlordSkillKind_Blizzard,
				SkillTargets: []zb.OverlordAbilityTarget_Enum{
					zb.OverlordAbilityTarget_Player,
					zb.OverlordAbilityTarget_OpponentCard,
				},
				Value: 2,
			}},
		},
	},
	Cards: []*zb.Card{
		{
			MouldId: 1,
			Faction: zb.Faction_Air,
			Name:    "Soothsayer",
			Rank:    zb.CreatureRank_Minion,
			Type:    zb.CreatureType_Walker,
			Damage:  2,
			Defense: 1,
			GooCost: 2,
			Abilities: []*zb.CardAbility{
				{
					Type:         zb.CardAbilityType_DrawCard,
					ActivityType: zb.CardAbilityActivityType_Passive,
					Trigger:      zb.CardAbilityTrigger_Entry,
					Faction:      zb.Faction_None,
				},
			},
			PictureTransform: &zb.PictureTransform{
				Position: &zb.Vector3Float{
					X: 1.5,
					Y: 2.5,
					Z: 3.5,
				},
				Scale: &zb.Vector3Float{
					X: 0.5,
					Y: 0.5,
					Z: 0.5,
				},
			},
		},
		{
			MouldId: 2,
			Faction: zb.Faction_Air,
			Name:    "Azuraz",
			Rank:    zb.CreatureRank_Minion,
			Type:    zb.CreatureType_Walker,
			Damage:  1,
			Defense: 1,
			GooCost: 1,
			Abilities: []*zb.CardAbility{
				{
					Type:         zb.CardAbilityType_ModificatorStats,
					ActivityType: zb.CardAbilityActivityType_Passive,
					Trigger:      zb.CardAbilityTrigger_Permanent,
					TargetTypes: []zb.CardAbilityTarget_Enum{
						zb.CardAbilityTarget_None,
					},
					Stat:    zb.StatType_Damage,
					Faction: zb.Faction_Earth,
					Value:   1,
				},
			},
		},
	},
	DefaultDecks: []*zb.Deck{
		{
			Id:     0,
			HeroId: 2,
			Name:   "Default",
			Cards: []*zb.DeckCard{
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

func setup(c *battleground.ZombieBattleground, pubKeyHex string, addr *loom.Address, ctx *contract.Context) {

	pubKey, _ := hex.DecodeString(pubKeyHex)

	addr = &loom.Address{
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}

	*ctx = contract.WrapPluginContext(
		plugin.CreateFakeContext(*addr, *addr),
	)

	err := c.Init(*ctx, &initRequest)
	if err != nil {
		panic(err)
	}
}

func setupAccount(c *battleground.ZombieBattleground, ctx contract.Context, upsertAccountRequest *zb.UpsertAccountRequest) {
	err := c.CreateAccount(ctx, upsertAccountRequest)
	if err != nil {
		panic(err)
	}
}
func setupZBContract() {

	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address

	setup(zvContract, pubKeyHexString, &addr, &ctx)
	setupAccount(zvContract, ctx, &zb.UpsertAccountRequest{
		UserId:  "AccountUser",
		Image:   "PathToImage",
		Version: "v1",
	})

}
func listItemsForPlayer(playerId int) []string {
	res := []string{}

	cardCollection, err := zvContract.GetCollection(ctx, &zb.GetCollectionRequest{
		UserId: "AccountUser",
	})
	if err != nil {
		panic(err)
	}
	for _, v := range cardCollection.Cards {
		res = append(res, v.CardName)
	}

	return res
}

var zvContract *battleground.ZombieBattleground
var ctx contract.Context

func main() {
	zvContract = &battleground.ZombieBattleground{}
	setupZBContract()

	runGocui()
	return
}

package main

import (
	"encoding/hex"
	"fmt"
	"github.com/loomnetwork/gamechain/battleground"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/gamechain/types/zb/zb_calls"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)

/*var initRequest = zb_calls.InitRequest{
	Version: "v1",
	DefaultCollection: []*zb_data.CardCollectionCard{
		{
			MouldId: 90,
			Amount:   4,
		},
		{
			MouldId: 91,
			Amount:   3,
		},
		{
			MouldId: 96,
			Amount:   5,
		},
		{
			MouldId: 3,
			Amount:   4,
		},
		{
			MouldId: 2,
			Amount:   3,
		},
		{
			MouldId: 92,
			Amount:   5,
		},
		{
			MouldId: 1,
			Amount:   4,
		},
		{
			MouldId: 93,
			Amount:   3,
		},
		{
			MouldId: 7,
			Amount:   5,
		},
		{
			MouldId: 94,
			Amount:   4,
		},
		{
			MouldId: 95,
			Amount:   3,
		},
		{
			MouldId: 5,
			Amount:   5,
		},
	},
	Overlords: []*zb_data.Overlord{
		{
			OverlordId: 0,
			Experience: 0,
			Level:      1,
			Skills: []*zb_data.Skill{{
				Title: "Attack",
				Skill: zb_enums.OverlordSkill_IceBolt,
				SkillTargets: []zb_enums.SkillTarget_Enum{
					zb_enums.SkillTarget_AllCards,
					zb_enums.SkillTarget_PlayerCard,
				},
				Value: 1,
			}},
		},
		{
			OverlordId: 1,
			Experience: 0,
			Level:      2,
			Skills: []*zb_data.Skill{{
				Title: "Deffence",
				Skill: zb_enums.OverlordSkill_Blizzard,
				SkillTargets: []zb_enums.SkillTarget_Enum{
					zb_enums.SkillTarget_Player,
					zb_enums.SkillTarget_OpponentCard,
				},
				Value: 2,
			}},
		},
	},
	Cards: []*zb_data.Card{
		{
			MouldId: 1,
			Faction: zb_enums.Faction_Air,
			Name:    "Soothsayer",
			Rank:    zb_enums.CreatureRank_Minion,
			Type:    zb_enums.CardType_Walker,
			Damage:  2,
			Defense: 1,
			Cost: 2,
			Abilities: []*zb_data.AbilityData{
				{
					Ability:  zb_enums.AbilityType_DrawCard,
					Activity: zb_enums.AbilityActivity_Passive,
					Trigger:  zb_enums.AbilityTrigger_Entry,
					Faction:  zb_enums.Faction_None,
				},
			},
			PictureTransform: &zb_data.PictureTransform{
				Position: &zb_data.Vector3Float{
					X: 1.5,
					Y: 2.5,
					Z: 3.5,
				},
				Scale: &zb_data.Vector3Float{
					X: 0.5,
					Y: 0.5,
					Z: 0.5,
				},
			},
		},
		{
			MouldId: 2,
			Faction: zb_enums.Faction_Air,
			Name:    "Azuraz",
			Rank:    zb_enums.CreatureRank_Minion,
			Type:    zb_enums.CardType_Walker,
			Damage:  1,
			Defense: 1,
			Cost: 1,
			Abilities: []*zb_data.AbilityData{
				{
					Ability:  zb_enums.AbilityType_ModificatorStats,
					Activity: zb_enums.AbilityActivity_Passive,
					Trigger:  zb_enums.AbilityTrigger_Permanent,
					Targets: []zb.Target_Enum{
						zb.Target_None,
					},
					Stat:    zb_enums.Stat_Damage,
					Faction: zb_enums.Faction_Earth,
					Value:   1,
				},
			},
		},
	},
	DefaultDecks: []*zb_data.Deck{
		{
			Id:         0,
			OverlordId: 2,
			Name:       "Default",
			Cards: []*zb_data.DeckCard{
				{
					MouldId: 90,
					Amount:   2,
				},
				{
					MouldId: 91,
					Amount:   2,
				},
				{
					MouldId: 96,
					Amount:   2,
				},
				{
					MouldId: 3,
					Amount:   2,
				},
				{
					MouldId: 2,
					Amount:   2,
				},
				{
					MouldId: 92,
					Amount:   2,
				},
				{
					MouldId: 1,
					Amount:   1,
				},
				{
					MouldId: 93,
					Amount:   1,
				},
				{
					MouldId: 7,
					Amount:   1,
				},
				{
					MouldId: 94,
					Amount:   1,
				},
				{
					MouldId: 95,
					Amount:   1,
				},
				{
					MouldId: 5,
					Amount:   1,
				},
			},
		},
	},
}*/

func setup(c *battleground.ZombieBattleground, pubKeyHex string, addr *loom.Address, ctx *contract.Context) {

	pubKey, _ := hex.DecodeString(pubKeyHex)

	addr = &loom.Address{
		Local: loom.LocalAddressFromPublicKey(pubKey),
	}

	*ctx = contract.WrapPluginContext(
		plugin.CreateFakeContext(*addr, *addr),
	)

	// FIXME
	/*err := c.Init(*ctx, &initRequest)
	if err != nil {
		panic(err)
	}*/
}

func setupAccount(c *battleground.ZombieBattleground, ctx contract.Context, upsertAccountRequest *zb_calls.UpsertAccountRequest) {
	err := c.CreateAccount(ctx, upsertAccountRequest)
	if err != nil {
		panic(err)
	}
}
func setupZBContract() {

	var pubKeyHexString = "e4008e26428a9bca87465e8de3a8d0e9c37a56ca619d3d6202b0567528786618"
	var addr loom.Address

	setup(zvContract, pubKeyHexString, &addr, &ctx)
	setupAccount(zvContract, ctx, &zb_calls.UpsertAccountRequest{
		UserId:  "AccountUser",
		Image:   "PathToImage",
		Version: "v1",
	})

}
func listItemsForPlayer(playerId int) []string {
	res := []string{}

	cardCollection, err := zvContract.GetCollection(ctx, &zb_calls.GetCollectionRequest{
		UserId: "AccountUser",
	})
	if err != nil {
		panic(err)
	}
	for _, v := range cardCollection.Cards {
		res = append(res, fmt.Sprintf("Mould Id %d", v.MouldId))
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

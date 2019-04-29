package battleground

import (
	"encoding/hex"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
	latypes "github.com/loomnetwork/loomauth/types"
	"github.com/stretchr/testify/assert"
)

var initRequest = zb.InitRequest{
	Version: "v1",
	DefaultCollection: []*zb.CardCollectionCard{
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
	Overlords: []*zb.Overlord{
		{
			OverlordId: 0,
			Experience: 0,
			Level:      1,
			Skills: []*zb.Skill{{
				Title: "Attack",
				Skill: zb.OverlordSkill_IceBolt,
				SkillTargets: []zb.SkillTarget_Enum{
					zb.SkillTarget_AllCards,
					zb.SkillTarget_PlayerCard,
				},
				Value: 1,
			}},
		},
		{
			OverlordId: 1,
			Experience: 0,
			Level:      2,
			Skills: []*zb.Skill{{
				Title: "Deffence",
				Skill: zb.OverlordSkill_Blizzard,
				SkillTargets: []zb.SkillTarget_Enum{
					zb.SkillTarget_Player,
					zb.SkillTarget_OpponentCard,
				},
				Value: 2,
			}},
		},
	},
	Cards: []*zb.Card{
		{MouldId: 1, Name: "Whizpar", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Air, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 34, Name: "Wheezy", Damage: 1, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Air, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 50, Name: "Soothsayer", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Air, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 142, Name: "Fumez", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Air, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 2, Name: "Pushhh", Damage: 3, Defense: 3, Kind: zb.CardKind_Creature, Faction: zb.Faction_Air, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 3, Name: "Ztormmcaller", Damage: 3, Defense: 3, Kind: zb.CardKind_Creature, Faction: zb.Faction_Air, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 32, Name: "Bouncer", Damage: 2, Defense: 3, Kind: zb.CardKind_Creature, Faction: zb.Faction_Air, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 143, Name: "Gaz", Damage: 2, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Air, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 96, Name: "Draft", Damage: 4, Defense: 5, Kind: zb.CardKind_Creature, Faction: zb.Faction_Air, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 97, Name: "MonZoon", Damage: 6, Defense: 6, Kind: zb.CardKind_Creature, Faction: zb.Faction_Air, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 4, Name: "Zeuz", Damage: 5, Defense: 6, Kind: zb.CardKind_Creature, Faction: zb.Faction_Air, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 94, Name: "Ztorm Shield", Damage: 4, Defense: 4, Kind: zb.CardKind_Creature, Faction: zb.Faction_Air, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 5, Name: "Rockky", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Earth, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 7, Name: "Bolderr", Damage: 1, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Earth, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 98, Name: "Blocker", Damage: 0, Defense: 3, Kind: zb.CardKind_Creature, Faction: zb.Faction_Earth, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 101, Name: "Slab", Damage: 3, Defense: 4, Kind: zb.CardKind_Creature, Faction: zb.Faction_Earth, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 144, Name: "Pit", Damage: 0, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Earth, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 6, Name: "Golem", Damage: 2, Defense: 6, Kind: zb.CardKind_Creature, Faction: zb.Faction_Earth, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 9, Name: "Walley", Damage: 2, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Earth, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 35, Name: "Tiny", Damage: 0, Defense: 7, Kind: zb.CardKind_Creature, Faction: zb.Faction_Earth, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 59, Name: "Spiker", Damage: 2, Defense: 3, Kind: zb.CardKind_Creature, Faction: zb.Faction_Earth, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 145, Name: "Crater", Damage: 1, Defense: 4, Kind: zb.CardKind_Creature, Faction: zb.Faction_Earth, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 36, Name: "Earthshaker", Damage: 4, Defense: 4, Kind: zb.CardKind_Creature, Faction: zb.Faction_Earth, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 61, Name: "IgneouZ", Damage: 3, Defense: 3, Kind: zb.CardKind_Creature, Faction: zb.Faction_Earth, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 103, Name: "Pyrite", Damage: 0, Defense: 8, Kind: zb.CardKind_Creature, Faction: zb.Faction_Earth, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 8, Name: "Mountain", Damage: 6, Defense: 8, Kind: zb.CardKind_Creature, Faction: zb.Faction_Earth, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 104, Name: "Gaea", Damage: 4, Defense: 7, Kind: zb.CardKind_Creature, Faction: zb.Faction_Earth, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 10, Name: "Pyromaz", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Fire, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 146, Name: "Quazi", Damage: 0, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Fire, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 11, Name: "Burrrnn", Damage: 2, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Fire, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 12, Name: "Cynderman", Damage: 2, Defense: 3, Kind: zb.CardKind_Creature, Faction: zb.Faction_Fire, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 38, Name: "Werezomb", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Fire, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 147, Name: "Modo", Damage: 2, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Fire, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 14, Name: "Fire-Maw", Damage: 3, Defense: 3, Kind: zb.CardKind_Creature, Faction: zb.Faction_Fire, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 67, Name: "Zhampion", Damage: 5, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Fire, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 13, Name: "Gargantua", Damage: 6, Defense: 8, Kind: zb.CardKind_Creature, Faction: zb.Faction_Fire, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 37, Name: "Cerberus", Damage: 7, Defense: 8, Kind: zb.CardKind_Creature, Faction: zb.Faction_Fire, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 18, Name: "Chainsaw", Damage: 0, Defense: 0, Kind: zb.CardKind_Spell, Faction: zb.Faction_Item, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 70, Name: "Goo Beaker", Damage: 0, Defense: 0, Kind: zb.CardKind_Spell, Faction: zb.Faction_Item, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 15, Name: "Stapler", Damage: 0, Defense: 0, Kind: zb.CardKind_Spell, Faction: zb.Faction_Item, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 16, Name: "Nail Bomb", Damage: 0, Defense: 0, Kind: zb.CardKind_Spell, Faction: zb.Faction_Item, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 17, Name: "Goo Bottles", Damage: 0, Defense: 0, Kind: zb.CardKind_Spell, Faction: zb.Faction_Item, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 117, Name: "Fresh Meat", Damage: 0, Defense: 0, Kind: zb.CardKind_Spell, Faction: zb.Faction_Item, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 19, Name: "Azuraz", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Life, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 75, Name: "Bloomer", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Life, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 148, Name: "Zap", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Life, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 20, Name: "Shroom", Damage: 4, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Life, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 21, Name: "Vindrom", Damage: 2, Defense: 3, Kind: zb.CardKind_Creature, Faction: zb.Faction_Life, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 23, Name: "Puffer", Damage: 2, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Life, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 44, Name: "Sapper", Damage: 2, Defense: 4, Kind: zb.CardKind_Creature, Faction: zb.Faction_Life, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 45, Name: "Keeper", Damage: 1, Defense: 3, Kind: zb.CardKind_Creature, Faction: zb.Faction_Life, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 149, Name: "Cactuz", Damage: 2, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Life, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 22, Name: "Shammann", Damage: 5, Defense: 6, Kind: zb.CardKind_Creature, Faction: zb.Faction_Life, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 42, Name: "Z-Virus", Damage: 0, Defense: 0, Kind: zb.CardKind_Creature, Faction: zb.Faction_Life, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 43, Name: "Yggdrazil", Damage: 4, Defense: 5, Kind: zb.CardKind_Creature, Faction: zb.Faction_Life, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 100, Name: "Zombie 1/1", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Others, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 101, Name: "Zombie 2/2", Damage: 2, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Others, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 102, Name: "Zombie Feral", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Others, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 155, Name: "Tainted Goo", Damage: 0, Defense: 0, Kind: zb.CardKind_Spell, Faction: zb.Faction_Others, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 156, Name: "Corrupted Goo", Damage: 0, Defense: 0, Kind: zb.CardKind_Spell, Faction: zb.Faction_Others, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 78, Name: "Rainz", Damage: 3, Defense: 4, Kind: zb.CardKind_Creature, Faction: zb.Faction_Others, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 125, Name: "Blight", Damage: 5, Defense: 4, Kind: zb.CardKind_Creature, Faction: zb.Faction_Others, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 131, Name: "Zteroid", Damage: 5, Defense: 4, Kind: zb.CardKind_Creature, Faction: zb.Faction_Others, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 108, Name: "BurZt", Damage: 4, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Others, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 135, Name: "Vortex", Damage: 6, Defense: 7, Kind: zb.CardKind_Creature, Faction: zb.Faction_Others, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 60, Name: "Defender", Damage: 4, Defense: 6, Kind: zb.CardKind_Creature, Faction: zb.Faction_Others, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 24, Name: "Poizom", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Toxic, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 26, Name: "Hazmaz", Damage: 1, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Toxic, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 46, Name: "Zpitter", Damage: 2, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Toxic, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 150, Name: "Zeptic", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Toxic, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 25, Name: "Ghoul", Damage: 3, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Toxic, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 47, Name: "Zeeter", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Toxic, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 151, Name: "Hazzard", Damage: 4, Defense: 4, Kind: zb.CardKind_Creature, Faction: zb.Faction_Toxic, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 85, Name: "Zludge", Damage: 4, Defense: 4, Kind: zb.CardKind_Creature, Faction: zb.Faction_Toxic, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 127, Name: "Ectoplasm", Damage: 2, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Toxic, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 27, Name: "Cherno-bill", Damage: 7, Defense: 9, Kind: zb.CardKind_Creature, Faction: zb.Faction_Toxic, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 132, Name: "GooZilla", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Toxic, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 129, Name: "Zlopper", Damage: 3, Defense: 5, Kind: zb.CardKind_Creature, Faction: zb.Faction_Toxic, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 28, Name: "Izze", Damage: 1, Defense: 1, Kind: zb.CardKind_Creature, Faction: zb.Faction_Water, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 49, Name: "Znowman", Damage: 0, Defense: 5, Kind: zb.CardKind_Creature, Faction: zb.Faction_Water, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 152, Name: "Ozmoziz", Damage: 1, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Water, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 29, Name: "Jetter", Damage: 3, Defense: 4, Kind: zb.CardKind_Creature, Faction: zb.Faction_Water, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 30, Name: "Freezzee", Damage: 2, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Water, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 153, Name: "Geyzer", Damage: 2, Defense: 3, Kind: zb.CardKind_Creature, Faction: zb.Faction_Water, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 90, Name: "Blizzard", Damage: 3, Defense: 4, Kind: zb.CardKind_Creature, Faction: zb.Faction_Water, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 139, Name: "Froztbite", Damage: 0, Defense: 6, Kind: zb.CardKind_Creature, Faction: zb.Faction_Water, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 48, Name: "Zhatterer", Damage: 1, Defense: 2, Kind: zb.CardKind_Creature, Faction: zb.Faction_Water, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 141, Name: "Maelstrom", Damage: 5, Defense: 5, Kind: zb.CardKind_Creature, Faction: zb.Faction_Water, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 31, Name: "Tzunamy", Damage: 6, Defense: 6, Kind: zb.CardKind_Creature, Faction: zb.Faction_Water, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}},
		{MouldId: 999999, Name: "Germs", Damage: 6, Defense: 6, Kind: zb.CardKind_Creature, Faction: zb.Faction_Water, PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}}}, // added this card for TestDeserializeGameStateChangeActions2 test
	},
	DefaultDecks: []*zb.Deck{
		{
			Id:         0,
			OverlordId: 2,
			Name:       "Default",
			Cards: []*zb.DeckCard{
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
	AiDecks: []*zb.AIDeck{
		{
			Deck: &zb.Deck{
				Id:         1,
				OverlordId: 2,
				Name:       "AI Decks",
				Cards: []*zb.DeckCard{
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
			Type: zb.AIType_MIXED_AI,
		},
	},
}

var updateInitRequest = zb.UpdateInitRequest{
	InitData: &zb.InitData{
		Version: "v2",
		DefaultCollection: []*zb.CardCollectionCard{
			{CardName: "Soothsayer", Amount:15},
			{CardName: "Azuraz", Amount: 15},
		},
		Overlords: []*zb.Overlord{
			{
				OverlordId: 0,
				Experience: 0,
				Level:      1,
				Skills: []*zb.Skill{{
					Title: "Attack",
					Skill: zb.OverlordSkill_IceBolt,
					SkillTargets: []zb.SkillTarget_Enum{
						zb.SkillTarget_AllCards,
						zb.SkillTarget_PlayerCard,
					},
					Value: 1,
				}},
			},
			{
				OverlordId: 1,
				Experience: 0,
				Level:      2,
				Skills: []*zb.Skill{{
					Title: "Deffence",
					Skill: zb.OverlordSkill_Blizzard,
					SkillTargets: []zb.SkillTarget_Enum{
						zb.SkillTarget_Player,
						zb.SkillTarget_OpponentCard,
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
				Type:    zb.CardType_Walker,
				Damage:  2,
				Defense: 1,
				Cost: 2,
				Abilities: []*zb.AbilityData{
					{
						Ability:  zb.AbilityType_DrawCard,
						Activity: zb.AbilityActivity_Passive,
						Trigger:  zb.AbilityTrigger_Entry,
						Faction:  zb.Faction_None,
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
				MouldId:          2,
				Faction:          zb.Faction_Air,
				Name:             "Azuraz",
				Rank:             zb.CreatureRank_Minion,
				Type:             zb.CardType_Walker,
				Damage:           1,
				Defense:          1,
				Cost:          1,
				PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}},
				Abilities: []*zb.AbilityData{
					{
						Ability:  zb.AbilityType_ModificatorStats,
						Activity: zb.AbilityActivity_Passive,
						Trigger:  zb.AbilityTrigger_Permanent,
						Targets: []zb.Target_Enum{
							zb.Target_None,
						},
						Stat:    zb.Stat_Damage,
						Faction: zb.Faction_Earth,
						Value:   1,
					},
				},
			},
			{
				MouldId:          3,
				Faction:          zb.Faction_Air,
				Name:             "NewCard",
				Rank:             zb.CreatureRank_Minion,
				Type:             zb.CardType_Walker,
				Damage:           1,
				Defense:          1,
				Cost:          1,
				PictureTransform: &zb.PictureTransform{Position: &zb.Vector3Float{X: 1, Y: 1, Z: 1}, Scale: &zb.Vector3Float{X: 1, Y: 1, Z: 1}},
				Abilities: []*zb.AbilityData{
					{
						Ability:  zb.AbilityType_ModificatorStats,
						Activity: zb.AbilityActivity_Passive,
						Trigger:  zb.AbilityTrigger_Permanent,
						Targets: []zb.Target_Enum{
							zb.Target_None,
						},
						Stat:    zb.Stat_Damage,
						Faction: zb.Faction_Water,
						Value:   1,
					},
				},
			},
		},
		DefaultDecks: []*zb.Deck{
			{
				Id:         0,
				OverlordId: 2,
				Name:       "Default",
				Cards: []*zb.DeckCard{
					{CardName: "Soothsayer", Amount:15},
					{CardName: "Azuraz", Amount: 15},
				},
			},
		},
		AiDecks: []*zb.AIDeck{
			{
				Deck: &zb.Deck{
					Id:         1,
					OverlordId: 2,
					Name:       "AI Decks",
					Cards: []*zb.DeckCard{
						{CardName: "Soothsayer", Amount: 2},
						{CardName: "Azuraz", Amount: 2},
					},
				},
				Type: zb.AIType_MIXED_AI,
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

func TestCardCollectionCardOperations(t *testing.T) {
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

	CardCollectionCard, err := c.GetCollection(ctx, &zb.GetCollectionRequest{
		UserId: "CardUser",
	})
	assert.Nil(t, err)
	assert.Equal(t, 12, len(CardCollectionCard.Cards))

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
				Name:       "NewDeck",
				OverlordId: 1,
				Cards: []*zb.DeckCard{
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
				Name:       "NewDeck",
				OverlordId: 1,
				Cards: []*zb.DeckCard{
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
				Name:       "NewDeck",
				OverlordId: 1,
				Cards: []*zb.DeckCard{
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
				Name:       "Default",
				OverlordId: 1,
				Cards: []*zb.DeckCard{
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
				Name:       "nEWdECK",
				OverlordId: 1,
				Cards: []*zb.DeckCard{
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
				Id:         2,
				Name:       "Edited",
				OverlordId: 1,
				Cards: []*zb.DeckCard{
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
				Id:         2,
				Name:       "Edited",
				OverlordId: 1,
				Cards: []*zb.DeckCard{
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
				Id:         2,
				Name:       "Edited",
				OverlordId: 1,
				Cards: []*zb.DeckCard{
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
				Id:         2,
				Name:       "dEFAULT",
				OverlordId: 1,
				Cards: []*zb.DeckCard{
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
		// we expect all the cards in InitRequest.Cards
		assert.Equal(t, 90, len(cardResponse.Cards))
	})

	t.Run("ListOverlordLibrary", func(t *testing.T) {
		overlordsResponse, err := c.ListOverlordLibrary(ctx, &zb.ListOverlordLibraryRequest{
			Version: "v1",
		})

		assert.Nil(t, err)
		assert.Equal(t, 2, len(overlordsResponse.Overlords))
	})
}

func TestOverlordsOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "7696b824516b283f81ea1747fbddbe73fe4b5fce0eac0728e47de41d8e306701"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "OverlordUser",
		Image:   "PathToImage",
		Version: "v1",
	}, t)

	t.Run("ListOverlords", func(t *testing.T) {
		overlordsResponse, err := c.ListOverlords(ctx, &zb.ListOverlordsRequest{
			UserId: "OverlordUser",
		})

		assert.Nil(t, err)
		assert.Equal(t, 2, len(overlordsResponse.Overlords))
	})

	t.Run("GetOverlord", func(t *testing.T) {
		overlordResponse, err := c.GetOverlord(ctx, &zb.GetOverlordRequest{
			UserId:     "OverlordUser",
			OverlordId: 1,
		})

		assert.Nil(t, err)
		assert.NotNil(t, overlordResponse.Overlord)
	})

	t.Run("GetOverlord (Overlord not exists)", func(t *testing.T) {
		_, err := c.GetOverlord(ctx, &zb.GetOverlordRequest{
			UserId:     "OverlordUser",
			OverlordId: 10,
		})

		assert.NotNil(t, err)
	})

	t.Run("AddOverlordExperience", func(t *testing.T) {
		resp, err := c.AddOverlordExperience(ctx, &zb.AddOverlordExperienceRequest{
			UserId:     "OverlordUser",
			OverlordId: 0,
			Experience: 2,
		})

		assert.Nil(t, err)
		assert.Equal(t, int64(2), resp.Experience)
	})

	t.Run("AddOverlordExperience (Negative experience)", func(t *testing.T) {
		_, err := c.AddOverlordExperience(ctx, &zb.AddOverlordExperienceRequest{
			UserId:     "OverlordUser",
			OverlordId: 0,
			Experience: -2,
		})

		assert.NotNil(t, err)
	})

	t.Run("AddOverlordExperience (Non existant overlord)", func(t *testing.T) {
		_, err := c.AddOverlordExperience(ctx, &zb.AddOverlordExperienceRequest{
			UserId:     "OverlordUser",
			OverlordId: 100,
			Experience: 2,
		})

		assert.NotNil(t, err)
	})

	t.Run("GetOverlordSkills", func(t *testing.T) {
		skillResponse, err := c.GetOverlordSkills(ctx, &zb.GetOverlordSkillsRequest{
			UserId:     "OverlordUser",
			OverlordId: 0,
		})

		assert.Nil(t, err)
		assert.Equal(t, 1, len(skillResponse.Skills))
	})

	t.Run("GetOverlordSkills (Non existant overlord)", func(t *testing.T) {
		_, err := c.GetOverlordSkills(ctx, &zb.GetOverlordSkillsRequest{
			UserId:     "OverlordUser",
			OverlordId: 100,
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

	t.Run("UpdateInit with missing data", func(t *testing.T) {
		updateInitRequestWithMissingCards := proto.Clone(&updateInitRequest).(*zb.UpdateInitRequest)
		updateInitRequestWithMissingCards.InitData.Cards = nil
		err := c.UpdateInit(ctx, updateInitRequestWithMissingCards)

		assert.NotNil(t, err)
		assert.Equal(t, "'cards' key missing", err.Error())
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
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
			},
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
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
			},
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
		response, err := c.GetMatch(ctx, &zb.GetMatchRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.Equal(t, zb.Match_Canceled, response.Match.Status)
	})
}

func TestCancelFindMatchOnEndedMatchOperations(t *testing.T) {
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
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
			},
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

	t.Run("EndMatch", func(t *testing.T) {
		response, err := c.EndMatch(ctx, &zb.EndMatchRequest{
			UserId:   "player-2",
			MatchId:  matchID,
			WinnerId: "player-2",
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
	})

	t.Run("GetMatch", func(t *testing.T) {
		response, err := c.GetMatch(ctx, &zb.GetMatchRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.Equal(t, zb.Match_Ended, response.Match.Status)
	})

	t.Run("CancelFindmatch", func(t *testing.T) {
		_, err := c.CancelFindMatch(ctx, &zb.CancelFindMatchRequest{
			UserId:  "player-1",
			MatchId: matchID,
		})
		assert.Nil(t, err)
	})

	t.Run("GetMatch", func(t *testing.T) {
		response, err := c.GetMatch(ctx, &zb.GetMatchRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.Equal(t, zb.Match_Ended, response.Match.Status)
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
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
			},
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
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1-tag",
				Version: "v1",
				Tags:    tags,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPoolTag", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2-tag",
				Version: "v1",
				Tags:    tags,
			},
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
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-3",
				Version: "v1",
				Tags:    tags,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPoolTag", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-3-tag",
				Version: "v1",
				Tags:    tags,
			},
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

func TestFindMatchWithTagGroupOperations(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)
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
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-4-tag",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-5-tag",
		Version: "v1",
	}, t)
	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "player-6-tag",
		Version: "v1",
	}, t)

	tags1 := []string{"tags1"}
	tags2 := []string{"tags2"}
	tags3 := []string{"tags3", "othertag"}

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1-tag",
				Version: "v1",
				Tags:    tags1,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-3-tag",
				Version: "v1",
				Tags:    tags2,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-5-tag",
				Version: "v1",
				Tags:    tags3,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-1-tag",
			Tags:   tags1,
		})
		assert.Nil(t, err)
		assert.Equal(t, false, response.MatchFound)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-3-tag",
			Tags:   tags2,
		})
		assert.Nil(t, err)
		assert.Equal(t, false, response.MatchFound)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-5-tag",
			Tags:   tags3,
		})
		assert.Nil(t, err)
		assert.Equal(t, false, response.MatchFound)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-6-tag",
				Version: "v1",
				Tags:    tags3,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2-tag",
				Version: "v1",
				Tags:    tags1,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-4-tag",
				Version: "v1",
				Tags:    tags2,
			},
		})
		assert.Nil(t, err)
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-1-tag",
			Tags:   tags1,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, "player-1-tag", response.Match.PlayerStates[0].Id, "Player should be player-1-tag")
		assert.Equal(t, "player-2-tag", response.Match.PlayerStates[1].Id, "Player should be player-2-tag")
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-3-tag",
			Tags:   tags2,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, "player-3-tag", response.Match.PlayerStates[0].Id, "Player should be player-3-tag")
		assert.Equal(t, "player-4-tag", response.Match.PlayerStates[1].Id, "Player should be player-4-tag")
	})

	t.Run("Findmatch", func(t *testing.T) {
		response, err := c.FindMatch(ctx, &zb.FindMatchRequest{
			UserId: "player-5-tag",
			Tags:   tags3,
		})
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Match.PlayerStates), "the player should see 2 player states")
		assert.Equal(t, zb.Match_Matching, response.Match.Status, "match status should be 'matching'")
		assert.Equal(t, "player-5-tag", response.Match.PlayerStates[0].Id, "Player should be player-5-tag")
		assert.Equal(t, "player-6-tag", response.Match.PlayerStates[1].Id, "Player should be player-6-tag")
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
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  fmt.Sprintf("player-%d", i+1),
				Version: "v1",
			},
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
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-3",
				Version: "v1",
			},
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
		response, err := c.GetMatch(ctx, &zb.GetMatchRequest{
			MatchId: matchID,
		})
		assert.Nil(t, err)
		assert.Equal(t, zb.Match_Timedout, response.Match.Status)
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
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
			},
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

		gameState, err := loadGameState(ctx, matchID)
		assert.EqualValues(t, 1, gameState.CurrentPlayerIndex, "player-2 should be active")
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

		gameState, err := loadGameState(ctx, matchID)
		assert.EqualValues(t, 0, gameState.CurrentPlayerIndex, "player-1 should be active")
	})
	t.Run("SendCardPlayPlayer1", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_CardPlay,
				PlayerId:   "player-1",
				Action: &zb.PlayerAction_CardPlay{
					CardPlay: &zb.PlayerActionCardPlay{
						Card: &zb.InstanceId{Id: 8},
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
						Card: &zb.InstanceId{Id: 1},
						Targets: []*zb.Unit{
							&zb.Unit{
								InstanceId: &zb.InstanceId{Id: 2},
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
						SkillId: 1,
						Target: &zb.Unit{
							InstanceId: &zb.InstanceId{Id: 2},
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
						Card: &zb.InstanceId{Id: 1},
						Targets: []*zb.Unit{
							&zb.Unit{
								InstanceId: &zb.InstanceId{Id: 2},
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

		gameState, err := loadGameState(ctx, matchID)
		assert.EqualValues(t, 1, gameState.CurrentPlayerIndex, "player-2 should be active")
	})
	t.Run("SendCardPlayPlayer2", func(t *testing.T) {
		response, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: matchID,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_CardPlay,
				PlayerId:   "player-2",
				Action: &zb.PlayerAction_CardPlay{
					CardPlay: &zb.PlayerActionCardPlay{
						Card: &zb.InstanceId{Id: 13},
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
						Attacker: &zb.InstanceId{Id: 13},
						Target: &zb.Unit{
							InstanceId: &zb.InstanceId{Id: 8},
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
						Card: &zb.InstanceId{Id: 1},
						Targets: []*zb.Unit{
							&zb.Unit{
								InstanceId: &zb.InstanceId{Id: 2},
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
							InstanceId: &zb.InstanceId{Id: 2},
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
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
			},
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
						Card: &zb.InstanceId{Id: 8},
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
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
				DebugCheats: zb.DebugCheatsConfiguration{
					Enabled:             true,
					UseCustomRandomSeed: true,
					CustomRandomSeed:    2,
				},
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
				DebugCheats: zb.DebugCheatsConfiguration{
					Enabled:             true,
					UseCustomRandomSeed: true,
					CustomRandomSeed:    2,
				},
			},
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

	t.Run("Player1_Leave", func(t *testing.T) {
		_, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: 1,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_LeaveMatch,
				PlayerId:   "player-1",
				Action: &zb.PlayerAction_LeaveMatch{
					LeaveMatch: &zb.PlayerActionLeaveMatch{
						Reason: zb.PlayerActionLeaveMatch_PlayerLeave,
					},
				},
			},
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
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
				DebugCheats: zb.DebugCheatsConfiguration{
					Enabled:             true,
					UseCustomRandomSeed: true,
					CustomRandomSeed:    2,
				},
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
				DebugCheats: zb.DebugCheatsConfiguration{
					Enabled:             true,
					UseCustomRandomSeed: true,
					CustomRandomSeed:    2,
				},
			},
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
		_, err := c.SendPlayerAction(ctx, &zb.PlayerActionRequest{
			MatchId: 1,
			PlayerAction: &zb.PlayerAction{
				ActionType: zb.PlayerActionType_LeaveMatch,
				PlayerId:   "player-1",
				Action: &zb.PlayerAction_LeaveMatch{
					LeaveMatch: &zb.PlayerActionLeaveMatch{
						Reason: zb.PlayerActionLeaveMatch_PlayerLeave,
					},
				},
			},
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
	aiDecks := []*zb.AIDeck{
		{
			Deck: &zb.Deck{
				Id:         1,
				OverlordId: 2,
				Name:       "AI Decks",
				Cards: []*zb.DeckCard{
					{CardName: "Wheezy", Amount: 2},
					{CardName: "Ztormmcaller", Amount: 2},
					{CardName: "Ztormmcaller", Amount: 2},
					{CardName: "Soothsayer", Amount: 2},
					{CardName: "Wheezy", Amount: 2},
					{CardName: "Gaz", Amount: 2},
					{CardName: "Whizpar", Amount: 1},
					{CardName: "Bolderr", Amount: 1},
					{CardName: "Bouncer", Amount: 1},
					{CardName: "Blocker", Amount: 1},
					{CardName: "Pit", Amount: 1},
					{CardName: "Golem", Amount: 1},
				},
			},
			Type: zb.AIType_MIXED_AI,
		},
	}

	t.Run("Get AI Decks", func(t *testing.T) {
		err := saveAIDecks(ctx, "v1", &zb.AIDeckList{Decks: aiDecks})
		_, err = c.GetAIDecks(ctx, &zb.GetAIDecksRequest{
			Version: "v1",
		})
		assert.Nil(t, err)
	})

	aiDecks = []*zb.AIDeck{
		{
			Deck: &zb.Deck{
				Id:         1,
				OverlordId: 2,
				Name:       "AI Decks",
				Cards: []*zb.DeckCard{
					{CardName: "NoName", Amount: 2},
				},
			},
			Type: zb.AIType_MIXED_AI,
		},
	}
}

func TestKeepAlive(t *testing.T) {
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
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-1",
				Version: "v1",
				DebugCheats: zb.DebugCheatsConfiguration{
					Enabled:             true,
					UseCustomRandomSeed: true,
					CustomRandomSeed:    2,
				},
			},
		})
		assert.Nil(t, err)
	})

	t.Run("RegisterPlayerPool", func(t *testing.T) {
		_, err := c.RegisterPlayerPool(ctx, &zb.RegisterPlayerPoolRequest{
			RegistrationData: &zb.PlayerProfileRegistrationData{
				DeckId:  1,
				UserId:  "player-2",
				Version: "v1",
				DebugCheats: zb.DebugCheatsConfiguration{
					Enabled:             true,
					UseCustomRandomSeed: true,
					CustomRandomSeed:    2,
				},
			},
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

	t.Run("KeepAlivePlayer1_Success", func(t *testing.T) {
		_, err := c.KeepAlive(ctx, &zb.KeepAliveRequest{
			MatchId: 1,
			UserId:  "player-1",
		})
		assert.Nil(t, err)
	})

	t.Run("KeepAlivePlayer2_Success", func(t *testing.T) {
		_, err := c.KeepAlive(ctx, &zb.KeepAliveRequest{
			MatchId: 1,
			UserId:  "player-2",
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

	// keep player2 alive
	fc.SetTime(now.Add(time.Second * 5))
	t.Run("KeepAlivePlayer2_Success", func(t *testing.T) {
		_, err := c.KeepAlive(ctx, &zb.KeepAliveRequest{
			MatchId: 1,
			UserId:  "player-2",
		})
		assert.Nil(t, err)
	})

	// player1 fails to keep alive
	// move time forward to expire the player's turn
	fc.SetTime(now.Add(KeepAliveTimeout + time.Second*5))

	t.Run("KeepAlivePlayer2_Winner", func(t *testing.T) {
		_, err := c.KeepAlive(ctx, &zb.KeepAliveRequest{
			MatchId: 1,
			UserId:  "player-2",
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

	t.Run("KeepAliveAfterGameEndedShouldNot_Fail", func(t *testing.T) {
		_, err := c.KeepAlive(ctx, &zb.KeepAliveRequest{
			MatchId: 1,
			UserId:  "player-1",
		})
		assert.Nil(t, err)
	})

	t.Run("KeepAliveAfterGameEndedShouldNot_Fail", func(t *testing.T) {
		_, err := c.KeepAlive(ctx, &zb.KeepAliveRequest{
			MatchId: 1,
			UserId:  "player-2",
		})
		assert.Nil(t, err)
	})

	t.Run("SendAnyActionShould_Fail", func(t *testing.T) {
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

func TestRewardTutorialCompleted(t *testing.T) {
	c := &ZombieBattleground{}
	var pubKeyHexString = "3866f776276246e4f9998aa90632931d89b0d3a5930e804e02299533f55b39e1"
	var addr loom.Address
	var ctx contract.Context

	setup(c, pubKeyHexString, &addr, &ctx, t)

	setupAccount(c, ctx, &zb.UpsertAccountRequest{
		UserId:  "loom1",
		Version: "v1",
	}, t)

	privateKeyStr = "757fc001c98d83eb8288d6c5294f31c284f1c83dbdbc516e3062365f682ffd8a"

	// make sure we have the correct reward_contract_version and tutorial_reward_amount
	resp, err := c.GetState(ctx, &zb.GetGamechainStateRequest{})
	assert.Nil(t, err)
	assert.EqualValues(t, 1, resp.State.RewardContractVersion)
	assert.EqualValues(t, 1, resp.State.TutorialRewardAmount)
	var rewardTutorialCompletedResp *zb.RewardTutorialCompletedResponse

	// get reward on completing tutorial
	t.Run("RewardTutorialCompleted", func(t *testing.T) {
		var err error
		rewardTutorialCompletedResp, err = c.RewardTutorialCompleted(ctx, &zb.RewardTutorialCompletedRequest{})
		assert.Nil(t, err)
		assert.NotNil(t, rewardTutorialCompletedResp)
		assert.Equal(t, "0xcc69885fda6bcc1a4ace058b4a62bf5e179ea78fd58a1ccd71c22cc9b688792f", rewardTutorialCompletedResp.Hash)
		assert.Equal(t, "0x25d667d7a5afa5fc9cfacf67aeaf35bff6b3d54bf40e8b3b5a90fe86ce596732", rewardTutorialCompletedResp.R)
		assert.Equal(t, "0x02a59d7fd35f07738ac2931193b92e8f31941c15a123d07f0990c28a1fef2760", rewardTutorialCompletedResp.S)
		assert.Equal(t, uint64(27), rewardTutorialCompletedResp.V)
	})

	// send the correct request
	t.Run("ConfirmRewardClaimed", func(t *testing.T) {
		err := c.ConfirmRewardTutorialClaimed(ctx, &zb.ConfirmRewardTutorialClaimedRequest{})
		assert.Nil(t, err)
	})

	// attempt to claim reward again should fail
	t.Run("RewardTutorialCompleted again should fail", func(t *testing.T) {
		var err error
		rewardTutorialCompletedResp, err = c.RewardTutorialCompleted(ctx, &zb.RewardTutorialCompletedRequest{})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "reward already claimed")
		assert.Nil(t, rewardTutorialCompletedResp)
	})

	// attempt to get reward again should fail
	t.Run("repeated RewardTutorialCompleted fails", func(t *testing.T) {
		resp, err := c.RewardTutorialCompleted(ctx, &zb.RewardTutorialCompletedRequest{})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "reward already claimed")
		assert.Nil(t, resp)
	})
}

func TestHashSignature(t *testing.T) {
	privateKeyStr = "921660bf3e5c7a404beed663f00462645fd8d50751d21e262f6f1a3b7e5b5da3"
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	assert.Nil(t, err)

	verifySignResult, err := generateVerifyHash(5, 1, privateKey)
	assert.Nil(t, err)
	assert.Equal(t, "0xe2689cd4a84e23ad2f564004f1c9013e9589d260bde6380aba3ca7e09e4df40c", verifySignResult.Hash)
	assert.Equal(t, "0x0e729720a48e6164451792e657bd4c68c25c67884e8852790d0cf2fbe999f2bb2833d41711917583d70dacbb5c80206edafa3bc4bca079ceafbb2fc50af1c8fc1b", verifySignResult.Signature)
}

// createJwtToken creates a jwt token from a User model
func createJwtToken(userID uint, secret string) (string, error) {
	// create the token
	claims := latypes.UserClaims{
		UserID: userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign and get the complete encoded token as string
	return token.SignedString([]byte(secret))
}

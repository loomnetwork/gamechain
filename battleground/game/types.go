//go:generate $PWD/bin/pbgraphserialization-gen --targetPackagePath github.com/loomnetwork/gamechain/battleground/game/ --targetPackageName game --protoPackageName zb --outputPath types_serialization.go

package game

import _ "github.com/loomnetwork/gamechain/types/zb"
import _ "github.com/loomnetwork/gamechain/library/pbgraphserialization"

//pbgraphserialization:enable
type DeckCard struct {
	cardName string
	amount int64
}

//pbgraphserialization:enable
type Deck struct {
	id             int64
	name           string
	heroId         int64
	cards          []*DeckCard
	primarySkill   OverlordSkillKind_Enum
	secondarySkill OverlordSkillKind_Enum
}

type OverlordSkillKind_Enum int32

const (
	OverlordSkillKind_None OverlordSkillKind_Enum = 0
	// AIR
	OverlordSkillKind_Push       OverlordSkillKind_Enum = 1
	OverlordSkillKind_Draw       OverlordSkillKind_Enum = 2
	OverlordSkillKind_WindShield OverlordSkillKind_Enum = 3
	OverlordSkillKind_WindWall   OverlordSkillKind_Enum = 4
	OverlordSkillKind_Retreat    OverlordSkillKind_Enum = 5
	// EARTH
	OverlordSkillKind_Harden    OverlordSkillKind_Enum = 6
	OverlordSkillKind_StoneSkin OverlordSkillKind_Enum = 7
	OverlordSkillKind_Fortify   OverlordSkillKind_Enum = 8
	OverlordSkillKind_Phalanx   OverlordSkillKind_Enum = 9
	OverlordSkillKind_Fortress  OverlordSkillKind_Enum = 10
	// FIRE
	OverlordSkillKind_FireBolt     OverlordSkillKind_Enum = 11
	OverlordSkillKind_Rabies       OverlordSkillKind_Enum = 12
	OverlordSkillKind_Fireball     OverlordSkillKind_Enum = 13
	OverlordSkillKind_MassRabies   OverlordSkillKind_Enum = 14
	OverlordSkillKind_MeteorShower OverlordSkillKind_Enum = 15
	// LIFE
	OverlordSkillKind_HealingTouch OverlordSkillKind_Enum = 16
	OverlordSkillKind_Mend         OverlordSkillKind_Enum = 17
	OverlordSkillKind_Ressurect    OverlordSkillKind_Enum = 18
	OverlordSkillKind_Enhance      OverlordSkillKind_Enum = 19
	OverlordSkillKind_Reanimate    OverlordSkillKind_Enum = 20
	// TOXIC
	OverlordSkillKind_PoisonDart OverlordSkillKind_Enum = 21
	OverlordSkillKind_ToxicPower OverlordSkillKind_Enum = 22
	OverlordSkillKind_Breakout   OverlordSkillKind_Enum = 23
	OverlordSkillKind_Infect     OverlordSkillKind_Enum = 24
	OverlordSkillKind_Epidemic   OverlordSkillKind_Enum = 25
	// WATER
	OverlordSkillKind_Freeze   OverlordSkillKind_Enum = 26
	OverlordSkillKind_IceBolt  OverlordSkillKind_Enum = 27
	OverlordSkillKind_IceWall  OverlordSkillKind_Enum = 28
	OverlordSkillKind_Shatter  OverlordSkillKind_Enum = 29
	OverlordSkillKind_Blizzard OverlordSkillKind_Enum = 30
)
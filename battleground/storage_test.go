package battleground

import (
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/battleground/battleground_nullable"
	"github.com/loomnetwork/gamechain/types/nullable/nullable_pb"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCard = zb_data.Card{
	MouldId:     3,
	Kind:        zb_enums.CardKind_Creature,
	Faction:     zb_enums.Faction_Earth,
	Name:        "Zpitter",
	Description: "Amazing zpit of unfathomeable power.",
	FlavorText:  "Zpittity-zpit",
	Picture:     "zpitter.png",
	Rank:        zb_enums.CreatureRank_Commander,
	Type:        zb_enums.CardType_Feral,
	Frame:       "normal-frame.png",
	Damage:      3,
	Defense:     4,
	Cost:        5,
	PictureTransform: &zb_data.PictureTransform{
		Position: &zb_data.Vector3Float{
			X: 0.1,
			Y: 0.2,
			Z: 0.3,
		},
		Scale: &zb_data.Vector3Float{
			X: 0.7,
			Y: 0.8,
			Z: 0.9,
		},
	},
	Abilities: []*zb_data.AbilityData{
		{
			Name:   "Super Ability",
			Cost:   3,
			Effect: zb_enums.AbilityEffect_HealDirectly,
		},
	},
	UniqueAnimation: zb_enums.UniqueAnimation_ChernoBillArrival,
	Hidden:          true,
	SourceMouldId:   0,
	Overrides:       nil,
}

func TestSourceIdBasic(t *testing.T) {
	targetCard := zb_data.Card{
		MouldId:       4,
		SourceMouldId: testCard.MouldId,
	}
	var cardLibrary = []*zb_data.Card{
		&testCard,
		&targetCard,
	}

	err := validateCardLibraryCards(cardLibrary)
	assert.Nil(t, err)

	mouldIdToCard, err := getMouldIdToCardMap(cardLibrary)
	assert.Nil(t, err)

	for _, card := range cardLibrary {
		err = applySourceMouldIdAndOverrides(card, mouldIdToCard)
		assert.Nil(t, err)
	}

	assert.Equal(t, 4, int(targetCard.MouldId))
	assert.Equal(t, testCard.MouldId, targetCard.SourceMouldId)
	assert.Equal(t, "Zpitter", targetCard.Name)

	json, err := protoMessageToJSON(&zb_data.CardList{Cards: cardLibrary})
	fmt.Println(json)
	//assert.Nil(t, json)
}

func TestSourceIdOverride(t *testing.T) {
	targetCard := zb_data.Card{
		MouldId:       4,
		SourceMouldId: testCard.MouldId,
		Overrides: &zb_data.CardOverrides{
			Name:       &nullable_pb.StringValue{Value: "Legendary Zpitter"},
			FlavorText: &nullable_pb.StringValue{Value: "Zpittity-zpit, now with more zpit"},
			Picture:    &nullable_pb.StringValue{Value: "zpitter_legendary.png"},
			Rank:       &battleground_nullable.CreatureRankEnumValue{Value: zb_enums.CreatureRank_General},
			Type:       &battleground_nullable.CardTypeEnumValue{Value: zb_enums.CardType_Heavy},
			Frame:      &nullable_pb.StringValue{Value: "legendary-frame.png"},
			Hidden:     &nullable_pb.BoolValue{Value: false},
		},
	}
	var cardLibrary = []*zb_data.Card{
		&testCard,
		&targetCard,
	}

	err := validateCardLibraryCards(cardLibrary)
	assert.Nil(t, err)

	mouldIdToCard, err := getMouldIdToCardMap(cardLibrary)
	assert.Nil(t, err)

	for _, card := range cardLibrary {
		err = applySourceMouldIdAndOverrides(card, mouldIdToCard)
		assert.Nil(t, err)
	}

	assert.Equal(t, 4, int(targetCard.MouldId))
	assert.Equal(t, testCard.MouldId, targetCard.SourceMouldId)
	assert.Equal(t, "Legendary Zpitter", targetCard.Name)
	assert.Equal(t, "Zpittity-zpit, now with more zpit", targetCard.FlavorText)
	assert.Equal(t, "zpitter_legendary.png", targetCard.Picture)
	assert.Equal(t, zb_enums.CreatureRank_General, targetCard.Rank)
	assert.Equal(t, zb_enums.CardType_Heavy, targetCard.Type)
	assert.Equal(t, "legendary-frame.png", targetCard.Frame)
	assert.Equal(t, false, targetCard.Hidden)

	json, err := protoMessageToJSON(&zb_data.CardList{Cards: cardLibrary})
	fmt.Println(json)
}

func protoMessageToJSON(pb proto.Message) (string, error) {
	m := jsonpb.Marshaler{
		OrigName:     false,
		Indent:       "  ",
		EmitDefaults: true,
	}

	json, err := m.MarshalToString(pb)
	if err != nil {
		return "", fmt.Errorf("error marshaling Proto to JSON: %s", err.Error())
	}

	return json, nil
}

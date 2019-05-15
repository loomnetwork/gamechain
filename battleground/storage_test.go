package battleground

import (
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSourceIdBasic(t *testing.T) {
	var cardLibrary = []*zb.Card{
		{
			MouldId: 3,
			Kind: zb_enums.CardKind_Creature,
			Faction: zb.Faction_Earth,
			Name: "Zpitter",
			Description: "Amazing zpit of unfathomeable power.",
			FlavorText: "Zpittity-zpit",
			Picture: "zpitter.png",
			Rank: zb.CreatureRank_Commander,
			Type: zb.CardType_Feral,
			Frame: "normal-frame.png",
			Damage: 3,
			Defense: 4,
			Cost: 5,
			PictureTransform: &zb.PictureTransform{
				Position: &zb.Vector3Float{
					X: 0.1,
					Y: 0.2,
					Z: 0.3,
				},
			},
			Abilities: []*zb.AbilityData{}, // FIXME
			UniqueAnimation: zb.UniqueAnimation_ChernoBillArrival,
			Hidden: true,
			SourceMouldId: 0,
			Overrides: nil,
		},
		{
			SourceMouldId: 3,
		},
	}

	err := validateCardLibraryCards(cardLibrary)
	assert.Nil(t, err)

	mouldIdToCard, err := getMouldIdToCardMap(cardLibrary)
	assert.Nil(t, err)

	for _, card := range cardLibrary {
		err = applySourceMouldIdAndOverrides(card, mouldIdToCard)
		assert.Nil(t, err)
	}

	json, err := protoMessageToJSON(&zb_data.CardList{Cards: cardLibrary})
	fmt.Println(json)
	assert.Nil(t, json)
}

func protoMessageToJSON(pb proto.Message) (string, error) {
	m := jsonpb.Marshaler{
		OrigName:     false,
		Indent:       "",
		EmitDefaults: true,
	}

	json, err := m.MarshalToString(pb)
	if err != nil {
		return "", fmt.Errorf("error marshaling Proto to JSON: %s", err.Error())
	}

	return json, nil
}
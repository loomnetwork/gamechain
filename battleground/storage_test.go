package battleground

import (
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSourceIdBasic(t *testing.T) {
	var cardLibrary = []*zb_data.Card{
		{
			MouldId: 3,
			Kind: zb_enums.CardKind_Creature,
			Faction: zb_enums.Faction_Earth,
			Name: "Zpitter",
			Description: "Amazing zpit of unfathomeable power.",
			FlavorText: "Zpittity-zpit",
			Picture: "zpitter.png",
			Rank: zb_enums.CreatureRank_Commander,
			Type: zb_enums.CardType_Feral,
			Frame: "normal-frame.png",
			Damage: 3,
			Defense: 4,
			Cost: 5,
			PictureTransform: &zb_data.PictureTransform{
				Position: &zb_data.Vector3Float{
					X: 0.1,
					Y: 0.2,
					Z: 0.3,
				},
			},
			Abilities: []*zb_data.AbilityData{}, // FIXME
			UniqueAnimation: zb_enums.UniqueAnimation_ChernoBillArrival,
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
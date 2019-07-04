package battleground

import (
	battleground_proto "github.com/loomnetwork/gamechain/battleground/proto"
	"github.com/loomnetwork/gamechain/battleground/proto/battleground_nullable"
	"github.com/loomnetwork/gamechain/types/nullable/nullable_pb"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	"github.com/loomnetwork/gamechain/types/zb/zb_enums"
	"testing"

	assert "github.com/stretchr/testify/require"
)

var testCard = zb_data.Card{
	CardKey:     battleground_proto.CardKey{MouldId: 3},
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
	Overrides:       nil,
}

func TestValidateDeckCollection(t *testing.T) {
	userHas := []*zb_data.CardCollectionCard{
		{CardKey: battleground_proto.CardKey{MouldId: 90}, Amount: 4},
		{CardKey: battleground_proto.CardKey{MouldId: 91}, Amount: 3},
		{CardKey: battleground_proto.CardKey{MouldId: 96}, Amount: 5},
		{CardKey: battleground_proto.CardKey{MouldId: 3}, Amount: 4},
	}

	t.Run("Successful validation", func(t *testing.T) {
		newCollection := []*zb_data.DeckCard{
			{CardKey: battleground_proto.CardKey{MouldId: 90}, Amount: 4},
			{CardKey: battleground_proto.CardKey{MouldId: 91}, Amount: 3},
			{CardKey: battleground_proto.CardKey{MouldId: 96}, Amount: 5},
			{CardKey: battleground_proto.CardKey{MouldId: 3}, Amount: 4},
		}
		assert.Nil(t, validateDeckAgainstUserCardCollection(userHas, newCollection))
	})

	t.Run("Successful validation", func(t *testing.T) {
		newCollection := []*zb_data.DeckCard{
			{CardKey: battleground_proto.CardKey{MouldId: 90}, Amount: 0},
			{CardKey: battleground_proto.CardKey{MouldId: 91}, Amount: 0},
			{CardKey: battleground_proto.CardKey{MouldId: 96}, Amount: 0},
			{CardKey: battleground_proto.CardKey{MouldId: 3}, Amount: 0},
		}
		assert.Nil(t, validateDeckAgainstUserCardCollection(userHas, newCollection))
	})

	t.Run("Successful validation", func(t *testing.T) {
		newCollection := []*zb_data.DeckCard{}
		assert.Nil(t, validateDeckAgainstUserCardCollection([]*zb_data.CardCollectionCard{}, newCollection))
	})

	t.Run("Failed validation", func(t *testing.T) {
		newCollection := []*zb_data.DeckCard{
			{CardKey: battleground_proto.CardKey{MouldId: 90}, Amount: 8},
			{CardKey: battleground_proto.CardKey{MouldId: 91}, Amount: 10},
		}
		assert.NotNil(t, validateDeckAgainstUserCardCollection(userHas, newCollection))
	})

	t.Run("Failed validation", func(t *testing.T) {
		newCollection := []*zb_data.DeckCard{
			{CardKey: battleground_proto.CardKey{MouldId: -2}, Amount: 0},
			{CardKey: battleground_proto.CardKey{MouldId: -3}, Amount: 0},
		}
		assert.NotNil(t, validateDeckAgainstUserCardCollection(userHas, newCollection))
	})

	t.Run("Failed validation", func(t *testing.T) {
		newCollection := []*zb_data.DeckCard{
			{CardKey: battleground_proto.CardKey{MouldId: 90}, Amount: 8},
			{CardKey: battleground_proto.CardKey{MouldId: 91}, Amount: 10},
		}
		assert.NotNil(t, validateDeckAgainstUserCardCollection([]*zb_data.CardCollectionCard{}, newCollection))
	})
}

func TestValidateDeckName(t *testing.T) {
	tests := []struct {
		name string
		id   int64
		werr error
	}{
		{name: "deck1", werr: ErrDeckNameExists},
		{name: "deck2", werr: nil},
		{name: "การ์ดอันที่หนึ่ง", werr: nil},
		{name: "", werr: ErrDeckNameEmpty},
		{name: "卡1", werr: nil},
		{name: "12345678901234567890123456789012345678901234567890", werr: ErrDeckNameTooLong},
		{name: "การ์ดอันที่หนึ่ง ชื่อยาวเกินไปจริงๆนะ ชื่อยาวเกินไปจริงๆนะ ชื่อยาวเกินไปจริงๆนะ ", werr: ErrDeckNameTooLong},
		{name: "deck1", id: 1, werr: nil},
		{name: "deck1", id: 2, werr: ErrDeckNameExists},
		{name: "DECK1", id: 2, werr: ErrDeckNameExists},
	}

	existingDecks := []*zb_data.Deck{
		{Id: 1, Name: "deck1"},
	}
	for _, test := range tests {
		err := validateDeckName(existingDecks, &zb_data.Deck{Name: test.name, Id: test.id})
		assert.Equal(t, test.werr, err)
	}
}

func TestVariantBasic(t *testing.T) {
	targetCard := zb_data.Card{
		CardKey: battleground_proto.CardKey{
			MouldId: 3,
			Variant: zb_enums.CardVariant_Limited,
		},
		Type:    zb_enums.CardType_Walker,
		Kind:    zb_enums.CardKind_Creature,
		Faction: zb_enums.Faction_Earth,
		Rank:    zb_enums.CreatureRank_Minion,
	}
	var cardLibrary = []*zb_data.Card{
		&testCard,
		&targetCard,
	}

	err := validateCardLibraryCards(cardLibrary)
	assert.Nil(t, err)

	mouldIdToCard, err := getCardKeyToCardMap(cardLibrary)
	assert.Nil(t, err)

	for _, card := range cardLibrary {
		err = applySourceMouldIdAndOverrides(card, mouldIdToCard)
		assert.Nil(t, err)
	}

	assert.Equal(t, 3, int(targetCard.CardKey.MouldId))
	assert.Equal(t, testCard.CardKey.MouldId, targetCard.CardKey.MouldId)
	assert.Equal(t, zb_enums.CardVariant_Limited, targetCard.CardKey.Variant)
	assert.Equal(t, "Zpitter", targetCard.Name)
}

func TestVariantOverride(t *testing.T) {
	targetCard := zb_data.Card{
		CardKey: battleground_proto.CardKey{
			MouldId: 3,
			Variant: zb_enums.CardVariant_Limited,
		},
		Type:    zb_enums.CardType_Walker,
		Kind:    zb_enums.CardKind_Creature,
		Faction: zb_enums.Faction_Earth,
		Rank:    zb_enums.CreatureRank_Minion,
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

	mouldIdToCard, err := getCardKeyToCardMap(cardLibrary)
	assert.Nil(t, err)

	for _, card := range cardLibrary {
		err = applySourceMouldIdAndOverrides(card, mouldIdToCard)
		assert.Nil(t, err)
	}

	assert.Equal(t, 3, int(targetCard.CardKey.MouldId))
	assert.Equal(t, testCard.CardKey.MouldId, targetCard.CardKey.MouldId)
	assert.Equal(t, zb_enums.CardVariant_Limited, targetCard.CardKey.Variant)
	assert.Equal(t, "Legendary Zpitter", targetCard.Name)
	assert.Equal(t, "Zpittity-zpit, now with more zpit", targetCard.FlavorText)
	assert.Equal(t, "zpitter_legendary.png", targetCard.Picture)
	assert.Equal(t, zb_enums.CreatureRank_General, targetCard.Rank)
	assert.Equal(t, zb_enums.CardType_Heavy, targetCard.Type)
	assert.Equal(t, "legendary-frame.png", targetCard.Frame)
	assert.Equal(t, false, targetCard.Hidden)
}

func TestValidateDeckCardVariants(t *testing.T) {
	cardLibrary := &zb_data.CardLibrary{
		Cards: []*zb_data.Card{
			{
				CardKey: battleground_proto.CardKey{
					MouldId: 5,
					Variant: zb_enums.CardVariant_Standard,
				},
				Type:    zb_enums.CardType_Walker,
				Kind:    zb_enums.CardKind_Creature,
				Faction: zb_enums.Faction_Earth,
				Rank:    zb_enums.CreatureRank_Minion,
			},
			{
				CardKey: battleground_proto.CardKey{
					MouldId: 5,
					Variant: zb_enums.CardVariant_Limited,
				},
				Type:    zb_enums.CardType_Walker,
				Kind:    zb_enums.CardKind_Creature,
				Faction: zb_enums.Faction_Earth,
				Rank:    zb_enums.CreatureRank_Minion,
			},
			{
				CardKey: battleground_proto.CardKey{
					MouldId: 6,
					Variant: zb_enums.CardVariant_Standard,
				},
				Type:    zb_enums.CardType_Walker,
				Kind:    zb_enums.CardKind_Creature,
				Faction: zb_enums.Faction_Earth,
				Rank:    zb_enums.CreatureRank_Minion,
			},
		},
	}

	cardKeyToCardMap, err := getCardKeyToCardMap(cardLibrary.Cards)
	assert.Nil(t, err)

	t.Run("Both Normal and Limited exist in card library", func(t *testing.T) {
		deck := &zb_data.Deck{
			Cards: []*zb_data.DeckCard{
				{
					CardKey: battleground_proto.CardKey{
						MouldId: 5,
						Variant: zb_enums.CardVariant_Standard,
					},
					Amount: 3,
				},
				{
					CardKey: battleground_proto.CardKey{
						MouldId: 5,
						Variant: zb_enums.CardVariant_Limited,
					},
					Amount: 4,
				},
			},
		}

		changed := fixDeckCardVariants(deck, cardKeyToCardMap)
		assert.False(t, changed)
		assert.Equal(t, 2, len(deck.Cards))
		assert.Equal(t, int64(3), deck.Cards[0].Amount)
		assert.Equal(t, int64(4), deck.Cards[1].Amount)
	})

	t.Run("Only Normal exists in card library, but both Normal and Limited is in deck", func(t *testing.T) {
		deck := &zb_data.Deck{
			Cards: []*zb_data.DeckCard{
				{
					CardKey: battleground_proto.CardKey{
						MouldId: 6,
						Variant: zb_enums.CardVariant_Standard,
					},
					Amount: 3,
				},
				{
					CardKey: battleground_proto.CardKey{
						MouldId: 6,
						Variant: zb_enums.CardVariant_Limited,
					},
					Amount: 4,
				},
			},
		}

		changed := fixDeckCardVariants(deck, cardKeyToCardMap)
		assert.True(t, changed)
		assert.Equal(t, 1, len(deck.Cards))
		assert.Equal(t, int64(7), deck.Cards[0].Amount)
		assert.Equal(t, zb_enums.CardVariant_Standard, deck.Cards[0].CardKey.Variant)
	})

	t.Run("Only Normal exists in card library, but only Limited is in deck", func(t *testing.T) {
		deck := &zb_data.Deck{
			Cards: []*zb_data.DeckCard{
				{
					CardKey: battleground_proto.CardKey{
						MouldId: 6,
						Variant: zb_enums.CardVariant_Limited,
					},
					Amount: 4,
				},
			},
		}

		changed := fixDeckCardVariants(deck, cardKeyToCardMap)
		assert.True(t, changed)
		assert.Equal(t, 1, len(deck.Cards))
		assert.Equal(t, int64(4), deck.Cards[0].Amount)
		assert.Equal(t, zb_enums.CardVariant_Standard, deck.Cards[0].CardKey.Variant)
	})

	t.Run("Normal doesn't exist in card library, but Limited is in deck", func(t *testing.T) {
		deck := &zb_data.Deck{
			Cards: []*zb_data.DeckCard{
				{
					CardKey: battleground_proto.CardKey{
						MouldId: 7,
						Variant: zb_enums.CardVariant_Standard,
					},
					Amount: 3,
				},
				{
					CardKey: battleground_proto.CardKey{
						MouldId: 7,
						Variant: zb_enums.CardVariant_Limited,
					},
					Amount: 4,
				},
			},
		}

		changed := fixDeckCardVariants(deck, cardKeyToCardMap)
		assert.True(t, changed)
		assert.Equal(t, 0, len(deck.Cards))
	})

	t.Run("No variants exists in card library, but card is in deck", func(t *testing.T) {
		deck := &zb_data.Deck{
			Cards: []*zb_data.DeckCard{
				{
					CardKey: battleground_proto.CardKey{
						MouldId: 10,
						Variant: zb_enums.CardVariant_Standard,
					},
					Amount: 3,
				},
				{
					CardKey: battleground_proto.CardKey{
						MouldId: 11,
						Variant: zb_enums.CardVariant_Limited,
					},
					Amount: 4,
				},
			},
		}

		changed := fixDeckCardVariants(deck, cardKeyToCardMap)
		assert.True(t, changed)
		assert.Equal(t, 0, len(deck.Cards))
	})
}

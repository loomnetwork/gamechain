package battleground

import (
	"testing"

	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/stretchr/testify/assert"
)

func TestValidateDeckCollection(t *testing.T) {
	userHas := []*zb_data.CardCollectionCard{
		{MouldId: 90, Amount: 4},
		{MouldId: 91, Amount: 3},
		{MouldId: 96, Amount: 5},
		{MouldId: 3, Amount: 4},
	}

	t.Run("Successful validation", func(t *testing.T) {
		newCollection := []*zb_data.CardCollectionCard{
			{MouldId: 90, Amount: 4},
			{MouldId: 91, Amount: 3},
			{MouldId: 96, Amount: 5},
			{MouldId: 3, Amount: 4},
		}
		assert.Nil(t, validateDeckCollections(userHas, newCollection))
	})

	t.Run("Successful validation", func(t *testing.T) {
		newCollection := []*zb_data.CardCollectionCard{
			{MouldId: 90, Amount: 0},
			{MouldId: 91, Amount: 0},
			{MouldId: 96, Amount: 0},
			{MouldId: 3, Amount: 0},
		}
		assert.Nil(t, validateDeckCollections(userHas, newCollection))
	})

	t.Run("Successful validation", func(t *testing.T) {
		newCollection := []*zb_data.CardCollectionCard{}
		assert.Nil(t, validateDeckCollections([]*zb_data.CardCollectionCard{}, newCollection))
	})

	t.Run("Failed validation", func(t *testing.T) {
		newCollection := []*zb_data.CardCollectionCard{
			{MouldId: 90, Amount: 8},
			{MouldId: 91, Amount: 10},
		}
		assert.NotNil(t, validateDeckCollections(userHas, newCollection))
	})

	t.Run("Failed validation", func(t *testing.T) {
		newCollection := []*zb_data.CardCollectionCard{
			{MouldId: -2, Amount: 0},
			{MouldId: -3, Amount: 0},
		}
		assert.NotNil(t, validateDeckCollections(userHas, newCollection))
	})

	t.Run("Failed validation", func(t *testing.T) {
		newCollection := []*zb_data.CardCollectionCard{
			{MouldId: 90, Amount: 8},
			{MouldId: 91, Amount: 10},
		}
		assert.NotNil(t, validateDeckCollections([]*zb_data.CardCollectionCard{}, newCollection))
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

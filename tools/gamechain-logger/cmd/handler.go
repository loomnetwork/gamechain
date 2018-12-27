package cmd

import (
	"fmt"
	"log"

	"github.com/gogo/protobuf/proto"
	"github.com/jinzhu/gorm"
	"github.com/loomnetwork/gamechain/types/zb"
	"github.com/loomnetwork/go-loom/plugin/types"
	"github.com/loomnetwork/loomauth/models"
	"github.com/pkg/errors"
)

type TopicHandler func(eventData *types.EventData, db *gorm.DB) error

func FindMatchHandler(eventData *types.EventData, db *gorm.DB) error {
	var event zb.PlayerActionEvent
	if err := proto.Unmarshal(eventData.EncodedBody, &event); err != nil {
		return err
	}

	match := event.Match
	if match == nil {
		return fmt.Errorf("match is nil")
	}
	if len(match.PlayerStates) < 2 {
		return fmt.Errorf("expected player state length 2")
	}

	m := models.Match{
		ID:              match.Id,
		Player1ID:       match.PlayerStates[0].Id,
		Player2ID:       match.PlayerStates[1].Id,
		Player1Accepted: match.PlayerStates[0].MatchAccepted,
		Player2Accepted: match.PlayerStates[1].MatchAccepted,
		Player1DeckID:   match.PlayerStates[0].Deck.Id,
		Player2DeckID:   match.PlayerStates[1].Deck.Id,
		Status:          match.Status.String(),
		Version:         match.Version,
		RandomSeed:      match.RandomSeed,
	}

	if err := db.Save(&m).Error; err != nil {
		return err
	}

	return nil
}

func AcceptMatchHandler(eventData *types.EventData, db *gorm.DB) error {
	var event zb.PlayerActionEvent
	if err := proto.Unmarshal(eventData.EncodedBody, &event); err != nil {
		return err
	}

	match := event.Match
	if match == nil {
		return fmt.Errorf("match is nil")
	}
	if len(match.PlayerStates) < 2 {
		return fmt.Errorf("expected player state length 2")
	}

	m := models.Match{
		ID:              match.Id,
		Player1ID:       match.PlayerStates[0].Id,
		Player2ID:       match.PlayerStates[1].Id,
		Player1Accepted: match.PlayerStates[0].MatchAccepted,
		Player2Accepted: match.PlayerStates[1].MatchAccepted,
		Player1DeckID:   match.PlayerStates[0].Deck.Id,
		Player2DeckID:   match.PlayerStates[1].Deck.Id,
		Status:          match.Status.String(),
		Version:         match.Version,
		RandomSeed:      match.RandomSeed,
	}
	if err := db.Save(&m).Error; err != nil {
		return err
	}

	return nil
}

func CreateDeckHandler(eventData *types.EventData, db *gorm.DB) error {
	var event zb.CreateDeckEvent
	if err := proto.Unmarshal(eventData.EncodedBody, &event); err != nil {
		return err
	}

	deck := event.Deck
	if deck == nil {
		return fmt.Errorf("deck is nil")
	}

	log.Printf("Saving deck with deck ID %d, userid %s, name %s to DB", event.Deck.Id, event.UserId, event.Deck.Name)

	cards := []models.DeckCard{}
	for _, card := range event.Deck.Cards {
		cards = append(cards, models.DeckCard{
			UserID:   event.UserId,
			CardName: card.CardName,
			Amount:   card.Amount,
		})
	}
	d := models.Deck{
		UserID:           event.UserId,
		DeckID:           deck.Id,
		Name:             deck.Name,
		HeroID:           deck.HeroId,
		DeckCards:        cards,
		PrimarySkillID:   int(deck.PrimarySkill),
		SecondarySkillID: int(deck.SecondarySkill),
		Version:          event.Version,
		SenderAddress:    event.SenderAddress,
	}
	if err := db.Save(&d).Error; err != nil {
		return err
	}

	return nil
}

func EditDeckHandler(eventData *types.EventData, db *gorm.DB) error {
	var event zb.EditDeckEvent
	if err := proto.Unmarshal(eventData.EncodedBody, &event); err != nil {
		return err
	}
	deck := event.Deck
	if deck == nil {
		return fmt.Errorf("deck is nil")
	}

	var edeck models.Deck
	err := db.Where(&models.Deck{UserID: event.UserId, DeckID: event.Deck.Id}).First(&edeck).Error
	if err != nil {
		return err
	}

	cards := []models.DeckCard{}
	for _, card := range event.Deck.Cards {
		cards = append(cards, models.DeckCard{
			CardName: card.CardName,
			Amount:   card.Amount,
		})
	}
	d := models.Deck{
		ID:               edeck.ID,
		UserID:           event.UserId,
		DeckID:           deck.Id,
		Name:             deck.Name,
		HeroID:           deck.HeroId,
		DeckCards:        cards,
		PrimarySkillID:   int(deck.PrimarySkill),
		SecondarySkillID: int(deck.SecondarySkill),
		Version:          event.Version,
		SenderAddress:    event.SenderAddress,
	}
	if err := db.Save(&d).Error; err != nil {
		return err
	}
	return nil
}

func DeleteDeckHandler(eventData *types.EventData, db *gorm.DB) error {
	var event zb.DeleteDeckEvent
	if err := proto.Unmarshal(eventData.EncodedBody, &event); err != nil {
		return err
	}
	log.Printf("Deleting deck with deck ID %d, userid %s from DB", event.DeckId, event.UserId)

	err := db.Where(&models.Deck{UserID: event.UserId, DeckID: event.DeckId}).
		Delete(models.Deck{}).Error
	if err != nil {
		return err
	}
	log.Printf("Deleted deck with deck ID %d, userid %s from DB", event.DeckId, event.UserId)
	return nil
}

func EndgameHandler(eventData *types.EventData, db *gorm.DB) error {
	var event zb.PlayerActionEvent
	if err := proto.Unmarshal(eventData.EncodedBody, &event); err != nil {
		return err
	}

	match := models.Match{}
	err := db.Where(&models.Match{ID: event.Match.Id}).First(&match).Error
	if err != nil {
		return err
	}

	match.WinnerID = event.Block.List[0].GetEndGame().WinnerId
	match.Status = event.Match.Status.String()

	if err := db.Save(&match).Error; err != nil {
		return err
	}
	return nil
}

func MatchHandler(eventData *types.EventData, db *gorm.DB) error {
	var event zb.PlayerActionEvent
	if err := proto.Unmarshal(eventData.EncodedBody, &event); err != nil {
		return err
	}

	match := event.Match
	if match == nil {
		return fmt.Errorf("match is nil")
	}
	topic := fmt.Sprintf("match:%d", match.Id)
	replay, err := writeReplayFile(topic, event)
	if err != nil {
		return errors.Wrap(err, "Error writing replay file: ")
	}

	// update match status
	err = db.Model(&models.Match{}).
		Where(&models.Match{ID: match.Id}).
		Updates(models.Match{Status: match.Status.String()}).
		Error
	if err != nil {
		return err
	}

	dbReplay := models.Replay{}
	err = db.Where(&models.Replay{MatchID: match.Id}).First(&dbReplay).Error
	if err == nil {
		db.First(&dbReplay)
		dbReplay.ReplayJSON = replay
		db.Save(&dbReplay)
	} else if gorm.IsRecordNotFoundError(err) {
		// insert
		dbReplay.MatchID = match.Id
		dbReplay.ReplayJSON = replay
		db.Create(&dbReplay)
	} else {
		return err
	}

	return nil
}

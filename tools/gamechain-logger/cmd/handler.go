package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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
		BlockHeight:     eventData.BlockHeight,
		BlockTime:       time.Unix(eventData.BlockTime, 0),
		CreatedAt:       time.Now(),
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
		BlockHeight:     eventData.BlockHeight,
		BlockTime:       time.Unix(eventData.BlockTime, 0),
	}
	if err := db.Omit("created_at").Save(&m).Error; err != nil {
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

	cards, err := json.Marshal(event.Deck.Cards)
	if err != nil {
		return err
	}

	d := models.Deck{
		UserID:           event.UserId,
		DeckID:           deck.Id,
		Name:             deck.Name,
		HeroID:           deck.OverlordId,
		Cards:            cards,
		PrimarySkillID:   int(deck.PrimarySkill),
		SecondarySkillID: int(deck.SecondarySkill),
		Version:          event.Version,
		SenderAddress:    event.SenderAddress,
		BlockHeight:      eventData.BlockHeight,
		IsDeleted:        false,
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

	cards, err := json.Marshal(event.Deck.Cards)
	if err != nil {
		return err
	}

	d := models.Deck{
		UserID:           event.UserId,
		DeckID:           deck.Id,
		Name:             deck.Name,
		HeroID:           deck.OverlordId,
		Cards:            cards,
		PrimarySkillID:   int(deck.PrimarySkill),
		SecondarySkillID: int(deck.SecondarySkill),
		Version:          event.Version,
		SenderAddress:    event.SenderAddress,
		BlockHeight:      eventData.BlockHeight,
		IsDeleted:        false,
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

	err := db.Model(&models.Deck{}).
		Where(&models.Deck{UserID: event.UserId, DeckID: event.DeckId}).
		Updates(models.Deck{IsDeleted: true}).
		Error
	if err != nil {
		return err
	}

	return nil
}

// TODO: seems this is not used anymore at all? can it be removed?
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
	match.BlockHeight = eventData.BlockHeight

	if err := db.Omit("created_at").Save(&match).Error; err != nil {
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
	var winnerID string
	if event.Block != nil && len(event.Block.List) > 0 && event.Block.List[0].GetEndGame() != nil {
		winnerID = event.Block.List[0].GetEndGame().WinnerId
	}

	var existingMatch models.Match
	notfound := db.Where(&models.Match{ID: match.Id}).
		First(&existingMatch).
		RecordNotFound()

	if !notfound {
		existingMatch.WinnerID = winnerID
		existingMatch.Status = match.Status.String()
		// since we just add blocktime, hav to check for compatibility
		if existingMatch.BlockTime.IsZero() {
			existingMatch.BlockTime = time.Unix(eventData.BlockTime, 0)
		}
		if err := db.Save(&existingMatch).Error; err != nil {
			return err
		}
	}

	dbReplay := models.Replay{}
	err = db.Where(&models.Replay{MatchID: match.Id}).First(&dbReplay).Error
	if err == nil {
		if replay != nil {
			dbReplay.ReplayJSON = replay
		}
		dbReplay.BlockHeight = eventData.BlockHeight
		dbReplay.BlockTime = time.Unix(eventData.BlockTime, 0)
		db.Save(&dbReplay)
	} else if gorm.IsRecordNotFoundError(err) {
		// insert
		dbReplay.MatchID = match.Id
		if replay != nil {
			dbReplay.ReplayJSON = replay
		}
		dbReplay.BlockHeight = eventData.BlockHeight
		dbReplay.BlockTime = time.Unix(eventData.BlockTime, 0)
		db.Create(&dbReplay)
	} else {
		return err
	}

	return nil
}

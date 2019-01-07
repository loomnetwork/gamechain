package cmd

import (
	"log"
	"strings"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/loomnetwork/go-loom/plugin/types"
	"github.com/loomnetwork/loomauth/models"
	"github.com/stretchr/testify/assert"
)

const event1 = `{
	"topics": [
		"zombiebattleground:createaccount"
	],
	"caller": {
		"chain_id": "default",
		"local": "PMEMWadk5kzr4wTe3SQjUSUjQ9Q="
	},
	"address": {
		"chain_id": "default",
		"local": "4ojW7scVDWoi/eM/CqLYHgZZHE0="
	},
	"plugin_name": "zombiebattleground:1.0.0",
	"block_height": 170961,
	"encoded_body": "eyJPd25lciI6ImxvY2s1IiwiTWV0aG9kIjoiY3JlYXRlYWNjb3VudCIsIkFkZHIiOiJQTUVNV2FkazVrenI0d1RlM1NRalVTVWpROVE9In0=",
	"original_request": "Cg1DcmVhdGVBY2NvdW50EhQKBWxvY2s1KgVJbWFnZUgBUgJ2MQ=="
	}`

const FindMatchEvent1 = `{"topics":["match:1","zombiebattleground:findmatch"],"caller":{"chain_id":"default","local":"xmOaJA2CIOecyHvkyJNo8WJjCkg="},"address":{"chain_id":"default","local":"4ojW7scVDWoi/eM/CqLYHgZZHE0="},"plugin_name":"zombiebattleground:1.0.0","block_height":18,"encoded_body":"ErsECAESB21hdGNoOjEa+QEKB3BsYXllcjYa7QEIARIHRGVmYXVsdBgEIgoKBlBvaXpvbRAEIgoKBkhhem1hehACIgoKBlplZXRlchABIg0KCVpoYXR0ZXJlchAEIgoKBlpsdWRnZRACIgsKB1pub3dtYW4QBCILCgdaLVZpcnVzEAEiCAoER2FlYRABIg8KC0VhcnRoc2hha2VyEAEiCgoGWmVwdGljEAEiCQoFR2hvdWwQASILCgdIYXp6YXJkEAEiDQoJRWN0b3BsYXNtEAEiCwoHQm91bmNlchACIgsKB1N0YXBsZXIQASIPCgtHb28gQm90dGxlcxABIgwKCENoYWluc2F3EAIa+QEKB3BsYXllcjca7QEIARIHRGVmYXVsdBgEIgoKBlBvaXpvbRAEIgoKBkhhem1hehACIgoKBlplZXRlchABIg0KCVpoYXR0ZXJlchAEIgoKBlpsdWRnZRACIgsKB1pub3dtYW4QBCILCgdaLVZpcnVzEAEiCAoER2FlYRABIg8KC0VhcnRoc2hha2VyEAEiCgoGWmVwdGljEAEiCQoFR2hvdWwQASILCgdIYXp6YXJkEAEiDQoJRWN0b3BsYXNtEAEiCwoHQm91bmNlchACIgsKB1N0YXBsZXIQASIPCgtHb28gQm90dGxlcxABIgwKCENoYWluc2F3EAIgASoCdjMwhom04QVAhom04QVSDwoHcGxheWVyNhCGibThBVIPCgdwbGF5ZXI3EIaJtOEFWgBaAA==","original_request":"CglGaW5kTWF0Y2gSCQoHcGxheWVyNg=="}`
const AcceptMatchEvent = `{"topics":["match:1","zombiebattleground:acceptmatch"],"caller":{"chain_id":"default","local":"xmOaJA2CIOecyHvkyJNo8WJjCkg="},"address":{"chain_id":"default","local":"4ojW7scVDWoi/eM/CqLYHgZZHE0="},"plugin_name":"zombiebattleground:1.0.0","block_height":20,"encoded_body":"Er0ECAESB21hdGNoOjEa+wEKB3BsYXllcjYQARrtAQgBEgdEZWZhdWx0GAQiCgoGUG9pem9tEAQiCgoGSGF6bWF6EAIiCgoGWmVldGVyEAEiDQoJWmhhdHRlcmVyEAQiCgoGWmx1ZGdlEAIiCwoHWm5vd21hbhAEIgsKB1otVmlydXMQASIICgRHYWVhEAEiDwoLRWFydGhzaGFrZXIQASIKCgZaZXB0aWMQASIJCgVHaG91bBABIgsKB0hhenphcmQQASINCglFY3RvcGxhc20QASILCgdCb3VuY2VyEAIiCwoHU3RhcGxlchABIg8KC0dvbyBCb3R0bGVzEAEiDAoIQ2hhaW5zYXcQAhr5AQoHcGxheWVyNxrtAQgBEgdEZWZhdWx0GAQiCgoGUG9pem9tEAQiCgoGSGF6bWF6EAIiCgoGWmVldGVyEAEiDQoJWmhhdHRlcmVyEAQiCgoGWmx1ZGdlEAIiCwoHWm5vd21hbhAEIgsKB1otVmlydXMQASIICgRHYWVhEAEiDwoLRWFydGhzaGFrZXIQASIKCgZaZXB0aWMQASIJCgVHaG91bBABIgsKB0hhenphcmQQASINCglFY3RvcGxhc20QASILCgdCb3VuY2VyEAIiCwoHU3RhcGxlchABIg8KC0dvbyBCb3R0bGVzEAEiDAoIQ2hhaW5zYXcQAiABKgJ2MzCGibThBUCGibThBVIPCgdwbGF5ZXI2EIaJtOEFUg8KB3BsYXllcjcQhom04QVaAFoA","original_request":"CgtBY2NlcHRNYXRjaBILCgdwbGF5ZXI2EAE="}`
const CreateDeckEvent = `{"topics":["zombiebattleground:createdeck"],"caller":{"chain_id":"default","local":"xmOaJA2CIOecyHvkyJNo8WJjCkg="},"address":{"chain_id":"default","local":"4ojW7scVDWoi/eM/CqLYHgZZHE0="},"plugin_name":"zombiebattleground:1.0.0","block_height":26,"encoded_body":"CgdwbGF5ZXIxEioweGM2NjM5YTI0MGQ4MjIwZTc5Y0M4N2JlNEM4OTM2OEYxNjI2MzBBNDgaKAgCEghOZXdEZWNrNxgBIgsKB1B5cm9tYXoQAiILCgdCdXJycm5uEAEiAnYz","original_request":"CgpDcmVhdGVEZWNrEjUKB3BsYXllcjESJhIITmV3RGVjazcYASILCgdQeXJvbWF6EAIiCwoHQnVycnJubhABIgJ2Mw=="}`
const EditDeckEvent = `{"topics":["zombiebattleground:editdeck"],"caller":{"chain_id":"default","local":"xmOaJA2CIOecyHvkyJNo8WJjCkg="},"address":{"chain_id":"default","local":"4ojW7scVDWoi/eM/CqLYHgZZHE0="},"plugin_name":"zombiebattleground:1.0.0","block_height":32,"encoded_body":"CgdwbGF5ZXIxEioweGM2NjM5YTI0MGQ4MjIwZTc5Y0M4N2JlNEM4OTM2OEYxNjI2MzBBNDgaHwgCEgxOZXdEZWNrN2FzZGYYASILCgdQeXJvbWF6EAMiAnYz","original_request":"CghFZGl0RGVjaxIuCgdwbGF5ZXIxEh8IAhIMTmV3RGVjazdhc2RmGAEiCwoHUHlyb21hehADIgJ2Mw=="}`
const DeleteDeckEvent = `{"topics":["zombiebattleground:deletedeck"],"caller":{"chain_id":"default","local":"xmOaJA2CIOecyHvkyJNo8WJjCkg="},"address":{"chain_id":"default","local":"4ojW7scVDWoi/eM/CqLYHgZZHE0="},"plugin_name":"zombiebattleground:1.0.0","block_height":34,"encoded_body":"CgdwbGF5ZXIxEioweGM2NjM5YTI0MGQ4MjIwZTc5Y0M4N2JlNEM4OTM2OEYxNjI2MzBBNDgYAg==","original_request":"CgpEZWxldGVEZWNrEgsKB3BsYXllcjEQAg=="}`

func ConnectDB(dbName string) *gorm.DB {
	db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	if err != nil {
		log.Printf("Got error when connect database, the error is '%v'", err)
	}
	return db
}

func DropDB() {
	db := ConnectDB("mysql")
	defer db.Close()
	db.Exec("drop database if exists gamechain_logger_test;")
}

func CreateDB() {
	db := ConnectDB("mysql")
	defer db.Close()
	db.Exec("create database gamechain_logger_test;")
}

func InitDB() *gorm.DB {
	CreateDB()
	db := ConnectDB("gamechain_logger_test")

	db.LogMode(true)

	if err := db.AutoMigrate(&models.Match{}).Error; err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&models.Deck{}).Error; err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&models.DeckCard{}).Error; err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&models.Replay{}).Error; err != nil {
		panic(err)
	}

	return db
}
func TestEventHandlers(t *testing.T) {
	DropDB()
	db := InitDB()

	var eventData types.EventData
	var unmarshaler jsonpb.Unmarshaler

	// FindMatch event
	err := unmarshaler.Unmarshal(strings.NewReader(FindMatchEvent1), &eventData)
	assert.Nil(t, err)

	err = FindMatchHandler(&eventData, db)
	assert.Nil(t, err)

	match := models.Match{}
	err = db.First(&match, 1).Error
	assert.Nil(t, err)
	assert.Equal(t, "player6", match.Player1ID)
	assert.Equal(t, "player7", match.Player2ID)
	assert.Equal(t, false, match.Player1Accepted)
	assert.Equal(t, false, match.Player2Accepted)

	// AcceptMatch event
	err = unmarshaler.Unmarshal(strings.NewReader(AcceptMatchEvent), &eventData)
	assert.Nil(t, err)

	err = AcceptMatchHandler(&eventData, db)
	assert.Nil(t, err)

	err = db.First(&match, 1).Error
	assert.Nil(t, err)
	assert.Equal(t, "player6", match.Player1ID)
	assert.Equal(t, "player7", match.Player2ID)
	assert.Equal(t, true, match.Player1Accepted)
	assert.Equal(t, false, match.Player2Accepted)

	// CreateDeck event
	err = unmarshaler.Unmarshal(strings.NewReader(CreateDeckEvent), &eventData)
	assert.Nil(t, err)

	err = CreateDeckHandler(&eventData, db)
	assert.Nil(t, err)

	deck := models.Deck{}
	err = db.First(&deck).Error
	assert.Nil(t, err)
	assert.Equal(t, "NewDeck7", deck.Name)
	assert.Equal(t, "player1", deck.UserID)

	// EditDeck event
	err = unmarshaler.Unmarshal(strings.NewReader(EditDeckEvent), &eventData)
	assert.Nil(t, err)

	err = EditDeckHandler(&eventData, db)
	assert.Nil(t, err)

	err = db.First(&deck, 1).Error
	assert.Nil(t, err)
	assert.Equal(t, "NewDeck7asdf", deck.Name)
	assert.Equal(t, "player1", deck.UserID)

	// DeleteDeck event
	err = unmarshaler.Unmarshal(strings.NewReader(DeleteDeckEvent), &eventData)
	assert.Nil(t, err)

	err = DeleteDeckHandler(&eventData, db)
	assert.Nil(t, err)

	db.Close()
	DropDB()
}

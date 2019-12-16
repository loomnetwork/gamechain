package models

import (
	"encoding/json"
	"time"
)

type Deck struct {
	ID               int64           `json:"id" gorm:"PRIMARY_KEY"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	UserID           string          `json:"user_id" gorm:"INDEX"`
	DeckID           int64           `json:"deck_id" gorm:"INDEX"`
	Name             string          `json:"name" gorm:"INDEX"`
	HeroID           int64           `json:"hero_id"`
	Cards            json.RawMessage `json:"cards" sql:"type:mediumtext;"`
	PrimarySkillID   int             `json:"primary_skill_id"`
	SecondarySkillID int             `json:"secondary_skill_id"`
	Version          string          `json:"version"`
	SenderAddress    string          `json:"sender_address"`
	BlockHeight      uint64          `json:"block_height" gorm:"INDEX"`
	IsDeleted        bool            `json:"is_deleted"`
}

type Match struct {
	ID              int64     `json:"id" gorm:"PRIMARY_KEY,auto_increment:false"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Player1ID       string    `json:"player1_id"`
	Player2ID       string    `json:"player2_id"`
	Player1Accepted bool      `json:"player1_accepted"`
	Player2Accepted bool      `json:"player2_accepted"`
	Player1DeckID   int64     `json:"player1_deck_id"`
	Player2DeckID   int64     `json:"player2_deck_id"`
	Status          string    `json:"status" gorm:"index:status"`
	Version         string    `json:"version"`
	RandomSeed      int64     `json:"random_seed"`
	WinnerID        string    `json:"winner_id"`
	BlockHeight     uint64    `json:"block_height" gorm:"INDEX"`
	BlockTime       time.Time `json:"block_time"`
	Turns           int       `json:"turns"`
	CardPlays       int       `json:"card_plays"`
	CardAttacks     int       `json:"card_attacks"`
}

type Replay struct {
	ID          int64     `gorm:"PRIMARY_KEY" json:"id"`
	MatchID     int64     `json:"match_id"`
	ReplayJSON  []byte    `json:"replay_json" sql:"type:mediumtext;"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	BlockHeight uint64    `json:"block_height"`
	BlockTime   time.Time `json:"block_time"`
}

type Card struct {
	ID          int64  `json:"id" gorm:"PRIMARY_KEY"`
	MouldID     string `json:"mould_id" gorm:"UNIQUE_INDEX:idx_mouldid_version"`
	Version     string `json:"version" gorm:"UNIQUE_INDEX:idx_mouldid_version"`
	Kind        string `json:"kind"`
	Set         string `json:"set"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FlavorText  string `json:"flavor_text"`
	Picture     string `json:"picture"`
	Rank        string `json:"rank"`
	Type        string `json:"type"`
	Rarity      string `json:"rarity"`
	Frame       string `json:"frame"`
	Damage      int32  `json:"damage"`
	Health      int32  `json:"health"`
	Cost        int32  `json:"cost"`
	Ability     string `json:"ability"`
	BlockHeight uint64 `json:"block_height"`
	ImageURL    string `json:"image_url"`
}

type ZbHeightCheck struct {
	Key             int64     `gorm:"not null;UNIQUE_INDEX:height_key;default=1" json:"key"`
	LastBlockHeight uint64    `json:"last_block_height"`
	UpdatedAt       time.Time `json:"updated_at"`
}
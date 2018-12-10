package models

import "time"

type Deck struct {
	ID               int64  `gorm:"PRIMARY_KEY"`
	UserID           string `gorm:"UNIQUE_INDEX:idx_userid_deckid"`
	DeckID           int64  `gorm:"UNIQUE_INDEX:idx_userid_deckid"`
	Name             string
	HeroID           int64
	Cards            []DeckCard
	PrimarySkillID   int
	SecondarySkillID int
	Version          string
	SenderAddress    string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type DeckCard struct {
	ID        int64 `gorm:"PRIMARY_KEY"`
	DeckID    uint
	CardName  string
	Amount    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

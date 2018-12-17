package models

import "time"

type Deck struct {
	ID               int64 `gorm:"PRIMARY_KEY"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	UserID           string `gorm:"UNIQUE_INDEX:idx_userid_deckid"`
	DeckID           int64  `gorm:"UNIQUE_INDEX:idx_userid_deckid"`
	Name             string
	HeroID           int64
	Cards            []DeckCard
	PrimarySkillID   int
	SecondarySkillID int
	Version          string
	SenderAddress    string
}

type DeckCard struct {
	ID        int64 `gorm:"PRIMARY_KEY"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    string `gorm:"type:varchar(255);index:userId"`
	DeckID    uint   `gorm:"index:deckId"`
	CardName  string
	Amount    int64
}

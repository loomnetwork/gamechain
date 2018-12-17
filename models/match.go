package models

import "time"

type Match struct {
	ID              int64 `gorm:"PRIMARY_KEY,auto_increment:false"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Player1ID       string
	Player2ID       string
	Player1Accepted bool
	Player2Accepted bool
	Player1DeckID   int64
	Player2DeckID   int64
	Status          string
	Version         string
	RandomSeed      int64
	Replay          Replay
	Deck            Deck
	WinnerID        string
}

type Replay struct {
	ID         int64 `gorm:"PRIMARY_KEY"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	MatchID    int64
	ReplayJSON []byte `sql:"type:mediumtext;"`
}

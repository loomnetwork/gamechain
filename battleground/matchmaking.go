package battleground

import (
	"sort"
	"time"

	"github.com/loomnetwork/gamechain/types/zb"
)

const (
	// MMRetries defines how many times the player keep retrying on match making
	MMRetries = 3
	// MMWaitTime defines how long the player will wait if there is no other player in player pool
	MMWaitTime = 3000 * time.Millisecond
	// MMTimeout determines how long the player should be in the player pool
	MMTimeout = 120 * time.Second
)

// MatchMakingFunc calculates the score based on the given profile target and candidate
type MatchMakingFunc func(target *zb.PlayerProfile, candidate *zb.PlayerProfile) float64

// mmf is the gobal match making function that calculates the match making score
var mmf MatchMakingFunc = func(target *zb.PlayerProfile, candidate *zb.PlayerProfile) float64 {
	return 1
}

// PlayerScore simply maintains the player id and score tuple
type PlayerScore struct {
	score float64
	id    string
}

type byPlayersScore []*PlayerScore

func (p byPlayersScore) Len() int { return len(p) }

func (p byPlayersScore) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p byPlayersScore) Less(i, j int) bool {
	return p[i].score > p[j].score
}

func sortByPlayerScore(ps []*PlayerScore) []*PlayerScore {
	sort.Sort(byPlayersScore(ps))
	return ps
}

func findPlayerProfileByID(pool *zb.PlayerPool, id string) *zb.PlayerProfile {
	for _, pp := range pool.PlayerProfiles {
		if pp.UserId == id {
			return pp
		}
	}
	return nil
}

func removePlayerFromPool(pool *zb.PlayerPool, id string) *zb.PlayerPool {
	var newpool zb.PlayerPool
	for _, pp := range pool.PlayerProfiles {
		if pp.UserId != id {
			newpool.PlayerProfiles = append(newpool.PlayerProfiles, pp)
		}
	}
	return &newpool
}

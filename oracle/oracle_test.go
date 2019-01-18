package oracle

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRecentHashPool(t *testing.T) {
	recentHashPool := newRecentHashPool(4 * time.Second)
	recentHashPool.startCleanupRoutine()

	require.True(t, recentHashPool.addHash([]byte{1, 2, 3}), "adding hash for first time should succed")

	require.False(t, recentHashPool.addHash([]byte{1, 2, 3}), "adding duplicate hash shouldnt be allowed")

	time.Sleep(5 * time.Second)

	require.True(t, recentHashPool.addHash([]byte{1, 2, 3}), "after timeout, hash should be allowed")
}

func TestTransferGatewayOracleMainnetEventSort(t *testing.T) {
	events := []*plasmachainEventInfo{
		&plasmachainEventInfo{BlockNum: 5, TxIdx: 0},
		&plasmachainEventInfo{BlockNum: 5, TxIdx: 1},
		&plasmachainEventInfo{BlockNum: 5, TxIdx: 4},
		&plasmachainEventInfo{BlockNum: 3, TxIdx: 3},
		&plasmachainEventInfo{BlockNum: 3, TxIdx: 7},
		&plasmachainEventInfo{BlockNum: 3, TxIdx: 1},
		&plasmachainEventInfo{BlockNum: 8, TxIdx: 4},
		&plasmachainEventInfo{BlockNum: 8, TxIdx: 1},
		&plasmachainEventInfo{BlockNum: 9, TxIdx: 0},
		&plasmachainEventInfo{BlockNum: 10, TxIdx: 5},
		&plasmachainEventInfo{BlockNum: 1, TxIdx: 2},
	}
	sortedEvents := []*plasmachainEventInfo{
		&plasmachainEventInfo{BlockNum: 1, TxIdx: 2},
		&plasmachainEventInfo{BlockNum: 3, TxIdx: 1},
		&plasmachainEventInfo{BlockNum: 3, TxIdx: 3},
		&plasmachainEventInfo{BlockNum: 3, TxIdx: 7},
		&plasmachainEventInfo{BlockNum: 5, TxIdx: 0},
		&plasmachainEventInfo{BlockNum: 5, TxIdx: 1},
		&plasmachainEventInfo{BlockNum: 5, TxIdx: 4},
		&plasmachainEventInfo{BlockNum: 8, TxIdx: 1},
		&plasmachainEventInfo{BlockNum: 8, TxIdx: 4},
		&plasmachainEventInfo{BlockNum: 9, TxIdx: 0},
		&plasmachainEventInfo{BlockNum: 10, TxIdx: 5},
	}
	sortPlasmachainEvents(events)
	require.EqualValues(t, sortedEvents, events, "wrong sort order")
}

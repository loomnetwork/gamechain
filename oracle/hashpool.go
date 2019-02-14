package oracle

import (
	"encoding/hex"
	"sync"
	"time"
)

type recentHashPool struct {
	hashMap         map[string]bool
	cleanupInterval time.Duration
	ticker          *time.Ticker
	stopCh          chan struct{}

	accessMutex sync.RWMutex
}

func newRecentHashPool(cleanupInterval time.Duration) *recentHashPool {
	return &recentHashPool{
		hashMap:         make(map[string]bool),
		cleanupInterval: cleanupInterval,
	}
}

func (r *recentHashPool) addHash(hash []byte) bool {
	r.accessMutex.Lock()
	defer r.accessMutex.Unlock()

	hexEncodedHash := hex.EncodeToString(hash)

	if _, ok := r.hashMap[hexEncodedHash]; ok {
		// If we are returning false, this means we have already seen hash
		return false
	}

	r.hashMap[hexEncodedHash] = true
	return true
}

func (r *recentHashPool) seenHash(hash []byte) bool {
	r.accessMutex.RLock()
	defer r.accessMutex.RUnlock()

	hexEncodedHash := hex.EncodeToString(hash)

	_, ok := r.hashMap[hexEncodedHash]
	return ok
}

func (r *recentHashPool) startCleanupRoutine() {
	r.ticker = time.NewTicker(r.cleanupInterval)
	r.stopCh = make(chan struct{})

	go func() {
		for {
			select {
			case <-r.stopCh:
				return
			case <-r.ticker.C:
				r.accessMutex.Lock()
				r.hashMap = make(map[string]bool)
				r.accessMutex.Unlock()
				break
			}
		}
	}()

}

func (r *recentHashPool) stopCleanupRoutine() {
	close(r.stopCh)
	r.ticker.Stop()
}

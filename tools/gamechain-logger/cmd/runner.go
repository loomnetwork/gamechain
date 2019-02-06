package cmd

import (
	"log"
	"strings"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"
	"github.com/loomnetwork/gamechain/battleground"
	"github.com/loomnetwork/go-loom/plugin/types"
	"github.com/loomnetwork/loomauth/models"
	"github.com/pkg/errors"
)

type Runner struct {
	db                *gorm.DB
	stopC             chan struct{}
	errC              chan error
	URL               string
	URLType           string
	reconnectInterval time.Duration
	blockInterval     int
	contractName      string
}

func NewRunner(URL string, db *gorm.DB, reconnectInterval time.Duration, blockInterval int, contractName string) *Runner {
	return &Runner{
		URL:               URL,
		db:                db,
		stopC:             make(chan struct{}),
		errC:              make(chan error),
		reconnectInterval: reconnectInterval,
		blockInterval:     blockInterval,
		contractName:      contractName,
	}
}

// Start runs the loop to watch topic. It's a blocking call.
func (r *Runner) Start() {
	for {
		err := r.watchTopic()
		if err == nil {
			break
		}
		log.Printf("error: %v", err)
		raven.CaptureErrorAndWait(err, map[string]string{})
		// delay before connecting again
		time.Sleep(r.reconnectInterval)
	}
}

func (r *Runner) Stop() {
	close(r.stopC)
}

func (r *Runner) Error() chan error {
	return r.errC
}

func (r *Runner) watchTopic() error {
	ticker := time.NewTicker(r.reconnectInterval)
	for {
		select {
		case <-ticker.C:
			height := models.ZbHeightCheck{}
			err := r.db.Where(&models.ZbHeightCheck{Key: 1}).First(&height).Error
			if err != nil && !gorm.IsRecordNotFoundError(err) {
				return err
			}
			fromBlock := height.LastBlockHeight + 1
			toBlock := fromBlock + uint64(r.blockInterval) - 1

			lastBlockHeight, err := queryBlockHeight(r.URL, r.contractName)
			if err != nil {
				return err
			}

			if toBlock > lastBlockHeight {
				continue
			}
			result, err := queryEventStore(r.URL, fromBlock, toBlock, r.contractName)
			if err != nil {
				return err
			}
			if err := r.batchProcessEvents(result.Events); err != nil {
				return err
			}
			if err = updateBlockHeight(r.db, result.ToBlock); err != nil {
				return err
			}
		case <-r.stopC:
			ticker.Stop()
			return nil
		}

	}
}

func (r *Runner) batchProcessEvents(events []*types.EventData) error {
	if len(events) == 0 {
		return nil
	}
	// need to create transaction to make sure all the data goes into db
	tx := r.db.Begin()
	for _, e := range events {
		for _, topic := range e.Topics {
			var topicHandler TopicHandler
			switch topic {
			case battleground.TopicFindMatchEvent:
				topicHandler = FindMatchHandler
			case battleground.TopicAcceptMatchEvent:
				topicHandler = AcceptMatchHandler
			case battleground.TopicCreateDeckEvent:
				topicHandler = CreateDeckHandler
			case battleground.TopicEditDeckEvent:
				topicHandler = EditDeckHandler
			case battleground.TopicDeleteDeckEvent:
				topicHandler = DeleteDeckHandler
			default:
				if strings.HasPrefix(topic, "match:") {
					topicHandler = MatchHandler
				}
			}

			if topicHandler != nil {
				err := topicHandler(e, tx)
				if err != nil {
					tx.Rollback()
					err = errors.Wrapf(err, "error calling topic handler")
					log.Printf("error: %s from event: %+v", err, e)
					return err
				}
			}
		}
	}
	tx.Commit()
	return nil
}

func updateBlockHeight(db *gorm.DB, blockHeight uint64) error {
	query := db.Model(&models.ZbHeightCheck{}).Where(&models.ZbHeightCheck{Key: 1}).Update("last_block_height", blockHeight)

	err, rows := query.Error, query.RowsAffected
	if err != nil {
		return err
	}
	if rows < 1 {
		db.Save(&models.ZbHeightCheck{Key: 1, LastBlockHeight: blockHeight})
	}

	return nil
}

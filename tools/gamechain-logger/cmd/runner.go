package cmd

import (
	"log"
	"strings"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"
	"github.com/loomnetwork/gamechain/battleground"
	"github.com/loomnetwork/go-loom/client"
	"github.com/loomnetwork/go-loom/plugin/types"
	"github.com/loomnetwork/gamechain/tools/gamechain-logger/modles"
	"github.com/pkg/errors"
)

type Config struct {
	ChainID           string
	ReadURI           string
	WriteURI          string
	ReconnectInterval time.Duration
	PollInterval      time.Duration
	BlockInterval     int
	ContractName      string
}

type Runner struct {
	db              *gorm.DB
	stopC           chan struct{}
	errC            chan error
	cfg             *Config
	dappchainClient *client.DAppChainRPCClient
}

func NewRunner(db *gorm.DB, config *Config) *Runner {
	return &Runner{
		db:              db,
		cfg:             config,
		stopC:           make(chan struct{}),
		errC:            make(chan error),
		dappchainClient: client.NewDAppChainRPCClient(config.ChainID, config.WriteURI, config.ReadURI),
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
		time.Sleep(r.cfg.ReconnectInterval)
	}
}

func (r *Runner) Stop() {
	close(r.stopC)
}

func (r *Runner) Error() chan error {
	return r.errC
}

func (r *Runner) watchTopic() error {
	ticker := time.NewTicker(r.cfg.PollInterval)
	for {
		select {
		case <-ticker.C:
			height := models.ZbHeightCheck{}
			err := r.db.Where(&models.ZbHeightCheck{Key: 1}).
				Attrs(models.ZbHeightCheck{Key: 1}).
				FirstOrCreate(&height).
				Error
			if err != nil {
				return err
			}

			fromBlock := height.LastBlockHeight + 1
			toBlock := fromBlock + uint64(r.cfg.BlockInterval) - 1

			lastBlockHeight, err := r.dappchainClient.GetBlockHeight()
			if err != nil {
				return err
			}

			if toBlock > lastBlockHeight {
				continue
			}
			result, err := r.dappchainClient.GetContractEvents(fromBlock, toBlock, r.cfg.ContractName)
			if err != nil {
				return err
			}

			tx := r.db.Begin()
			if err := r.batchProcessEvents(tx, result.Events); err != nil {
				tx.Rollback()
				return err
			}

			height.LastBlockHeight = result.ToBlock
			err = tx.Model(&models.ZbHeightCheck{}).
				Update(height).
				Error
			if err != nil {
				tx.Rollback()
				return err
			}

			tx.Commit()
		case <-r.stopC:
			ticker.Stop()
			return nil
		}

	}
}

func (r *Runner) batchProcessEvents(db *gorm.DB, events []*types.EventData) error {
	if len(events) == 0 {
		return nil
	}
	for _, event := range events {
		for _, topic := range event.Topics {
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
				err := topicHandler(event, db)
				if err != nil {
					err = errors.Wrapf(err, "error calling topic handler")
					log.Printf("error: %s from event: %+v", err, event)
					return err
				}
			}
		}
	}

	return nil
}

package main

import (
	"context"
	"strings"
	"sync"

	"github.com/polygon-io/go-app-ticker-wall/models"
	"github.com/sirupsen/logrus"
	tombv2 "gopkg.in/tomb.v2"
)

// TickerWallLeader manages this leaders state.
type TickerWallLeader struct {
	sync.RWMutex

	// config
	cfg *ServiceConfig

	// This keeps the settings
	clusterConfig *models.ScreenCluster

	// Our list of tickers we want to display / keep updated.
	Tickers []*models.Ticker

	// Our list of client screens currently connected.
	ScreenClients []*ScreenClient

	// Used for internally passing messages between websockets and parser.
	tickerUpdate chan []byte
}

type ScreenClient struct {
	UUID    string
	Updates chan *models.Update
	Screen  *models.Screen
	Stream  models.TickerWallLeader_RegisterAndListenForUpdatesServer
}

// NewTickerWallLeader creates a new ticker wall leader.
func NewTickerWallLeader(cfg *ServiceConfig) *TickerWallLeader {
	return &TickerWallLeader{
		cfg:          cfg,
		tickerUpdate: make(chan []byte, 1000), // Buffered channel to account for bursts.
		clusterConfig: &models.ScreenCluster{
			TickerBoxWidth: int32(cfg.TickerBoxWidthPx),
			ScrollSpeed:    int32(cfg.ScrollSpeed),
		},
	}
}

func (t *TickerWallLeader) Run(ctx context.Context) error {
	for _, ticker := range strings.Split(t.cfg.TickerList, ",") {
		// Make sure context hasn't closed.
		if ctx.Err() != nil {
			return ctx.Err()
		}

		newTickerObj, err := t.loadInitialTickerData(ctx, ticker)
		if err != nil {
			return err
		}

		t.Tickers = append(t.Tickers, newTickerObj)
		logrus.WithFields(logrus.Fields{
			"ticker":        ticker,
			"price":         newTickerObj.Price,
			"previousClose": newTickerObj.PreviousClosePrice,
		}).Debug("Added ticker..")
	}

	logrus.Info("Ticker data loaded..")

	// Create new tomb for this process.
	tomb, ctx := tombv2.WithContext(ctx)

	tomb.Go(func() error {
		return t.listenForTickerUpdates(ctx)
	})

	tomb.Go(func() error {
		return t.queueTickerUpdates(ctx)
	})

	tomb.Go(func() error {
		return t.runHTTPServer(ctx)
	})

	return tomb.Wait()
}

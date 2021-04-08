package main

import (
	"context"
	"fmt"
	"io"

	"github.com/polygon-io/go-app-ticker-wall/models"
	tickerManager "github.com/polygon-io/go-app-ticker-wall/ticker_manager"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// TickerWallClient manages the connection / interaction with leader.
type TickerWallClient struct {
	cfg    *ServiceConfig
	conn   *grpc.ClientConn
	client models.TickerWallLeaderClient
	screen *models.Screen

	manager tickerManager.TickerManager
}

// NewTickerWallClient creates a new ticker wall client.
func NewTickerWallClient(cfg *ServiceConfig, manager tickerManager.TickerManager) *TickerWallClient {
	obj := &TickerWallClient{
		cfg:     cfg,
		manager: manager,
		screen: &models.Screen{
			Width:              int32(cfg.ScreenWidth),
			Height:             int32(cfg.ScreenHeight),
			Index:              int32(cfg.ScreenIndex),
			ScreenGlobalOffset: 0,
		},
	}
	return obj
}

// Run starts all of our go routines / listeners.
func (t *TickerWallClient) Run(ctx context.Context) error {
	updateListener, err := t.client.RegisterAndListenForUpdates(ctx, t.screen)
	if err != nil {
		return err
	}

	for {
		// Read message.
		update, err := updateListener.Recv()
		if err == io.EOF {
			break
		}

		switch update.UpdateType {
		case models.UpdateTypeScreenCluster:
			t.updateScreenCluster(update)
		case models.UpdateTypeScreenTicker:
			// t.updateTicker(update)
		default:
			logrus.WithField("updateType", update.UpdateType).Warning("Unknown update type message.")
		}
	}

	return nil
}

func (t *TickerWallClient) updateScreenCluster(update *models.Update) error {
	logrus.Debug("Updating screen cluster information..")
	return nil
}

func (t *TickerWallClient) LoadTickers(ctx context.Context) error {
	tickers, err := t.client.GetTickers(ctx, t.screen)
	if err != nil {
		return err
	}

	// Add our tickers to the manager.
	for _, ticker := range tickers.Tickers {
		t.manager.AddTicker(ticker.Ticker, ticker.Price, (1 - (ticker.Price / ticker.PreviousClosePrice)), ticker.Ticker)
	}

	return nil
}

// Close cleans up our current grpc connection.
func (t *TickerWallClient) Close() error {
	return t.conn.Close()
}

// CreateGRPCClient creates a new GRPC client connection.
func (t *TickerWallClient) CreateGRPCClient() error {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMessageSize)))

	conn, err := grpc.Dial(cfg.Leader, opts...)
	if err != nil {
		return fmt.Errorf("not able to connect to grpc ticker wall leader: %w", err)
	}

	// Set our attributes.
	t.conn = conn
	t.client = models.NewTickerWallLeaderClient(t.conn)

	return nil
}

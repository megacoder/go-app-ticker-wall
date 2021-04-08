package main

import (
	"context"
	"fmt"

	"github.com/polygon-io/go-app-ticker-wall/models"
	"github.com/sirupsen/logrus"
)

func (t *TickerWallLeader) RegisterAndListenForUpdates(screen *models.Screen, stream models.TickerWallLeader_RegisterAndListenForUpdatesServer) error {
	logrus.Info("Got Screen: ", screen.Index)
	screenClient := &ScreenClient{
		Screen:  screen,
		Stream:  stream,
		Updates: make(chan *models.Update, 10), // dont block.
	}

	// Add new screen client.
	if err := t.addScreenToCluster(screenClient); err != nil {
		return fmt.Errorf("unable to add new screen client: %w", err)
	}

	// Remove this screen when we close the request.
	defer func() {
		t.removeScreenFromCluster(screenClient) // When we disconnect, remove from cluster.
	}()

	errChan := make(chan error, 2)
	go func() {
		// Start update listener until channel is closed.
		for update := range screenClient.Updates {
			if err := screenClient.Stream.Send(update); err != nil {
				errChan <- err
				return
			}
		}
	}()

	return nil
}

// GetTickers returns our current state of ticker data.
func (t *TickerWallLeader) GetTickers(ctx context.Context, screen *models.Screen) (*models.Tickers, error) {
	t.RLock()
	defer t.RUnlock()

	return &models.Tickers{
		Tickers: t.Tickers,
	}, nil
}

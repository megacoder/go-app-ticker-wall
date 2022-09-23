package client

import (
	"context"
	"errors"
	"sort"

	"github.com/megacoder/go-app-ticker-wall/models"
)

// tickerPriceUpdate handles updating a tickers price & market cap.
func (t *ClusterClient) tickerPriceUpdate(update *models.PriceUpdate) error {
	t.Lock()
	defer t.Unlock()

	for _, t := range t.Tickers {
		if t.Ticker == update.Ticker {
			t.Price = update.Price
			t.MarketCap = float64(t.OutstandingShares) * update.Price
			t.PriceChangePercentage = ((t.Price / t.PreviousClosePrice) - 1) * 100
		}
	}

	return nil
}

// tickerAdded handles adding a ticker to our local state.
func (t *ClusterClient) tickerAdded(ticker *models.Ticker) error {
	t.Lock()

	// Try to update what we have, if we have it.
	didUpdate := false
	for i, tick := range t.Tickers {
		if tick.Ticker == ticker.Ticker {
			t.Tickers[i] = ticker
			didUpdate = true
			break
		}
	}

	// If we didn't have it already, add it.
	if !didUpdate {
		t.Tickers = append(t.Tickers, ticker)
	}

	t.Unlock()

	err := t.tickerPriceUpdate(&models.PriceUpdate{
		Ticker: ticker.Ticker,
		Price:  ticker.Price,
	})
	if err != nil {
		return err
	}

	t.sortAndTagTickers()

	return nil
}

// tickerAdded handles removing a ticker from our local state.
func (t *ClusterClient) tickerRemoved(ticker *models.Ticker) error {
	t.Lock()

	// Find index of the given ticker.
	tickerIndex := -1
	for i, tick := range t.Tickers {
		if tick.Ticker == ticker.Ticker {
			tickerIndex = i
		}
	}

	// We didn't find this ticker??
	if tickerIndex == -1 {
		t.Unlock()
		return errors.New("unable to find ticker when attempting to remove it")
	}

	// Remove the element from the slice.
	t.Tickers[tickerIndex] = t.Tickers[len(t.Tickers)-1]
	t.Tickers[len(t.Tickers)-1] = nil
	t.Tickers = t.Tickers[:len(t.Tickers)-1]

	t.Unlock()

	t.sortAndTagTickers()

	return nil
}

// LoadTickers requests the full list of tickers from leader.
func (t *ClusterClient) LoadTickers(ctx context.Context) error {
	// Request full list of tickers from the leader.
	tickers, err := t.client.GetTickers(ctx, &models.Empty{})
	if err != nil {
		return err
	}

	for _, ticker := range tickers.Tickers {
		if err := t.tickerAdded(ticker); err != nil {
			return err
		}
	}

	return nil
}

func (t *ClusterClient) sortAndTagTickers() {
	t.Lock()
	defer t.Unlock()

	// Sort tickers (asc).
	sort.Sort(models.TickerSlice(t.Tickers))

	// Tag each ticker with it's Index.
	for i, ticker := range t.Tickers {
		ticker.Index = int32(i)
	}
}

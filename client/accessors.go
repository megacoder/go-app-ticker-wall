package client

import "github.com/megacoder/go-app-ticker-wall/models"

// GetTickers returns all the tickers we have.
func (t *ClusterClient) GetTickers() []*models.Ticker {
	t.RLock()
	defer t.RUnlock()

	return t.Tickers
}

// GetSettings returns the presentation settings.
func (t *ClusterClient) GetSettings() *models.PresentationSettings {
	t.RLock()
	defer t.RUnlock()

	return t.Cluster.Settings
}

// GetCluster returns the entire screen cluster.
func (t *ClusterClient) GetCluster() *models.ScreenCluster {
	t.RLock()
	defer t.RUnlock()

	return t.Cluster
}

// GetScreen returns the local screen settings.
func (t *ClusterClient) GetScreen() *models.Screen {
	t.RLock()
	defer t.RUnlock()

	return t.Screen
}

// GetAnnouncements returns the channel of announcements to be displayed.
func (t *ClusterClient) GetAnnouncements() chan *models.Announcement {
	t.RLock()
	defer t.RUnlock()

	return t.Announcements
}

// GetStatus returns the clients status.
func (t *ClusterClient) GetStatus() *Status {
	t.RLock()
	defer t.RUnlock()

	return t.Status
}

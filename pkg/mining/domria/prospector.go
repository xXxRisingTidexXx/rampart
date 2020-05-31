package domria

import (
	"rampart/pkg/mining"
)

func NewProspector() mining.Prospector {
	return &prospector{newFetcher()}
}

type prospector struct {
	fetcher *fetcher
}

func (prospector *prospector) Prospect(housing mining.Housing) error {
	_, err := prospector.fetcher.fetchFlats(housing)
	return err
}

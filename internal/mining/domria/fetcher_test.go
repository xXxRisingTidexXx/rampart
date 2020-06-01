package domria

import (
	"rampart/internal/mining"
	"rampart/internal/mining/configs"
	"testing"
)

func TestFetcherUnmarshalSearchEmptyString(t *testing.T) {
	fetcher := newFetcher("rampart-test/v1.0", &configs.Fetcher{})
	flats, err := fetcher.unmarshalSearch([]byte(""), mining.Primary)
	if flats != nil || err == nil {
		t.Error("domria: empty string unmarshalling didn't yield an error")
	}
}

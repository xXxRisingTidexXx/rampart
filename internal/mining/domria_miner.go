package mining

import (
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"net/http"
)

func NewDomriaMiner(config config.DomriaMiner) Miner {
	return &domriaMiner{
		config.Name,
		config.Spec,
		&http.Client{Timeout: config.Timeout},
		0,
		config.RetryLimit,
		config.SearchPrefix,
		config.UserAgent,
	}
}

type domriaMiner struct {
	name         string
	spec         string
	client       *http.Client
	page         int
	retryLimit   int
	searchPrefix string
	userAgent    string
}

func (m *domriaMiner) Name() string {
	return m.name
}

func (m *domriaMiner) Spec() string {
	return m.spec
}

func (m *domriaMiner) MineFlat() (Flat, error) {

	return Flat{}, nil
}

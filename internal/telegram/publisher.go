package telegram

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewPublisher(
	config config.Publisher,
	db *sql.DB,
	logger log.FieldLogger,
) (*Publisher, error) {
	return &Publisher{logger}, nil
}

type Publisher struct {
	logger log.FieldLogger
}

func (publisher *Publisher) Run() {

}

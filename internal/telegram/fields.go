package telegram

import (
	log "github.com/sirupsen/logrus"
)

func NewFields(handler string) Fields {
	return Fields{"handler": handler}
}

type Fields log.Fields

func (fields Fields) Handler() string {
	return fields["handler"].(string)
}

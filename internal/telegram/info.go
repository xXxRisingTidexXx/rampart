package telegram

import (
	log "github.com/sirupsen/logrus"
)

func NewInfo(handler string) Info {
	return Info{handler, log.Fields{}}
}

type Info struct {
	Handler string
	Extras  log.Fields
}

func (i Info) Fields() log.Fields {
	fields := make(log.Fields, len(i.Extras)+1)
	for key, value := range i.Extras {
		fields[key] = value
	}
	fields["handler"] = i.Handler
	return fields
}

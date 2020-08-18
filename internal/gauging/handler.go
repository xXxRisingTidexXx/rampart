package gauging

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func newHandler(gauger *Gauger) *handler {
	return &handler{gauger}
}

type handler struct {
	gauger *Gauger
}

func (handler *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		log.Errorf("httpserve: received invalid request method, %s", request.Method)
	} else {
		log.Info("gauging: success")
	}
}

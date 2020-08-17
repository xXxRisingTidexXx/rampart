package httpserve

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func newHandler() *handler {
	return &handler{}
}

type handler struct{}

func (handler *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if url := request.URL.String(); url != "/" {
		writer.WriteHeader(http.StatusNotFound)
		log.Errorf("httpserve: received invalid request url, %s", url)
	} else if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		log.Errorf("httpserve: received invalid request method, %s", request.Method)
	} else {
		
	}
}

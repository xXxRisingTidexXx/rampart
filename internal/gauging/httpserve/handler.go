package httpserve

import (
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
)

func newHandler() *handler {
	return &handler{misc.Headers{"Content-Type": "application/json"}}
}

type handler struct {
	headers misc.Headers
}

func (handler *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler.headers.Inject()
}

package gauging

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
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
		log.Errorf("gauging: received invalid request method, %s", request.Method)
		return
	}
	locations := make([]*dto.Location, 0)
	if err := json.NewDecoder(request.Body).Decode(&locations); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Errorf("gauging: failed to unmarshal the locations, %v", err)
		return
	}
	handler.gauger.GaugeAmenities(locations)
}

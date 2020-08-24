package gauging

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"net/http"
)

type handler struct {
	scheduler *Scheduler
	logger    log.FieldLogger
}

func (handler *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		handler.logger.WithField("method", request.Method).Errorf(
			"gauging: handler received invalid request method",
		)
		return
	}
	flats := make([]*dto.Flat, 0)
	if err := json.NewDecoder(request.Body).Decode(&flats); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		handler.logger.Errorf("gauging: handler failed to unmarshal the flats, %v", err)
		return
	}
	go handler.scheduler.ScheduleFlats(flats)
}

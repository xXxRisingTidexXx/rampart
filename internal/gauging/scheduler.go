package gauging

import (
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

func NewScheduler() *Scheduler {
	scheduler := &Scheduler{misc.Set{"Київ": struct{}{}}}
	go scheduler.run()
	return scheduler
}

type Scheduler struct {
	subwayCities misc.Set
}

func (scheduler *Scheduler) run() {

}

func (scheduler *Scheduler) ScheduleFlats(flats []*dto.Flat) {
	for _, flat := range flats {
		if scheduler.subwayCities.Contains(flat.City) {

		}
	}
}

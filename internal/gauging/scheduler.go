package gauging

import (
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
	"time"
)

func NewScheduler() *Scheduler {
	client := &http.Client{Timeout: 35 * time.Second}
	scheduler := &Scheduler{
		gaugingChannel:               make(chan *intent, 600),
		updateChannel:                make(chan *intent, 600),
		period:                       time.Second,
		subwayCities:                 misc.Set{"Київ": struct{}{}},
		subwayStationDistanceGauger:  NewSubwayStationDistanceGauger(client),
		industrialZoneDistanceGauger: NewIndustrialZoneDistanceGauger(client),
		greenZoneDistanceGauger:      NewGreenZoneDistanceGauger(client),
	}
	go scheduler.runGauging()
	go scheduler.runUpdate()
	return scheduler
}

type Scheduler struct {
	gaugingChannel                chan *intent
	updateChannel                 chan *intent
	period                        time.Duration
	subwayCities                  misc.Set
	subwayStationDistanceGauger   Gauger
	industrialZoneDistanceGauger  Gauger
	greenZoneDistanceGauger       Gauger
	subwayStationDistanceUpdater  Updater
	industrialZoneDistanceUpdater Updater
	greenZoneDistanceUpdater      Updater
}

func (scheduler *Scheduler) runGauging() {
	ticker := time.NewTicker(scheduler.period)
	for intent := range scheduler.gaugingChannel {
		<-ticker.C
		go scheduler.gaugeFlat(intent)
	}
}

func (scheduler *Scheduler) gaugeFlat(i *intent) {
	scheduler.updateChannel <- &intent{i.flat, i.gauger.GaugeFlat(i.flat), i.gauger, i.updater}
}

func (scheduler *Scheduler) runUpdate() {
	for intent := range scheduler.updateChannel {
		intent.updater.UpdateFlat(intent.flat, intent.value)
	}
}

func (scheduler *Scheduler) ScheduleFlats(flats []*dto.Flat) {
	for _, flat := range flats {
		if scheduler.subwayCities.Contains(flat.City) {
			scheduler.gaugingChannel <- &intent{
				flat:    flat,
				gauger:  scheduler.subwayStationDistanceGauger,
				updater: scheduler.subwayStationDistanceUpdater,
			}
		}
		scheduler.gaugingChannel <- &intent{
			flat:    flat,
			gauger:  scheduler.industrialZoneDistanceGauger,
			updater: scheduler.industrialZoneDistanceUpdater,
		}
		scheduler.gaugingChannel <- &intent{
			flat:    flat,
			gauger:  scheduler.greenZoneDistanceGauger,
			updater: scheduler.greenZoneDistanceUpdater,
		}
	}
}

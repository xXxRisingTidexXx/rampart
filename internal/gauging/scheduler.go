package gauging

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
	"time"
)

func NewScheduler(db *sql.DB) *Scheduler {
	client := &http.Client{Timeout: 35 * time.Second}
	scheduler := &Scheduler{
		make(chan *intent, 600),
		make(chan *intent, 600),
		time.Second,
		misc.Set{"Київ": struct{}{}},
		NewSubwayStationDistanceGauger(client),
		NewIndustrialZoneDistanceGauger(client),
		NewGreenZoneDistanceGauger(client),
		NewSubwayStationDistanceUpdater(db),
		NewIndustrialZoneDistanceUpdater(db),
		NewGreenZoneDistanceUpdater(db),
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
	value, err := i.gauger.GaugeFlat(i.flat)
	if err != nil {
		log.WithField("url", i.flat.OriginURL).Error(err)
	}
	scheduler.updateChannel <- &intent{i.flat, value, i.gauger, i.updater}
}

func (scheduler *Scheduler) runUpdate() {
	for intent := range scheduler.updateChannel {
		if err := intent.updater.UpdateFlat(intent.flat, intent.value); err != nil {
			log.WithField("url", intent.flat.OriginURL).Error(err)
		}
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

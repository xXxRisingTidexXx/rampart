package gauging

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/gauging/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
	"time"
)

func NewScheduler(config *config.Scheduler, db *sql.DB, logger log.FieldLogger) *Scheduler {
	client := &http.Client{Timeout: config.Timeout}
	scheduler := &Scheduler{
		make(chan *intent, config.Capacity),
		make(chan *intent, config.Capacity),
		config.Period,
		config.SubwayCities,
		NewSubwayStationDistanceGauger(config.SubwayStationDistanceGauger, client),
		NewIndustrialZoneDistanceGauger(config.IndustrialZoneDistanceGauger, client),
		NewGreenZoneDistanceGauger(config.GreenZoneDistanceGauger, client),
		NewSubwayStationDistanceUpdater(db),
		NewIndustrialZoneDistanceUpdater(db),
		NewGreenZoneDistanceUpdater(db),
		logger,
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
	logger                        log.FieldLogger
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
		scheduler.logger.WithField("url", i.flat.OriginURL).Error(err)
	}
	scheduler.updateChannel <- &intent{i.flat, value, i.gauger, i.updater}
}

func (scheduler *Scheduler) runUpdate() {
	for intent := range scheduler.updateChannel {
		if err := intent.updater.UpdateFlat(intent.flat, intent.value); err != nil {
			scheduler.logger.WithField("url", intent.flat.OriginURL).Error(err)
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
		} else {
			metrics.SubwayStationDistance.WithLabelValues("subwayless").Inc()
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

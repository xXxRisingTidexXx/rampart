package gauging

import (
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
)

type Gauger interface {
	GaugeFlat(flat *dto.Flat) float64
}

package gauging

import (
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
)

type Updater interface {
	UpdateFlat(flat *dto.Flat, value float64)
}

package gauging

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
)

type intent struct {
	flat    *dto.Flat
	value   float64
	gauger  Gauger
	updater Updater
}

func (intent *intent) String() string {
	return fmt.Sprintf("{%v %.6f %v %v}", intent.flat, intent.value, intent.gauger, intent.updater)
}

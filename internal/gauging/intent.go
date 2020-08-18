package gauging

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
)

type intent struct {
	target target
	flat   *dto.Flat
}

func (intent *intent) String() string {
	return fmt.Sprintf("{%d %v}", intent.target, intent.flat)
}

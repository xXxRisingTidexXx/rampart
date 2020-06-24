package domria

import (
	"fmt"
	"time"
)

type origin struct {
	updateTime time.Time
	price      float64
}

func (origin *origin) String() string {
	return fmt.Sprintf("{%s %.1f}", origin.updateTime, origin.price)
}

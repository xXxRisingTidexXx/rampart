package domria

import (
	"fmt"
	"time"
)

type origin struct {
	updateTime time.Time
}

func (origin *origin) String() string {
	return fmt.Sprintf("{%s}", origin.updateTime)
}

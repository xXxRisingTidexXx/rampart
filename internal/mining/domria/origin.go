package domria

import (
	"fmt"
	"time"
)

type origin struct {
	id         int
	updateTime time.Time
}

func (origin *origin) String() string {
	return fmt.Sprintf("{%d %s}", origin.id, origin.updateTime)
}

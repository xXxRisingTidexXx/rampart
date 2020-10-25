package domria

import (
	"time"
)

type origin struct {
	id         int
	updateTime time.Time
	isFound    bool
}

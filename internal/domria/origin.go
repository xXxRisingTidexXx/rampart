package domria

import (
	"time"
)

type origin struct {
	id         int
	upsertTime time.Time
	isFound    bool
}

package metrics

import (
	"time"
)

type bucket struct {
	sum   float64
	count float64
}

func (b *bucket) span(start time.Time) {
	b.sum += time.Since(start).Seconds()
	b.count++
}

func (b *bucket) avg() float64 {
	if b.count == 0 {
		return 0
	}
	return b.sum / b.count
}

func (b *bucket) reset() {
	b.sum = 0
	b.count = 0
}

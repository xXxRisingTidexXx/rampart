package domria

import (
	"fmt"
	"time"
)

type similarity struct {
	id         int
	updateTime time.Time
	price      float64
}

func (similarity *similarity) String() string {
	return fmt.Sprintf("{%d %s %.1f}", similarity.id, similarity.updateTime, similarity.price)
}

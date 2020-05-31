package domria

import (
	"fmt"
)

type prices struct {
	USD price `json:"1"`
	EUR price `json:"2"`
	UAH price `json:"3"`
}

func (prices *prices) String() string {
	return fmt.Sprintf("{%.1f %.1f %.1f}", prices.USD, prices.EUR, prices.UAH)
}

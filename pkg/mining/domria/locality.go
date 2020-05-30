package domria

import (
	"fmt"
)

type locality struct {
	State   string   `json:"stateName"`
	City    string   `json:"cityName"`
	Payload *payload `json:"payload"`
}

func (locality *locality) String() string {
	return fmt.Sprintf("{%s %s %v}", locality.State, locality.City, locality.Payload)
}

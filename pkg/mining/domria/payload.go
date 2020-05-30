package domria

import (
	"fmt"
)

type payload struct {
	StateID int `json:"stateId"`
	CityID  int `json:"cityId"`
}

func (payload *payload) String() string {
	return fmt.Sprintf("{%d %d}", payload.StateID, payload.CityID)
}

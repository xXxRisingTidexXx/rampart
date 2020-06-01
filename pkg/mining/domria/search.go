package domria

import (
	"fmt"
)

type search struct {
	Items []*item `json:"items"`
}

func (search *search) String() string {
	return fmt.Sprintf("{%v}", search.Items)
}

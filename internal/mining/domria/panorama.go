package domria

import (
	"fmt"
)

type panorama struct {
	Img string `json:"img"`
}

func (panorama *panorama) String() string {
	return fmt.Sprintf("{%s}", panorama.Img)
}

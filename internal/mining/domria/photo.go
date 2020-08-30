package domria

import (
	"fmt"
)

type photo struct {
	File string `json:"file"`
}

func (photo *photo) String() string {
	return fmt.Sprintf("{%s}", photo.File)
}

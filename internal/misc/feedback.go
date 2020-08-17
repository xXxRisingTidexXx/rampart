package misc

import (
	"fmt"
)

type Feedback struct {
	Message string `json:"message"`
}

func (feedback *Feedback) String() string {
	return fmt.Sprintf("{%s}", feedback.Message)
}

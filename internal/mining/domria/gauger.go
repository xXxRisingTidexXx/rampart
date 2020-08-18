package domria

import (
	"net/http"
)

type Gauger struct {
	client *http.Client
}

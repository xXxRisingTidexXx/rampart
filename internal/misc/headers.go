package misc

import (
	"net/http"
)

// TODO: remove it and use explicit headers.
type Headers map[string]string

func (headers Headers) Inject(request *http.Request) {
	for key, value := range headers {
		request.Header.Set(key, value)
	}
}

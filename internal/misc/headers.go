package misc

import (
	"net/http"
)

type Headers map[string]string

func (headers Headers) Inject(request *http.Request) {
	for key, value := range headers {
		request.Header.Set(key, value)
	}
}

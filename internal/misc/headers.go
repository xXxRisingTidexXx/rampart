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

func (headers Headers) Apply(writer http.ResponseWriter)  {
	for key, value := range headers {
		writer.Header().Set(key, value)
	}
}

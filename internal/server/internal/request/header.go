package request

import "net/http"

func CheckContentType(types []string, headers http.Header) bool {
	for _, t := range types {
		if headers.Get("Content-Type") == t {
			return true
		}
	}
	return false
}

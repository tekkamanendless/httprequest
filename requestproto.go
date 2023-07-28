package httprequest

import (
	"net/http"
	"strings"
)

// Proto returns the desired protocol of the request, after taking into account
// proxies and such.
//
// First priority goes to the `Forwarded` header, if it has `proto` set.
// Second priority goes to the `X-Forwarded-Proto` header, if there is a value.
// Last priority goes to whether or not TLS is set on the request.
func Proto(request *http.Request) string {
	// If the "Forwarded" header is set, use that if it gives us anything.
	//
	// Information on the "Forwarded" header: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Forwarded
	if headerValue := request.Header.Get("Forwarded"); headerValue != "" {
		headerValue = strings.ToLower(headerValue)
		directives := strings.Split(headerValue, ";")
		for _, directive := range directives {
			directive = strings.TrimSpace(directive)
			multiples := strings.Split(directive, ",")
			for _, multiple := range multiples {
				multiple = strings.TrimSpace(multiple)
				parts := strings.SplitN(multiple, "=", 2)
				if parts[0] == "proto" && len(parts[1]) > 0 {
					return parts[1]
				}
			}
		}
	}

	// If the "X-Forwarded-Proto" header is set, use that if it gives us anything.
	//
	// Information on the "X-Forwarded-Proto" header: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-Proto
	if headerValue := request.Header.Get("X-Forwarded-Proto"); headerValue != "" {
		return headerValue
	}

	// Otherwise, just use whatever the request contains.
	if request.TLS != nil {
		return "https"
	}
	return "http"
}

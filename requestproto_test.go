package httpextra

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequestProto(t *testing.T) {
	rows := []struct {
		description string
		request     http.Request
		result      string
	}{
		{
			description: "No fancy headers",
			request:     http.Request{},
			result:      "http",
		},
		{
			description: "No fancy headers but with TLS",
			request: http.Request{
				TLS: &tls.ConnectionState{},
			},
			result: "https",
		},
		{
			description: "X-Forwarded-Proto set",
			request: http.Request{
				Header: http.Header{
					"X-Forwarded-Proto": []string{"x-proto"},
				},
			},
			result: "x-proto",
		},
		{
			description: "X-Forwarded-Proto empty",
			request: http.Request{
				Header: http.Header{
					"X-Forwarded-Proto": []string{""},
				},
			},
			result: "http",
		},
		{
			description: "Forwarded proto not set",
			request: http.Request{
				Header: http.Header{
					"Forwarded": []string{"for=somebody"},
				},
			},
			result: "http",
		},
		{
			description: "Forwarded proto empty",
			request: http.Request{
				Header: http.Header{
					"Forwarded": []string{"for=somebody;proto="},
				},
			},
			result: "http",
		},
		{
			description: "Forwarded",
			request: http.Request{
				Header: http.Header{
					"Forwarded": []string{"for=somebody;proto=f-proto"},
				},
			},
			result: "f-proto",
		},
		{
			description: "Forwarded has highest priority",
			request: http.Request{
				Header: http.Header{
					"Forwarded":         []string{"for=somebody;proto=f-proto"},
					"X-Forwarded-Proto": []string{"x-proto"},
				},
			},
			result: "f-proto",
		},
	}
	for rowIndex, row := range rows {
		t.Run(fmt.Sprintf("%d/%s", rowIndex, row.description), func(t *testing.T) {
			result := Proto(&row.request)
			assert.Equal(t, row.result, result)
		})
	}
}

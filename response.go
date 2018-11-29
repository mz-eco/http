package http

import (
	"net/http"
	"strings"

	"github.com/mz-eco/memoir"
)

type BodyReader interface {
	Bytes() []byte
}

type Response struct {
	Request *Request
	Status  int
	Header  http.Header
	Body    interface{}
}

func (m *Response) Bytes() []byte {

	x, replace := body(m.Body)

	if replace {
		m.Body = x
	}

	return x
}

func (m *Response) UI() memoir.Component {

	var (
		headers = func() memoir.KeyValue {

			kv := make(memoir.KeyValue)

			for name, value := range m.Header {
				kv[name] = strings.Join(value, ";")
			}

			return kv
		}
	)

	if m == nil {
		return memoir.Label("Response")
	}

	return memoir.Label(
		"Response",
		memoir.KeyValue{
			"StatusCode": m.Status,
			"Status":     http.StatusText(m.Status),
		},
		memoir.Label(
			"Header",
			headers()),
		memoir.Label(
			"Body",
			memoir.DataView(m.Body),
		),
	)

}

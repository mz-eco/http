package http

import (
	"net/http"
	"time"

	"github.com/mz-eco/memoir"
)

type Translate struct {
	Name     string
	Error    error
	Request  *Request
	Response *Response
	Created  time.Time
	Used     time.Duration
}

func (m *Translate) Summary() *Summary {

	var (
		u = m.Request.Url()
	)
	return &Summary{
		Method:     m.Request.Method,
		Host:       u.Host,
		Path:       u.Path,
		Error:      m.Error != nil,
		Status:     http.StatusText(m.Response.Status),
		StatusCode: m.Response.Status,
		Create:     m.Created,
		Used:       m.Used,
	}

}

func (m *Translate) message() string {
	if m.Error != nil {
		return m.Error.Error()
	}

	return ""
}

func (m *Translate) UI() memoir.Component {

	return memoir.NewDocument(
		memoir.DocHttpTranslate,
		m.Name,
		memoir.KeyValue{
			"Error":   m.Error != nil,
			"Message": m.message(),
			"Created": m.Created,
			"Used":    m.Used,
		},
		m.Request,
		m.Response,
	)
}

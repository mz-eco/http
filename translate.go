package http

import (
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

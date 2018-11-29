package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/mz-eco/memoir"

	"github.com/pkg/errors"
)

type Header map[string]interface{}
type Query map[string]interface{}

type Request struct {
	URL    string
	Method string
	Header Header
	Query  Query
	Body   interface{}
}

func (m *Request) GetRequest() (*Request, bool) {
	return m, true
}

func (m *Request) UI() memoir.Component {

	return memoir.Label(
		"Request",
		memoir.Label(
			"Basic",
			memoir.KeyValue{
				"Method": m.Method,
				"Host":   m.URL,
				"URL":    m.Url(),
			},
		),
		memoir.Label(
			"Query",
			memoir.KeyValue(m.Query)),
		memoir.Label(
			"Header",
			memoir.KeyValue(m.Header),
		),
		memoir.Label(
			"Body",
			memoir.DataView(m.Body),
		),
	)
}

func (m *Request) Bytes() []byte {
	bx, replace := body(m.Body)

	if replace {
		m.Body = bx
	}
	return bx
}

func (m *Request) getBody() io.Reader {
	return bytes.NewBuffer(m.Bytes())
}

func (m *Request) getHeaders() http.Header {

	var (
		headers = make(http.Header)
	)

	for name, value := range m.Header {

		switch v := value.(type) {
		case string:
			headers.Add(name, v)
		case []string:
			for _, x := range v {
				headers.Set(name, x)
			}

		default:
			panic("header only support string | []string")
		}
	}

	return headers

}

func (m *Request) makeQuery(values url.Values) string {

	for name, value := range m.Query {

		switch v := value.(type) {
		case string:
			values.Add(name, v)
		case []string:
			values[name] = v
		case int:
			values.Set(name, fmt.Sprintf("%d", value))
		default:
			panic("query only support string | []string | int")
		}
	}

	return values.Encode()

}

func (m *Request) Url() string {

	u, err := url.Parse(m.URL)

	if err != nil {
		panic(err.Error())
	}

	u.RawQuery = m.makeQuery(u.Query())

	return u.String()
}

func (m *Request) Do(c *Client) (rep *Response, err error) {

	var (
		check = func(e error) bool {

			if e != nil {
				rep = nil
				err = errors.WithMessagef(e, "when do for [%s]", m.URL)

				return true
			}

			return false
		}
	)

	if len(m.Method) == 0 {
		m.Method = "GET"
	}

	url := m.Url()

	ask, err := http.NewRequest(
		m.Method,
		url,
		m.getBody())

	ask.Header = m.getHeaders()

	if check(err) {
		return
	}

	ack, err := c.Do(ask)

	if check(err) {
		return
	}

	defer ack.Body.Close()

	bytes, err := ioutil.ReadAll(ack.Body)

	if check(err) {
		return
	}

	return &Response{
		Request: m,
		Status:  ack.StatusCode,
		Header:  ack.Header,
		Body:    bytes,
	}, nil

}

func (m *Request) SetURL(url string) *Request {

	m.URL = url

	return m
}

func Clone(src *http.Request) *Request {

	var (
		x = &Request{
			URL:    src.URL.String(),
			Method: src.Method,
			Body:   src.Body,
			Query:  make(Query),
			Header: make(Header),
		}
	)

	for name, values := range src.URL.Query() {
		x.Query[name] = values
	}

	for name, values := range src.Header {
		x.Header[name] = values
	}

	return x
}

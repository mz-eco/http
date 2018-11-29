package http

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type Client = http.Client

type Action int

var (
	c = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				var (
					d = net.Dialer{
						Timeout:   10 * time.Second,
						KeepAlive: 10 * time.Minute,
					}
				)

				return d.DialContext(ctx, network, addr)
			},
		},
	}
)

type IfStmt struct {
	stmt DoStmt
	cond func() bool
}

func (m IfStmt) GetRequest() (*Request, bool) {

	if m.cond() {
		return m.stmt.GetRequest()
	}

	return nil, false

}

func If(cond func() bool, stmt DoStmt) DoStmt {

	return IfStmt{
		stmt: stmt,
		cond: cond,
	}
}

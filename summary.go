package http

import "time"

type Summary struct {
	Index      int
	Method     string
	Host       string
	Path       string
	Error      bool
	Status     string
	StatusCode int
	Create     time.Time
	Used       time.Duration
}

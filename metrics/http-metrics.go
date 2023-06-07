package metrics

import "time"

type HttpMetrics struct {
	Handler    string
	Method     string
	StatusCode string
	StartedAt  time.Time
	FinishedAt time.Time
	Duration   float64
}

func NewHttpMetrics(
	handler string,
	method string,
) *HttpMetrics {
	return &HttpMetrics{
		Handler: handler,
		Method:  method,
	}
}

func (http *HttpMetrics) Started() {
	http.StartedAt = time.Now().UTC()
}

func (http *HttpMetrics) Finished() {
	http.FinishedAt = time.Now().UTC()
	http.Duration = time.Since(http.StartedAt).Seconds()
}

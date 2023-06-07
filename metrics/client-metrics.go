package metrics

import "time"

type ClientMetrics struct {
	Name       string
	StartedAt  time.Time
	FinishedAt time.Time
	Duration   float64
}

func NewClientMetrics(
	name string,
) *ClientMetrics {
	return &ClientMetrics{
		Name: name,
	}
}

func (client *ClientMetrics) Started() {
	client.StartedAt = time.Now().UTC()
}

func (client *ClientMetrics) Finished() {
	client.FinishedAt = time.Now().UTC()
	client.Duration = time.Since(client.StartedAt).Seconds()
}

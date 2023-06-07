package nats

import (
	"github.com/nats-io/nats.go"
)

type Publisher interface {
	Publish(subject string, data []byte) error
}

type publisher struct {
	js nats.JetStreamContext
}

func NewPublisher(
	js nats.JetStreamContext,
) *publisher {
	return &publisher{
		js: js,
	}
}

func (p *publisher) Publish(subject string, data []byte) error {
	_, err := p.js.Publish(subject, data)
	if err != nil {
		return err
	}

	return nil
}

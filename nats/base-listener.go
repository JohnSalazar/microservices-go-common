package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	ackWait       = 5 * time.Second
	maxAckPending = 25
)

type Listener interface {
	Listener(subject string, queueGroupName string, durableName string, handler nats.MsgHandler)
}

type listener struct {
	js nats.JetStreamContext
}

func NewListener(
	js nats.JetStreamContext,
) *listener {
	return &listener{
		js: js,
	}
}

func (l *listener) Listener(subject string, queueGroupName string, durableName string, handler nats.MsgHandler) {
	ticker := time.NewTicker(1500 * time.Millisecond)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				err := l.queueSubscribe(subject, queueGroupName, durableName, handler)
				if err == nil {
					<-quit
				}
				fmt.Println(fmt.Errorf("subject: %v, QueueSubscribe: %v, Error: %v", subject, queueGroupName, err))
				ticker.Reset(2000 * time.Millisecond)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (l *listener) queueSubscribe(subject string, queueGroupName string, durableName string, handler nats.MsgHandler) error {
	_, err := l.js.QueueSubscribe(
		subject,
		queueGroupName,
		handler,
		nats.ManualAck(),
		nats.AckWait(ackWait),
		nats.Durable(durableName),
		nats.MaxAckPending(maxAckPending),
		nats.DeliverAll(),
	)

	return err
}

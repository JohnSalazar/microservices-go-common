package nats

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/oceano-dev/microservices-go-common/config"
	common_service "github.com/oceano-dev/microservices-go-common/services"
	trace "github.com/oceano-dev/microservices-go-common/trace/otel"
	trace_span "go.opentelemetry.io/otel/trace"
)

type CommandErrorHelper struct {
	config *config.Config
	email  common_service.EmailService
}

func NewCommandErrorHelper(
	config *config.Config,
	email common_service.EmailService,
) *CommandErrorHelper {
	return &CommandErrorHelper{
		config: config,
		email:  email,
	}
}

func (c CommandErrorHelper) CheckUnmarshal(msg *nats.Msg, err error) error {
	if err != nil {
		log.Printf("error unmarshalling %s command: %v", msg.Subject, err)
		msgErr := fmt.Sprintf("appName: %s: error unmarshalling command: %s data: %s %s\n", c.config.AppName, msg.Subject, msg.Data, err.Error())

		go c.email.SendSupportMessage(msgErr)
	}

	return err
}

func (c CommandErrorHelper) CheckCommandError(span trace_span.Span, msg *nats.Msg, err error) {
	if err != nil {
		msgErr := fmt.Sprintf("appName: %s: error processing %s: data: %s %s\n", c.config.AppName, msg.Subject, msg.Data, err.Error())
		trace.FailSpan(span, msgErr)

		// errs := make(chan error, 1)
		// go func() {
		// 	errs <- c.email.SendSupportMessage(msgErr)
		// }()

		// if err := <-errs; err != nil {
		// 	log.Printf("error sending support message: %s", err.Error())
		// }

		go c.email.SendSupportMessage(msgErr)

		log.Println(msgErr)
		fmt.Printf("%s fail!!!\n", msg.Subject)
	} else {
		fmt.Printf("%s processed!!!\n", msg.Subject)
	}
}

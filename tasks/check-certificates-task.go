package tasks

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/security"
	"github.com/oceano-dev/microservices-go-common/services"
	trace "github.com/oceano-dev/microservices-go-common/trace/otel"
)

type CheckCertificatesTask struct {
	config              *config.Config
	managerCertificates security.ManagerCertificates
	emailService        services.EmailService
}

func NewCheckCertificatesTask(
	config *config.Config,
	managerCertificates security.ManagerCertificates,
	emailService services.EmailService,
) *CheckCertificatesTask {
	return &CheckCertificatesTask{
		config:              config,
		managerCertificates: managerCertificates,
		emailService:        emailService,
	}
}

func (task *CheckCertificatesTask) Start(ctx context.Context, certsDone chan bool) {
	ticker := time.NewTicker(2500 * time.Millisecond)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				_, span := trace.NewSpan(ctx, "CheckCertificatesTask.Start")
				defer span.End()

				certsIsValid := task.managerCertificates.VerifyCertificates()
				if !certsIsValid {
					caCertOK := task.checkCertificateCA()
					certHostOK := task.checkCertificateHost()

					if !caCertOK || !certHostOK {
						ticker.Reset(15 * time.Second)
						break
					}
				}
				fmt.Printf("start check certificates successfully: %s\n", time.Now().UTC())

				ticker.Reset(time.Duration(task.config.Certificates.MinutesToReloadCertificate) * time.Minute)
				certsDone <- true
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (task *CheckCertificatesTask) checkCertificateCA() bool {
	err := task.managerCertificates.GetCertificateCA()
	if err != nil {
		msg := fmt.Sprintln("EmailService - certificate CA error: ", err)
		err := task.emailService.SendSupportMessage(msg)
		if err != nil {
			log.Println(err)
		}
		log.Println(msg)
	}

	return err == nil
}

func (task *CheckCertificatesTask) checkCertificateHost() bool {
	err := task.managerCertificates.GetCertificate()
	if err != nil {
		msg := fmt.Sprintln("EmailService - certificate Host error: ", err)
		err := task.emailService.SendSupportMessage(msg)
		if err != nil {
			log.Println(err)
		}
		log.Println(msg)
	}

	return err == nil
}

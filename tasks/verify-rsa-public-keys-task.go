package tasks

import (
	"log"
	"time"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/security"
)

type VerifyRSAPublicKeysTask struct {
	config  *config.Config
	manager security.ManagerSecurityRSAKeys
}

func NewVerifyRSAPublicKeysTask(
	config *config.Config,
	manager security.ManagerSecurityRSAKeys,
) *VerifyRSAPublicKeysTask {
	return &VerifyRSAPublicKeysTask{
		config:  config,
		manager: manager,
	}
}

func (task *VerifyRSAPublicKeysTask) ReloadRSAPublicKeys() {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				keys := task.manager.GetAllRSAPublicKeys()

				if keys == nil {
					log.Printf("rsa public keys not success refreshed %s\n", time.Now().UTC())
					ticker.Reset(15 * time.Second)
					break
				}

				//fmt.Printf("rsa public keys success refreshed %s\n", time.Now().UTC())
				ticker.Reset(1 * time.Hour)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

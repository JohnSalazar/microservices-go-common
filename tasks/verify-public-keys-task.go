package tasks

import (
	"log"
	"time"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/security"
)

type VerifyPublicKeysTask struct {
	config  *config.Config
	manager security.ManagerSecurityKeys
}

func NewVerifyPublicKeysTask(
	config *config.Config,
	manager security.ManagerSecurityKeys,
) *VerifyPublicKeysTask {
	return &VerifyPublicKeysTask{
		config:  config,
		manager: manager,
	}
}

func (task *VerifyPublicKeysTask) ReloadPublicKeys() {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				keys := task.manager.GetAllPublicKeys()

				if keys == nil {
					log.Printf("public keys not success refreshed %s\n", time.Now().UTC())
					ticker.Reset(15 * time.Second)
					break
				}

				//fmt.Printf("public keys success refreshed %s\n", time.Now().UTC())
				ticker.Reset(1 * time.Hour)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

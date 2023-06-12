package security

import (
	"fmt"
	"sync"
	"time"

	"github.com/JohnSalazar/microservices-go-common/config"
	"github.com/JohnSalazar/microservices-go-common/models"
	"github.com/JohnSalazar/microservices-go-common/services"
)

type managerSecurityKeys struct {
	config  *config.Config
	service services.SecurityKeysService
}

var (
	muxKeys           sync.Mutex
	publicKeys        []*models.ECDSAPublicKey
	refreshPublicKeys = time.Now().UTC()
)

func NewManagerSecurityKeys(
	config *config.Config,
	service services.SecurityKeysService,
) *managerSecurityKeys {
	return &managerSecurityKeys{
		config:  config,
		service: service,
	}
}

func (m *managerSecurityKeys) GetAllPublicKeys() []*models.ECDSAPublicKey {
	if publicKeys == nil {
		m.refreshPublicKeys()
	}

	publicKeysRefresh := refreshPublicKeys.Before(time.Now().UTC())
	if publicKeysRefresh {
		m.refreshPublicKeys()
		fmt.Println("refresh public keys")
	}

	return publicKeys
}

func (m *managerSecurityKeys) refreshPublicKeys() {
	newestPublicKeys, err := m.service.GetAllPublicKeys()
	if err != nil {
		fmt.Println(err)
	}

	muxKeys.Lock()
	publicKeys = nil
	publicKeys = append(publicKeys, newestPublicKeys...)
	muxKeys.Unlock()

	refreshPublicKeys = time.Now().UTC().Add(time.Minute * time.Duration(m.config.SecurityKeys.MinutesToRefreshPublicKeys))
}

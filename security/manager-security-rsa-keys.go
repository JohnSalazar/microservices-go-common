package security

import (
	"crypto/rsa"
	"fmt"
	"sync"
	"time"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/models"
	"github.com/oceano-dev/microservices-go-common/services"
)

type managerSecurityRSAKeys struct {
	config  *config.Config
	service services.SecurityRSAKeysService
}

var (
	muxRSAKeys           sync.Mutex
	rsaPublicKeys        []*models.RSAPublicKey
	refreshRSAPublicKeys = time.Now().UTC()
)

func NewManagerSecurityRSAKeys(
	config *config.Config,
	service services.SecurityRSAKeysService,
) *managerSecurityRSAKeys {
	return &managerSecurityRSAKeys{
		config:  config,
		service: service,
	}
}

func (m *managerSecurityRSAKeys) GetAllRSAPublicKeys() []*models.RSAPublicKey {
	if publicKeys == nil {
		m.refreshRSAPublicKeys()
	}

	rsaPublicKeysRefresh := refreshRSAPublicKeys.Before(time.Now().UTC())
	if rsaPublicKeysRefresh {
		m.refreshRSAPublicKeys()
		fmt.Println("refresh RSA public keys")
	}

	return rsaPublicKeys
}

func (m *managerSecurityRSAKeys) Encrypt(msg string, publicKey *rsa.PublicKey) ([]byte, error) {
	return m.service.Encrypt(msg, publicKey)
}

func (m *managerSecurityRSAKeys) Decrypt(encryptedBytes []byte, privateKey *rsa.PrivateKey) (string, error) {
	return m.service.Decrypt(encryptedBytes, privateKey)
}

func (m *managerSecurityRSAKeys) refreshRSAPublicKeys() {
	newestRSAPublicKeys, err := m.service.GetAllRSAPublicKeys()
	if err != nil {
		fmt.Println(err)
	}

	muxRSAKeys.Lock()
	rsaPublicKeys = nil
	rsaPublicKeys = append(rsaPublicKeys, newestRSAPublicKeys...)
	muxRSAKeys.Unlock()

	refreshRSAPublicKeys = time.Now().UTC().Add(time.Minute * time.Duration(m.config.SecurityRSAKeys.MinutesToRefreshRSAPublicKeys))
}

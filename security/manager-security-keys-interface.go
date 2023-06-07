package security

import "github.com/oceano-dev/microservices-go-common/models"

type ManagerSecurityKeys interface {
	GetAllPublicKeys() []*models.ECDSAPublicKey
}

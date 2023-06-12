package security

import "github.com/JohnSalazar/microservices-go-common/models"

type ManagerSecurityKeys interface {
	GetAllPublicKeys() []*models.ECDSAPublicKey
}

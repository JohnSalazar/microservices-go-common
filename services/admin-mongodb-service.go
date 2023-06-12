package services

import (
	"context"
	"time"

	"github.com/JohnSalazar/microservices-go-common/config"
	"github.com/JohnSalazar/microservices-go-common/repositories"
)

type AdminMongoDbService struct {
	config                 *config.Config
	adminMongoDbRepository *repositories.AdminMongoDbRepository
}

func NewAdminMongoDbService(
	config *config.Config,
	adminMongoDbRepository *repositories.AdminMongoDbRepository,
) *AdminMongoDbService {
	return &AdminMongoDbService{
		config:                 config,
		adminMongoDbRepository: adminMongoDbRepository,
	}
}

func (service *AdminMongoDbService) VerifyMongoDBExporterUser() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()

	user := service.config.MongoDbExporter.User
	pwd := service.config.MongoDbExporter.Password

	userExisting, err := service.adminMongoDbRepository.FindMongoDBExporterUser(ctx, user)
	if err != nil {
		return false, err
	}

	if userExisting == user {
		return true, nil
	}

	err = service.adminMongoDbRepository.Create(ctx, user, pwd)
	if err != nil {
		return false, err
	}

	return true, nil
}

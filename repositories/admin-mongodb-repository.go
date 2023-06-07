package repositories

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminMongoDbRepository struct {
	database *mongo.Database
}

func NewAdminMongoDbRepository(
	database *mongo.Database,
) *AdminMongoDbRepository {
	return &AdminMongoDbRepository{
		database: database,
	}
}

func (r *AdminMongoDbRepository) FindMongoDBExporterUser(ctx context.Context, user string) (string, error) {
	cmd := bson.D{{Key: "usersInfo", Value: user}}

	databaseAdmin := r.database.Client().Database("admin")

	singleResult := databaseAdmin.RunCommand(ctx, cmd)
	if singleResult.Err() != nil {
		if singleResult.Err() != mongo.ErrNoDocuments {
			return "", singleResult.Err()
		}
	}

	var rawData bson.M
	err := singleResult.Decode(&rawData)
	if err != nil {
		fmt.Println(err)
	}

	userFound := r.getUser(rawData)

	return userFound, nil
}

func (r *AdminMongoDbRepository) Create(ctx context.Context, user string, pwd string) error {
	cmd := bson.D{
		{Key: "createUser", Value: user},
		{Key: "pwd", Value: pwd},
		{Key: "roles", Value: bson.A{
			bson.D{
				{Key: "role", Value: "clusterMonitor"},
				{Key: "db", Value: "admin"},
			},
			bson.D{
				{Key: "role", Value: "read"},
				{Key: "db", Value: "local"},
			},
		},
		},
	}

	databaseAdmin := r.database.Client().Database("admin")

	result := databaseAdmin.RunCommand(ctx, cmd)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (r *AdminMongoDbRepository) getUser(rawData map[string]interface{}) string {
	for _, value := range rawData {
		if pa, ok := value.(primitive.A); ok {
			valueMSI := []interface{}(pa)
			for _, v := range valueMSI {
				for key, provider := range v.(primitive.M) {
					if key == "user" {
						return provider.(string)
					}
				}
			}
		}
	}

	return ""
}

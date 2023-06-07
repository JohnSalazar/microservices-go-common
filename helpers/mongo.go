package helpers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StringToID(s string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(s)
	return id
}

func IsValidID(s string) bool {
	return primitive.IsValidObjectID(s)
}

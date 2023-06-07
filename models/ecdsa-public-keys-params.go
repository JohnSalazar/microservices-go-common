package models

import (
	"time"
)

type ECDSAPublicKeysParams struct {
	Alg       string            `bson:"alg" json:"alg"`
	Kid       string            `bson:"kid" json:"kid"`
	Use       string            `bson:"use" json:"use"`
	ExpiresAt time.Time         `bson:"expires_at" json:"expires_at"`
	Params    map[string]string `bson:"params" json:"params"`
}

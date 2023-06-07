package models

import (
	"crypto/rsa"
	"time"
)

type RSAPublicKey struct {
	Key       *rsa.PublicKey `json:"key"`
	Kid       string         `json:"kid"`
	ExpiresAt time.Time      `json:"expires_at"`
}

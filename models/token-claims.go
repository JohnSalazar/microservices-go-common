package models

import "github.com/golang-jwt/jwt/v4"

type TokenClaims struct {
	Sub    string   `json:"sub,omitempty"`   // ID usuário
	Email  string   `json:"email,omitempty"` // Email usuário
	Jti    string   `json:"jti,omitempty"`   // ID único para o token
	Nbf    int64    `json:"nbf,omitempty"`   // Especifica a partir de quando o token passa a ser válido
	Iat    int64    `json:"iat,omitempty"`   // Determina a idade do token
	Exp    int64    `json:"exp,omitempty"`   // Tempo para utilização do token
	Iss    string   `json:"iss,omitempty"`   // Quem emitiu o token
	Claims []Claims `json:"claims,omitempty"`
	jwt.RegisteredClaims
}

type Claims struct {
	Type  string `bson:"type" json:"type"`
	Value string `bson:"value" json:"value"`
}

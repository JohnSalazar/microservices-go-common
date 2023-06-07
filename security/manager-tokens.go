package security

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type ManagerTokens struct {
	config              *config.Config
	managerSecurityKeys ManagerSecurityKeys
}

func NewManagerTokens(
	config *config.Config,
	managerSecurityKeys ManagerSecurityKeys,
) *ManagerTokens {
	return &ManagerTokens{
		config:              config,
		managerSecurityKeys: managerSecurityKeys,
	}
}

func (m *ManagerTokens) ReadHeadAccessToken(c *gin.Context) (*models.TokenClaims, error) {
	var err error
	tokenString := c.Request.Header.Get("Authorization")
	if len(tokenString) == 0 {
		return nil, fmt.Errorf("token not found")
	}

	bearer := strings.Split(tokenString, " ")
	if bearer[0] != "Bearer" {
		return nil, fmt.Errorf("token is not Bearer")
	}

	tokenString = bearer[1]

	var keyFunc = m.getKeyFunc()
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, keyFunc)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	claims, ok := token.Claims.(*models.TokenClaims)
	expires := time.Now().UTC().After(time.Unix(claims.Exp, 0).UTC())
	if !ok || !token.Valid || token.Header["typ"] != "access" || claims.Iss != m.config.Token.Issuer || expires {
		return nil, errors.New("JWT failed validation")
	}

	return claims, nil
}

// func (m *ManagerTokens) ReadCookieAccessToken(c *gin.Context) (*models.TokenClaims, error) {

// 	log.Println("==========================INICIO REQUEST ==============================")
// 	log.Println(c.Request)
// 	log.Println("==========================FIM REQUEST ==============================")
// 	log.Println("")
// 	log.Println("==========================INICIO HEADER ==============================")
// 	log.Println(c.Request.Header)
// 	log.Println("==========================FIM HEADER ==============================")

// 	var err error
// 	tokenString, err := c.Cookie("accessToken")
// 	if err != nil {
// 		return nil, fmt.Errorf(err.Error())
// 	}

// 	var keyFunc = m.getKeyFunc()
// 	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, keyFunc)
// 	if err != nil {
// 		return nil, fmt.Errorf(err.Error())
// 	}

// 	claims, ok := token.Claims.(*models.TokenClaims)
// 	if !ok || !token.Valid || claims.Iss != m.config.Token.Issuer {
// 		return nil, errors.New("JWT failed validation")
// 	}

// 	return claims, nil
// }

func (m *ManagerTokens) ReadRefreshToken(c *gin.Context, tokenString string) (string, error) {
	var keyFunc = m.getKeyFunc()
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, keyFunc)
	if err != nil {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(*models.TokenClaims)
	expires := time.Now().UTC().After(time.Unix(claims.Exp, 0).UTC())
	if !ok || !token.Valid || token.Header["typ"] != "refresh" || claims.Iss != m.config.Token.Issuer || expires {
		return "", errors.New("invalid token")
	}

	return claims.Sub, nil
}

func (m *ManagerTokens) getKeyFunc() jwt.Keyfunc {
	keys := m.managerSecurityKeys.GetAllPublicKeys()

	var keyFunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("expecting JWT header to have string kid")
		}

		var key *ecdsa.PublicKey
		for i := range keys {
			if keys[i].Kid == keyID {
				key = keys[i].Key
				break
			}
		}

		if key == nil {
			return nil, fmt.Errorf("unable to parse public key")
		}

		return key, nil
	}

	return keyFunc
}

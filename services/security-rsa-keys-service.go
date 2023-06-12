package services

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/JohnSalazar/microservices-go-common/config"
	"github.com/JohnSalazar/microservices-go-common/models"
	"github.com/eapache/go-resiliency/breaker"
	"github.com/eapache/go-resiliency/retrier"
)

type SecurityRSAKeysService interface {
	GetAllRSAPublicKeys() ([]*models.RSAPublicKey, error)
	Encrypt(msg string, publicKey *rsa.PublicKey) ([]byte, error)
	Decrypt(encryptedBytes []byte, privateKey *rsa.PrivateKey) (string, error)
}

type securityRSAKeysService struct {
	config  *config.Config
	service CertificatesService
}

func NewSecurityRSAKeysService(
	config *config.Config,
	service CertificatesService,
) *securityRSAKeysService {
	return &securityRSAKeysService{
		config:  config,
		service: service,
	}
}

func (s *securityRSAKeysService) GetAllRSAPublicKeys() ([]*models.RSAPublicKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := s.requestRSAPublicKey(ctx)
	if err != nil {
		return nil, err
	}

	modelsPublicKeys := []*models.RSAPublicKey{}
	err = json.Unmarshal([]byte(data), &modelsPublicKeys)
	if err != nil {
		return nil, err
	}

	return modelsPublicKeys, nil
}

func (s *securityRSAKeysService) Encrypt(msg string, publicKey *rsa.PublicKey) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		[]byte(msg),
		nil)
	if err != nil {
		return nil, err
	}

	return encryptedBytes, nil
}

func (s *securityRSAKeysService) Decrypt(encryptedBytes []byte, privateKey *rsa.PrivateKey) (string, error) {
	decryptedBytes, err := privateKey.Decrypt(
		nil,
		encryptedBytes,
		&rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		panic(err)
	}

	return string(decryptedBytes), nil
}

func (s *securityRSAKeysService) requestRSAPublicKey(ctx context.Context) ([]byte, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true,
				GetCertificate: s.service.GetLocalCertificate,
				RootCAs:        s.service.GetLocalCertificateCA()},
		},
	}

	var err error
	request, err := http.NewRequestWithContext(ctx, "GET", s.config.SecurityRSAKeys.EndPointGetRSAPublicKeys, nil)
	if err != nil {
		log.Println("request error:", err)
		return nil, err
	}

	var response *http.Response
	r := retrier.New(retrier.ConstantBackoff(6, 10*time.Millisecond), nil)
	err = r.Run(func() error {
		b := breaker.New(6, 1, 5*time.Second)
		for {
			result := b.Run(func() error {
				response, err = client.Do(request)
				if err != nil {
					return err
				}

				return nil
			})

			switch result {
			case nil:
				return nil
			case breaker.ErrBreakerOpen:
				return err
			default:
				return err
			}
		}
	})

	if err != nil {
		log.Println("response error:", err)
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("data parse:", err)
		return nil, err
	}

	return data, nil
}

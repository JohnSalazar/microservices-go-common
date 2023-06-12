package services

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JohnSalazar/microservices-go-common/config"
	"github.com/JohnSalazar/microservices-go-common/helpers"
)

type CertificatesService interface {
	GetCertificateCA() ([]byte, error)
	GetCertificateHost() ([]byte, error)
	GetCertificateHostKey() ([]byte, error)
	GetPathsCertificateCAAndKey() (string, string)
	GetPathsCertificateHostAndKey() (string, string)
	ReadCertificateCA() ([]byte, error)
	ReadCertificate() (*x509.Certificate, error)
	GetLocalCertificateCA() *x509.CertPool
	GetLocalCertificate(info *tls.ClientHelloInfo) (*tls.Certificate, error)
}

type certificatesService struct {
	config *config.Config
}

func NewCertificatesService(
	config *config.Config,
) *certificatesService {
	return &certificatesService{
		config: config,
	}
}

func (s *certificatesService) GetCertificateCA() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := s.requestCertificateCA(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *certificatesService) GetCertificateHost() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := s.requestCertificateHost(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *certificatesService) GetCertificateHostKey() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := s.requestCertificateHostKey(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *certificatesService) GetPathsCertificateCAAndKey() (string, string) {
	caCertPath := fmt.Sprintf("%s/ca_%s", s.config.Certificates.FolderName, s.config.Certificates.FileNameCert)
	caKeyPath := fmt.Sprintf("%s/ca_%s", s.config.Certificates.FolderName, s.config.Certificates.FileNameKey)

	return caCertPath, caKeyPath
}

func (s *certificatesService) GetPathsCertificateHostAndKey() (string, string) {
	certPath := fmt.Sprintf("%s/%s", s.config.Certificates.FolderName, s.config.Certificates.FileNameCert)
	keyPath := fmt.Sprintf("%s/%s", s.config.Certificates.FolderName, s.config.Certificates.FileNameKey)

	return certPath, keyPath
}

func (s *certificatesService) GetLocalCertificateCA() *x509.CertPool {
	caCertPath, _ := s.GetPathsCertificateCAAndKey()
	if !helpers.FileExists(caCertPath) {
		fmt.Println("certificate CA not found")
		return nil
	}

	caCertBytes, err := s.ReadCertificateCA()
	if err != nil {
		return nil
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCertBytes)

	return pool
}

func (s *certificatesService) GetLocalCertificate(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	certPath, keyPath := s.GetPathsCertificateHostAndKey()
	if !helpers.FileExists(certPath) || !helpers.FileExists(keyPath) {
		return nil, errors.New("certificate host not found")
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &cert, nil
}

func (s *certificatesService) ReadCertificateCA() ([]byte, error) {
	caCertPath, _ := s.GetPathsCertificateHostAndKey()
	data, err := os.ReadFile(caCertPath)
	if err != nil {
		os.Exit(1)
		return nil, fmt.Errorf("read Certificate CA file error")
	}

	pemBlock, _ := pem.Decode(data)
	if pemBlock == nil {
		return nil, fmt.Errorf("decode Certificate CA error")
	}

	return pemBlock.Bytes, nil
}

func (s *certificatesService) ReadCertificate() (*x509.Certificate, error) {
	certPath, _ := s.GetPathsCertificateHostAndKey()
	data, err := os.ReadFile(certPath)
	if err != nil {
		os.Exit(1)
		return nil, fmt.Errorf("read Certificate file error")
	}

	pemBlock, _ := pem.Decode(data)
	if pemBlock == nil {
		return nil, fmt.Errorf("decode Certificate error")
	}

	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func (s *certificatesService) requestCertificateCA(ctx context.Context) ([]byte, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	hash := s.getHash()
	endPoint := fmt.Sprintf("%s/%s", s.config.Certificates.EndPointGetCertificateCA, hash)
	request, err := http.NewRequestWithContext(ctx, "GET", endPoint, nil)
	if err != nil {
		log.Println("request:", err)
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("response:", err)
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

func (s *certificatesService) requestCertificateHost(ctx context.Context) ([]byte, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	hash := s.getHash()
	endPoint := fmt.Sprintf("%s/%s", s.config.Certificates.EndPointGetCertificateHost, hash)
	request, err := http.NewRequestWithContext(ctx, "GET", endPoint, nil)
	if err != nil {
		log.Println("request:", err)
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("response:", err)
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

func (s *certificatesService) requestCertificateHostKey(ctx context.Context) ([]byte, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	hash := s.getHash()
	endPoint := fmt.Sprintf("%s/%s", s.config.Certificates.EndPointGetCertificateHostKey, hash)
	request, err := http.NewRequestWithContext(ctx, "GET", endPoint, nil)
	if err != nil {
		log.Println("request:", err)
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("response:", err)
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

func (s *certificatesService) getHash() string {
	return base64.StdEncoding.EncodeToString([]byte(s.config.Certificates.PasswordPermissionEndPoint))
}

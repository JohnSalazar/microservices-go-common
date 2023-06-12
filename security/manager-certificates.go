package security

import (
	"errors"
	"time"

	"github.com/JohnSalazar/microservices-go-common/config"
	"github.com/JohnSalazar/microservices-go-common/helpers"
	"github.com/JohnSalazar/microservices-go-common/services"

	"github.com/eapache/go-resiliency/breaker"
)

type managerCertificates struct {
	config  *config.Config
	service services.CertificatesService
}

var (
	caCertPath string
	certPath   string
	keyPath    string
)

func NewManagerCertificates(
	config *config.Config,
	service services.CertificatesService,
) *managerCertificates {
	caCertPath, _ = service.GetPathsCertificateCAAndKey()
	certPath, keyPath = service.GetPathsCertificateHostAndKey()
	return &managerCertificates{
		config:  config,
		service: service,
	}
}

func (m *managerCertificates) VerifyCertificates() bool {
	if helpers.FileExists(caCertPath) && helpers.FileExists(certPath) && helpers.FileExists(keyPath) {
		caCert, err := m.service.ReadCertificateCA()
		if caCert == nil || err != nil {
			return false
		}

		cert, err := m.service.ReadCertificate()
		if cert == nil || err != nil {
			return false
		}

		if cert == nil || cert.NotAfter.AddDate(0, 0, -7).Before(time.Now().UTC()) {
			return false
		}

		return true
	}

	return false
}

func (m *managerCertificates) GetCertificateCA() error {
	err := m.refreshCertificateCA()
	if err != nil {
		return err
	}

	return nil
}

func (m *managerCertificates) GetCertificate() error {
	err := m.refreshCertificate()
	if err != nil {
		return err
	}

	return nil
}

func (m *managerCertificates) refreshCertificateCA() error {
	err := m.requestCertificateCA()
	if err != nil {
		return err
	}

	return nil
}

func (m *managerCertificates) refreshCertificate() error {
	err := m.requestCertificate()
	if err != nil {
		return err
	}

	err = m.requestCertificateKey()
	if err != nil {
		return err
	}

	return nil
}

func (m managerCertificates) requestCertificateCA() error {
	b := breaker.New(3, 1, 5*time.Second)
	for {
		var caCert []byte
		var err error
		err = b.Run(func() error {
			caCert, err = m.service.GetCertificateCA()
			if err != nil {
				return err
			}

			return nil
		})

		switch err {
		case nil:
			if caCert == nil {
				return errors.New("certificate CA not found")
			}

			err := helpers.CreateFile(caCert, caCertPath)
			if err != nil {
				return err
			}

			return nil
		case breaker.ErrBreakerOpen:
			return err
		}
	}
}

func (m managerCertificates) requestCertificate() error {
	b := breaker.New(3, 1, 5*time.Second)
	for {
		var cert []byte
		var err error
		err = b.Run(func() error {
			cert, err = m.service.GetCertificateHost()
			if err != nil {
				return err
			}

			return nil
		})

		switch err {
		case nil:
			if cert == nil {
				return errors.New("certificate not found")
			}

			err := helpers.CreateFile(cert, certPath)
			if err != nil {
				return err
			}

			return nil
		case breaker.ErrBreakerOpen:
			return err
		}
	}
}

func (m *managerCertificates) requestCertificateKey() error {
	b := breaker.New(3, 1, 5*time.Second)
	for {
		var key []byte
		var err error
		err = b.Run(func() error {
			key, err = m.service.GetCertificateHostKey()
			if err != nil {
				return err
			}

			return nil
		})

		switch err {
		case nil:
			if key == nil {
				return errors.New("certificate key not found")
			}

			err := helpers.CreateFile(key, keyPath)
			if err != nil {
				return err
			}

			return nil
		case breaker.ErrBreakerOpen:
			return err
		}
	}
}

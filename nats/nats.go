package nats

import (
	"crypto/tls"
	"log"
	"time"

	"github.com/JohnSalazar/microservices-go-common/config"
	"github.com/JohnSalazar/microservices-go-common/services"
	"github.com/nats-io/nats.go"
)

func NewNats(config *config.Config, service services.CertificatesService) (*nats.Conn, error) {
	tls := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		InsecureSkipVerify: true,
		GetCertificate:     service.GetLocalCertificate,
		RootCAs:            service.GetLocalCertificateCA(),
	}
	nc, err := nats.Connect(
		config.Nats.Url,
		nats.Timeout(time.Second*time.Duration(config.Nats.ConnectWait)),
		nats.PingInterval(time.Second*time.Duration(config.Nats.Interval)),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(5),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Fatalf("Connection lost: %v", err)
		}),
		nats.Secure(tls),
	)

	return nc, err
}

func NewJetStream(nc *nats.Conn, streamName string, subjects []string) (nats.JetStreamContext, error) {
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	stream, _ := js.StreamInfo(streamName)
	if stream == nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: subjects,
		})
		if err != nil {
			return nil, err
		}
	} else {
		_, err = js.UpdateStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: subjects,
		})
		if err != nil {
			return nil, err
		}
	}

	return js, nil
}

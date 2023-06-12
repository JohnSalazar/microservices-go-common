package httputil

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/JohnSalazar/microservices-go-common/config"
	"github.com/JohnSalazar/microservices-go-common/services"
	"github.com/gin-gonic/gin"
)

type HttpServer interface {
	RunTLSServer() (*http.Server, error)
}

type httpServer struct {
	config  *config.Config
	router  *gin.Engine
	service services.CertificatesService
}

var srv *http.Server

func NewHttpServer(
	config *config.Config,
	router *gin.Engine,
	service services.CertificatesService,
) *httpServer {
	return &httpServer{
		config:  config,
		router:  router,
		service: service,
	}
}

func (s *httpServer) RunTLSServer() (*http.Server, error) {
	var err error
	if srv == nil {
		srv = s.mountTLSServer()

		go func() {
			if err = srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				log.Fatalf("err: %s\n", err)
			}
		}()

		log.Printf("Listening on port %s", s.config.ListenPort)
	}

	return srv, err
}

func (s *httpServer) mountTLSServer() *http.Server {
	return &http.Server{
		Addr:         s.config.ListenPort,
		Handler:      s.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		TLSConfig: &tls.Config{
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
			GetCertificate: s.service.GetLocalCertificate,
			ClientCAs:      s.service.GetLocalCertificateCA(),
		},
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
}

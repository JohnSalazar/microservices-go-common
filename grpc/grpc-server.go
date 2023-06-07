package proto

import (
	"crypto/tls"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	grpc_otel "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/middlewares"
	"github.com/oceano-dev/microservices-go-common/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type GrpcServer struct {
	config         *config.Config
	service        services.CertificatesService
	serviceMetrics services.Metrics
}

func NewGrpcServer(
	config *config.Config,
	service services.CertificatesService,
	serviceMetrics services.Metrics,
) *GrpcServer {
	return &GrpcServer{
		config:         config,
		service:        service,
		serviceMetrics: serviceMetrics,
	}
}

func (s *GrpcServer) CreateGrpcServer() (*grpc.Server, error) {
	grpcServer := grpc.NewServer(
		grpc.Creds(s.credentials()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(s.config.GrpcServer.MaxConnectionIdle) * time.Minute,
			Timeout:           time.Duration(s.config.GrpcServer.Timeout) * time.Second,
			MaxConnectionAge:  time.Duration(s.config.GrpcServer.MaxConnectionAge) * time.Minute,
			Time:              time.Duration(s.config.GrpcServer.Timeout) * time.Minute,
		}),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_otel.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_recovery.StreamServerInterceptor(),
			middlewares.StreamServerInterceptorMetrics(s.serviceMetrics),
		),
		),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_otel.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
			middlewares.UnaryServerInterceptorMetrics(s.serviceMetrics)),
		),
	)

	return grpcServer, nil
}

func (s *GrpcServer) credentials() credentials.TransportCredentials {
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
		GetCertificate: s.service.GetLocalCertificate,
		ClientCAs:      s.service.GetLocalCertificateCA(),
	}

	return credentials.NewTLS(tls)
}

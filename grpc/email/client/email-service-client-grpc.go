package proto

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"github.com/JohnSalazar/microservices-go-common/config"
	"github.com/JohnSalazar/microservices-go-common/services"
	trace "github.com/JohnSalazar/microservices-go-common/trace/otel"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type EmailServiceClientGrpc struct {
	config  *config.Config
	service services.CertificatesService
}

func NewEmailServiceClientGrpc(
	config *config.Config,
	service services.CertificatesService,
) *EmailServiceClientGrpc {
	return &EmailServiceClientGrpc{
		config:  config,
		service: service,
	}
}

var grpcClient EmailServiceClient

func (s *EmailServiceClientGrpc) SendPasswordCode(email string, code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	ctx, span := trace.NewSpan(ctx, "emailServiceGrpc.SendPasswordCodeReq")
	defer span.End()

	s.verifyClientGrpc(ctx)

	req := &PasswordCodeReq{
		Email: email,
		Code:  code,
	}

	validator := validator.New()
	if err := validator.StructCtx(ctx, req); err != nil {
		trace.AddSpanError(span, err)
		log.Printf("emailServiceGrpc.SendPasswordCodeReq: %v", err)
		return err
	}

	_, err := grpcClient.SendPasswordCode(ctx, req)
	if err != nil {
		return err
	}

	log.Print("email sent")

	return nil
}

func (s *EmailServiceClientGrpc) SendSupportMessage(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	ctx, span := trace.NewSpan(ctx, "emailServiceGrpc.SendSupportMessageReq")
	defer span.End()

	s.verifyClientGrpc(ctx)

	req := &SupportMessageReq{
		Message: message,
	}

	validator := validator.New()
	if err := validator.StructCtx(ctx, req); err != nil {
		trace.AddSpanError(span, err)
		log.Printf("emailServiceGrpc.SendSupportMessageReq: %v", err)
		return err
	}

	_, err := grpcClient.SendSupportMessage(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *EmailServiceClientGrpc) verifyClientGrpc(ctx context.Context) {
	if grpcClient == nil {
		s.createClientGrpc(ctx)
	}
}

func (s *EmailServiceClientGrpc) createClientGrpc(ctx context.Context) {
	conn, err := grpc.DialContext(
		ctx,
		s.config.EmailService.Host,
		grpc.WithTransportCredentials(s.credentials()),
		grpc.WithBlock())
	if err != nil {
		log.Fatalf("EmailServiceClientGrpc error connection: %v", err)
	}

	grpcClient = NewEmailServiceClient(conn)
}

func (s *EmailServiceClientGrpc) credentials() credentials.TransportCredentials {
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
		GetCertificate:     s.service.GetLocalCertificate,
		RootCAs:            s.service.GetLocalCertificateCA(),
	}

	return credentials.NewTLS(tls)
}

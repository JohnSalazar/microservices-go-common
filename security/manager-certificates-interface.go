package security

type ManagerCertificates interface {
	VerifyCertificates() bool
	GetCertificateCA() error
	GetCertificate() error
}

package services

type EmailService interface {
	SendPasswordCode(email string, code string) error
	SendSupportMessage(message string) error
}

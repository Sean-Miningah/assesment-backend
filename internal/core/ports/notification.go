package ports

import "context"

type NotificationRepository interface {
	SendSms(ctx context.Context, phoneNumber []string, senderID, message string) error
	SendEmail(ctx context.Context, email, subject, to string) error
}

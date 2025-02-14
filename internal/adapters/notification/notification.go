package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sean-miningah/sil-backend-assessment/pkg/utils"
	"go.opentelemetry.io/otel"
	gomail "gopkg.in/mail.v2"
)

type NotificationRepo struct {
	ATAPIKey             string
	NotificationUsername string
	GmailAppAPIKey       string
	ATAPIUrl             string
}

func NewNotificationRepo(ataPIKey, notificationUsername, gmailAppAPIKey, atAPIUrl string) *NotificationRepo {
	return &NotificationRepo{
		ATAPIKey:             ataPIKey,
		NotificationUsername: notificationUsername,
		GmailAppAPIKey:       gmailAppAPIKey,
		ATAPIUrl:             atAPIUrl,
	}
}

func (s *NotificationRepo) SendSms(ctx context.Context, phoneNumber []string, senderID, message string) error {
	payload := utils.SMSRequest{
		Username:     s.NotificationUsername,
		Message:      message,
		SenderID:     senderID,
		PhoneNumbers: phoneNumber,
	}
	if err := s.sendRequest(ctx, payload); err != nil {
		return err
	}
	return nil
}

func (s *NotificationRepo) sendRequest(ctx context.Context, payload utils.SMSRequest) error {
	ctx, span := otel.Tracer("").Start(ctx, "NotificationService.sendATSMS")
	defer span.End()

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", s.ATAPIUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", s.ATAPIUrl)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API responded with status: %d", resp.StatusCode)
	}

	fmt.Println("SMS sent successfully")
	return nil
}

func (s *NotificationRepo) SendEmail(ctx context.Context, email, subject, to string) error {
	ctx, span := otel.Tracer("").Start(ctx, "NotificationService.sendMail")
	defer span.End()
	m := gomail.NewMessage()
	m.SetHeader("From", "seanpminingah@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", email)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "seanpminingah@gmail.com", s.GmailAppAPIKey)

	if err := dialer.DialAndSend(m); err != nil {
		fmt.Println("Error", err)
		return err
	} else {
		fmt.Println("Email sent successfully")
	}

	return nil
}

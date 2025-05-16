package services

import (
	"log"

	"github.com/resend/resend-go/v2"
)

type EmailService struct {
	client      *resend.Client
	senderEmail string
}

func NewEmailService(apiKey, senderEmail string) *EmailService {
	return &EmailService{
		client:      resend.NewClient(apiKey),
		senderEmail: senderEmail,
	}
}

func (s *EmailService) SendArrivalNotification(to, studentName, time string) error {
	params := &resend.SendEmailRequest{
		From:    s.senderEmail,
		To:      []string{to},
		Subject: studentName + " has arrived at school",
		Html: `<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #2563eb;">School Attendance Notification</h2>
			<p>Dear Parent,</p>
			<p>Your child <strong>` + studentName + `</strong> has arrived at school at <strong>` + time + `</strong>.</p>
			<p>Thank you,</p>
			<p>The School Administration</p>
		</div>`,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("Failed to send arrival notification: %v", err)
		return err
	}

	return nil
}

func (s *EmailService) SendAbsenceNotification(to, studentName, date string) error {
	params := &resend.SendEmailRequest{
		From:    s.senderEmail,
		To:      []string{to},
		Subject: studentName + " was absent from school",
		Html: `<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #dc2626;">Absence Notification</h2>
			<p>Dear Parent,</p>
			<p>Your child <strong>` + studentName + `</strong> was marked absent on <strong>` + date + `</strong>.</p>
			<p>If this is incorrect, please contact the school office.</p>
			<p>Thank you,</p>
			<p>The School Administration</p>
		</div>`,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("Failed to send absence notification: %v", err)
		return err
	}

	return nil
}

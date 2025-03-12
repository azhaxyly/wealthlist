package service

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"
	"wealthlist/config"
	"wealthlist/internal/models"
)

type FeedbackService struct {
	smtpConfig config.SMTPConfig
	log        *slog.Logger
}

func NewFeedbackService(cfg *config.Config, log *slog.Logger) *FeedbackService {
	return &FeedbackService{
		smtpConfig: cfg.SMTP,
		log:        log,
	}
}

func (s *FeedbackService) SendFeedbackEmail(feedback models.FeedbackDto) error {
	s.log.Info("Sending feedback email", slog.String("from", feedback.Email))

	messageBody := s.createEmailBody(feedback)

	smtpAuth := smtp.PlainAuth("", s.smtpConfig.Username, s.smtpConfig.Password, s.smtpConfig.Host)

	subject := "Subject: New feedback from the website\n"
	from := fmt.Sprintf("From: %s\n", s.smtpConfig.From)
	to := fmt.Sprintf("To: %s\n", s.smtpConfig.To)
	message := []byte(from + to + subject + "\n" + messageBody)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.smtpConfig.Host, s.smtpConfig.Port),
		smtpAuth,
		s.smtpConfig.From,
		[]string{s.smtpConfig.To},
		message,
	)

	if err != nil {
		s.log.Error("Failed to send email", slog.String("error", err.Error()))
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.log.Info("Email successfully sent",
		slog.String("to", s.smtpConfig.To),
		slog.String("subject", "New feedback from the website"))

	return nil
}

func (s *FeedbackService) createEmailBody(feedback models.FeedbackDto) string {
	var sb strings.Builder

	sb.WriteString("New feedback from the website\n\n")
	sb.WriteString(fmt.Sprintf("Name: %s\n", feedback.Name))
	sb.WriteString(fmt.Sprintf("Email: %s\n", feedback.Email))

	if feedback.CityOrRegion != "" {
		sb.WriteString(fmt.Sprintf("City/Region: %s\n", feedback.CityOrRegion))
	}

	if feedback.Organization != "" {
		sb.WriteString(fmt.Sprintf("Organization: %s\n", feedback.Organization))
	}

	if feedback.Position != "" {
		sb.WriteString(fmt.Sprintf("Position: %s\n", feedback.Position))
	}

	if feedback.GratitudeExpression != "" {
		sb.WriteString(fmt.Sprintf("Expression of gratitude: %s\n", feedback.GratitudeExpression))
	}

	sb.WriteString(fmt.Sprintf("Message: %s\n", feedback.Message))

	return sb.String()
}

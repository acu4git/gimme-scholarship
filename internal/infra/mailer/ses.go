package mailer

import (
	"bytes"
	"context"
	"errors"
	"log"
	"text/template"

	"github.com/acu4git/gimme-scholarship/internal/service"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type sesMailer struct {
	client      *sesv2.Client
	fromAddress string
}

func NewSESMailer(ctx context.Context, fromAddress string) (service.Mailer, error) {
	if fromAddress == "" {
		return nil, errors.New("from address is required")
	}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &sesMailer{
		client:      sesv2.NewFromConfig(cfg),
		fromAddress: fromAddress,
	}, nil
}

func (m *sesMailer) SendEmail(ctx context.Context, to, subject, body string) error {
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(m.fromAddress),
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: aws.String(subject)},
				Body:    &types.Body{Text: &types.Content{Data: aws.String(body)}},
			},
		},
	}
	_, err := m.client.SendEmail(ctx, input)
	return err
}

func (m *sesMailer) SendBulkEmail(ctx context.Context, subject string, tmpl *template.Template, data []service.BulkEmailData) error {
	var sendErrors []error

	for _, item := range data {
		var body bytes.Buffer
		if err := tmpl.Execute(&body, item.TemplateData); err != nil {
			log.Printf("template execution failed for %s: %v", item.To, err)
			sendErrors = append(sendErrors, err)
			continue
		}

		if err := m.SendEmail(ctx, item.To, subject, body.String()); err != nil {
			log.Printf("failed to send email to %s: %v", item.To, err)
			sendErrors = append(sendErrors, err)
			continue
		}
	}

	if len(sendErrors) > 0 {
		return sendErrors[0]
	}

	return nil
}

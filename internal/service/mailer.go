package service

import (
	"context"

	"github.com/acu4git/gimme-scholarship/internal/domain/model"
)

type BulkEmailData struct {
	To           string
	TemplateData map[string]any
}

type Mailer interface {
	SendEmail(ctx context.Context, to, subject, body string) error
	SendBulkEmail(ctx context.Context, mailKey model.MailKey, data []BulkEmailData) error
}

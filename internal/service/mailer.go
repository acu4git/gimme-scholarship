package service

import (
	"context"
	"text/template"
)

type BulkEmailData struct {
	To           string
	TemplateData any
}

type Mailer interface {
	SendEmail(ctx context.Context, to, subject, body string) error
	SendBulkEmail(ctx context.Context, subject string, tmpl *template.Template, data []BulkEmailData) error
}

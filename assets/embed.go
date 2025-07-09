package assets

import (
	"embed"
	"text/template"
)

//go:embed all:mail-templates
var mailTemplateFiles embed.FS
var MailTemplatesTpl = template.Must(template.ParseFS(mailTemplateFiles, "mail-templates/*/*.txt")).Option("missingkey=error")

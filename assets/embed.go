package assets

import (
	"embed"
	"text/template"
)

//go:embed mail-templates/*/*
var mailTemplateFiles embed.FS
var MailTemplatesTpl = template.Must(template.ParseFS(mailTemplateFiles, "mail-templates/*/*")).Option("missingkey=error")

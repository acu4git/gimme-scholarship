package model

import (
	"bytes"
	"fmt"

	"github.com/acu4git/gimme-scholarship/assets"
)

type MailKey string

const (
	MailKeyNotifyScholarshipDeadline MailKey = "notify_scholarship_deadline"
	MailKeyNotifyScholarshipUpdate   MailKey = "notify_scholarship_update"
	MailKeyNotifyUserWelcome         MailKey = "notify_user_welcome"
)

func (m MailKey) getMailFooter() (string, error) {
	var disclaimerBuf bytes.Buffer
	if err := assets.MailTemplatesTpl.ExecuteTemplate(&disclaimerBuf, "footer_disclaimer.txt", nil); err != nil {
		return "", fmt.Errorf("failed to execute template 'footer_disclaimer.txt': %w", err)
	}

	var appInfoBuf bytes.Buffer
	if err := assets.MailTemplatesTpl.ExecuteTemplate(&appInfoBuf, "footer_app_info.txt", nil); err != nil {
		return "", fmt.Errorf("failed to execute template 'footer_app_info.txt': %w", err)
	}

	data := make(map[string]string)
	data["Disclaimer"] = disclaimerBuf.String()
	data["AppInfo"] = appInfoBuf.String()

	var footerBuf bytes.Buffer
	if err := assets.MailTemplatesTpl.ExecuteTemplate(&footerBuf, "footer.txt", data); err != nil {
		return "", fmt.Errorf("failed to execute template 'footer.txt': %w", err)
	}

	return footerBuf.String(), nil
}

func (mk MailKey) MailBodyWithFooter(data map[string]any) (string, string, error) {
	footer, err := mk.getMailFooter()
	if err != nil {
		return "", "", err
	}

	if data == nil {
		data = make(map[string]any)
	}
	data["Footer"] = footer

	var subjectBuf bytes.Buffer
	if err := assets.MailTemplatesTpl.ExecuteTemplate(&subjectBuf, fmt.Sprintf("%s_subject.txt", mk), data); err != nil {
		return "", "", err
	}

	var bodyBuf bytes.Buffer
	if err := assets.MailTemplatesTpl.ExecuteTemplate(&bodyBuf, fmt.Sprintf("%s_body.txt", mk), data); err != nil {
		return "", "", err
	}

	return subjectBuf.String(), bodyBuf.String(), nil
}

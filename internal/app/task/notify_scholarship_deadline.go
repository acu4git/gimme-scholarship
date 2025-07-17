package task

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/acu4git/gimme-scholarship/assets"
	"github.com/acu4git/gimme-scholarship/internal/domain/model"
	"github.com/acu4git/gimme-scholarship/internal/domain/repository"
	"github.com/acu4git/gimme-scholarship/internal/service"
)

type NotifyScholarshipDeadlineExecutor struct {
	repository repository.Repository
	mailer     service.Mailer
}

func NewNotifyScholarshipDeadlineExecutor(repository repository.Repository, mailer service.Mailer) *NotifyScholarshipDeadlineExecutor {
	return &NotifyScholarshipDeadlineExecutor{
		repository: repository,
		mailer:     mailer,
	}
}

func (e *NotifyScholarshipDeadlineExecutor) Execute() error {
	log.Println("Starting task: Notify Scholarship Deadline")

	userScholarships, err := e.repository.FindUsersToNotifyForUpcomingDeadlines()
	if err != nil {
		return fmt.Errorf("failed to find users to notify: %w", err)
	}

	if len(userScholarships) == 0 {
		log.Println("No users to notify today.")
		return nil
	}

	tmpl, err := assets.MailTemplatesTpl.Clone()
	if err != nil {
		return fmt.Errorf("failed to clone template: %w", err)
	}

	bodyTmpl := tmpl.Lookup("deadline_notification_body.txt")
	if bodyTmpl == nil {
		return fmt.Errorf("template 'deadline_notification_body.txt' not found")
	}

	subjectTmpl := tmpl.Lookup("deadline_notification_subject.txt")
	if subjectTmpl == nil {
		return fmt.Errorf("template 'deadline_notification_subject.txt' not found")
	}

	var subjectBuf bytes.Buffer
	if err := subjectTmpl.Execute(&subjectBuf, nil); err != nil {
		return fmt.Errorf("failed to execute subject template: %w", err)
	}

	var disclaimerBuf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&disclaimerBuf, "footer_disclaimer.txt", nil); err != nil {
		return fmt.Errorf("failed to execute disclaimer template: %w", err)
	}

	var appInfoBuf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&appInfoBuf, "footer_app_info", nil); err != nil {
		return fmt.Errorf("failed to execute app info template: %w", err)
	}

	bulkData := make([]service.BulkEmailData, 0, len(userScholarships))
	for email, scholarships := range userScholarships {
		templateData := struct {
			Scholarships []model.Scholarship
			Disclaimer   string
			AppInfo      string
		}{
			Scholarships: scholarships,
			Disclaimer:   disclaimerBuf.String(),
			AppInfo:      appInfoBuf.String(),
		}
		bulkData = append(bulkData, service.BulkEmailData{
			To:           email,
			TemplateData: templateData,
		})
	}

	ctx := context.Background()
	log.Printf("Sending deadline notifications to %d users...\n", len(bulkData))
	if err := e.mailer.SendBulkEmail(ctx, subjectBuf.String(), bodyTmpl, bulkData); err != nil {
		return fmt.Errorf("failed to send bulk email: %w", err)
	}

	log.Println("Successfully sent all deadline notifications.")
	return nil
}

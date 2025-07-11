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

	bulkData := make([]service.BulkEmailData, 0, len(userScholarships))
	for email, scholarships := range userScholarships {
		templateData := struct {
			Scholarships []model.Scholarship
		}{
			Scholarships: scholarships,
		}
		bulkData = append(bulkData, service.BulkEmailData{
			To:           email,
			TemplateData: templateData,
		})
	}

	tmpl, err := assets.MailTemplatesTpl.Clone()
	if err != nil {
		return fmt.Errorf("failed to clone template: %w", err)
	}

	bodyTmpl := tmpl.Lookup("mail-templates/deadline_notification/body.txt")
	if bodyTmpl == nil {
		return fmt.Errorf("template 'mail-templates/deadline_notification/body.txt' not found")
	}

	subjectTmpl := tmpl.Lookup("mail-templates/deadline_notification/subject.txt")
	if subjectTmpl == nil {
		return fmt.Errorf("template 'mail-templates/deadline_notification/subject.txt' not found")
	}

	var subjectBuf bytes.Buffer
	if err := subjectTmpl.Execute(&subjectBuf, nil); err != nil {
		return fmt.Errorf("failed to execute subject template: %w", err)
	}

	ctx := context.Background()
	log.Printf("Sending deadline notifications to %d users...", len(bulkData))
	if err := e.mailer.SendBulkEmail(ctx, subjectBuf.String(), bodyTmpl, bulkData); err != nil {
		return fmt.Errorf("failed to send bulk email: %w", err)
	}

	log.Println("Successfully sent all deadline notifications.")
	return nil
}

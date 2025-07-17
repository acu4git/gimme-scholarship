package task

import (
	"context"
	"fmt"
	"log"
	"strings"

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
		text := ""
		for i := range scholarships {
			text += fmt.Sprintf("=== %d件目\n", i+1)
			text += fmt.Sprintf("【奨学金名】\n%s\n", scholarships[i].Name)
			text += fmt.Sprintf("【対象】\n%s\n", strings.Join(scholarships[i].Targets, "・"))
			text += fmt.Sprintf("【形式】\n%s\n\n", scholarships[i].TypeDetail)
		}
		templateData := map[string]any{
			"Scholarships": text,
		}
		bulkData = append(bulkData, service.BulkEmailData{
			To:           email,
			TemplateData: templateData,
		})
	}

	ctx := context.Background()
	log.Printf("Sending deadline notifications to %d users...\n", len(bulkData))
	if err := e.mailer.SendBulkEmail(ctx, model.MailKeyNotifyScholarshipDeadline, bulkData); err != nil {
		return fmt.Errorf("failed to send bulk email: %w", err)
	}

	log.Println("Successfully sent all deadline notifications.")
	return nil
}

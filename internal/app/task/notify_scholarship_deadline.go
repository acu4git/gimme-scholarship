package task

import (
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

func (e *NotifyScholarshipDeadlineExecutor) Exectute() error {
	return nil
}

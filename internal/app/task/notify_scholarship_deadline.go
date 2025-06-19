package task

import "github.com/acu4git/gimme-scholarship/internal/domain/repository"

type NotifyScholarshipDeadlineExecutor struct {
	repository repository.Repository
}

func NewNotifyScholarshipDeadlineExecutor(repository repository.Repository) *NotifyScholarshipDeadlineExecutor {
	return &NotifyScholarshipDeadlineExecutor{repository: repository}
}

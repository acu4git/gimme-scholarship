package main

import (
	"context"
	"flag"
	"log"

	"github.com/acu4git/gimme-scholarship/internal/app/task"
	"github.com/acu4git/gimme-scholarship/internal/infra/mailer"
	"github.com/acu4git/gimme-scholarship/internal/service"
)

var (
	taskName = flag.String("task", "", "task name")
)

const (
	notifyScholarshipDeadline = "NOTIFY_SCHOLARSHIP_DEADLINE"
)

func main() {
	repo, err := service.CreateRepository()
	if err != nil {
		log.Fatal(err)
	}

	switch *taskName {
	case notifyScholarshipDeadline:
		ctx := context.Background()
		from := "no-reply@kit-gimme-scholarship.com"
		mailer, err := mailer.NewSESMailer(ctx, from)
		if err != nil {
			log.Fatal(err)
		}

		executor := task.NewNotifyScholarshipDeadlineExecutor(repo, mailer)
		if err := executor.Execute(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("invalid task name: %s", *taskName)
	}
}

package main

import (
	"flag"
	"log"

	"github.com/acu4git/gimme-scholarship/internal/app/task"
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
		executor := task.NewNotifyScholarshipDeadlineExecutor(repo)
	default:
		log.Fatalf("invalid task name: %s", *taskName)
	}
}

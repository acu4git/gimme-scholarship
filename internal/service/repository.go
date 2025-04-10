package service

import (
	"github.com/yummynaan/kit-enrollment-helper/internal/database"
	domain "github.com/yummynaan/kit-enrollment-helper/internal/domain"
)

func CreateRepository() (domain.Repository, error) {
	return database.NewDatabase()
}

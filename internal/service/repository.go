package service

import (
	"github.com/acu4git/gimme-scholarship/internal/database"
	"github.com/acu4git/gimme-scholarship/internal/domain/repository"
)

func CreateRepository() (repository.Repository, error) {
	return database.NewDatabase()
}

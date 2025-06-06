package repository

import "github.com/acu4git/gimme-scholarship/internal/domain/model"

type Repository interface {
	CreateUser(UserInput) error
	GetScholarships(FilterOption) ([]model.Scholarship, error)
	UserFavoriteAction(UserFavoriteInput) error
}

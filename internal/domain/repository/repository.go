package repository

import "github.com/acu4git/gimme-scholarship/internal/domain/model"

type Repository interface {
	CreateUser(UserInput) error
	FindScholarships(FilterOption) (scholarships []model.Scholarship, favoriteMap map[int64]bool, err error)
	UserFavoriteAction(UserFavoriteInput) error
}

package repository

import "github.com/acu4git/gimme-scholarship/internal/domain/model"

type Repository interface {
	CreateUser(model.User) error
}

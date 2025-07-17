package database

import (
	"github.com/acu4git/gimme-scholarship/internal/domain/model"
	"github.com/acu4git/gimme-scholarship/internal/domain/repository"
	"github.com/acu4git/gimme-scholarship/internal/testutil"
)

const idPrefix = "user_"

func (fdb *FakeDatabase) TestInitUsers() error {
	inputs := []repository.UserInput{
		{
			// Name:        "テスト太郎",
			ID:          idPrefix + testutil.RandLetters(32),
			Email:       "test-taro@example.com",
			Level:       model.EducationLevelMaster.Str(),
			Grade:       1,
			AcceptEmail: true,
		}, {
			// Name:        "テスト次郎",
			ID:          idPrefix + testutil.RandLetters(32),
			Email:       "test-jiro@example.com",
			Level:       model.EducationLevelBacholor.Str(),
			Grade:       3,
			AcceptEmail: true,
		}, {
			// Name:        "テスト・クルーズ",
			ID:          idPrefix + testutil.RandLetters(32),
			Email:       "test-cruise@example.com",
			Level:       model.EducationLevelOther.Str(),
			Grade:       4,
			AcceptEmail: true,
		}, {
			// Name:        "テスト花子",
			ID:          idPrefix + testutil.RandLetters(32),
			Email:       "test-hanako@example.com",
			Level:       model.EducationLevelBacholor.Str(),
			Grade:       3,
			AcceptEmail: true,
		},
	}

	for _, input := range inputs {
		if err := fdb.CreateUser(input); err != nil {
			return err
		}
	}

	return nil
}

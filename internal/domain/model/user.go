package model

import (
	"errors"
	"unicode/utf8"
)

type User struct {
	Name           UserName
	Email          string
	EducationLevel EducationLevel
}

type UserName string

func (un UserName) String() string {
	return string(un)
}

func (u User) Validate() error {
	n := utf8.RuneCountInString(u.Name.String())
	if n == 0 {
		return errors.New("username is required")
	}
	if n > 15 {
		return errors.New("username length must be 15 characters or fewer")
	}
	if !u.EducationLevel.Validate() {
		return errors.New("education level is invalid")
	}
	return nil
}

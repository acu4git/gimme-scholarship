package model

import (
	"errors"
	"unicode/utf8"
)

type User struct {
	Name           UserName
	Email          string
	Password       string
	EnrollmentYear int
	Department     Department
}

type UserName string

func (un UserName) String() string {
	return string(un)
}

func (un UserName) Validate() error {
	n := utf8.RuneCountInString(un.String())
	if n == 0 {
		return errors.New("username is required")
	}
	if n > 15 {
		return errors.New("username length must be 15 characters or fewer")
	}
	return nil
}

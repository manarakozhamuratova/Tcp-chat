package model

import (
	"errors"
	"net/http"
)

type ErrorWeb struct {
	StatusCode int
	StatusText string
}

var (
	ErrUserNotFound       = errors.New("there is no such user")
	ErrUserExists         = errors.New("a user with this username or email exists")
	ErrPostNotFound       = errors.New("post does not exist")
	ErrCommentNotFound    = errors.New("there are no comments with this ID")
	ErrNoSession          = errors.New("the user does not have a session")
	ErrValueNotSet        = errors.New("the required values are not assigned")
	ErrInsertFailed       = errors.New("create record was failed")
	ErrUpdateFailed       = errors.New("update table was failed")
	ErrDeleteFromDBFailed = errors.New("delete from db was failed")
	ErrMessageInvalid     = errors.New("enter data correctly")
)

func NewErrorWeb(statusCode int) *ErrorWeb {
	return &ErrorWeb{
		StatusCode: statusCode,
		StatusText: http.StatusText(statusCode),
	}
}

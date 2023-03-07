package services

import (
	"database/sql"
	"errors"
	"fmt"

	"forum/internal/model"
	"forum/internal/repository"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignIn(user *model.User, session *model.Session) error
	SignUp(user *model.User) error
	SignOut(session *model.Session) error
}

type authService struct {
	repository.UserQuery
	repository.SessionQuery
}

func NewAuthService(dao repository.DAO) AuthService {
	return &authService{
		UserQuery:    dao.NewUserQuery(),
		SessionQuery: dao.NewSessionQuery(),
	}
}

func (a *authService) SignIn(user *model.User, session *model.Session) error {
	err := a.UserQuery.UserVerification(user)
	if err != nil {
		return fmt.Errorf("signIn: %w", err)
	}

	tempSession := model.Session{User: *user}
	err = a.SessionQuery.GetSessionByUserID(&tempSession)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("signIn: %w", err)
	} else if err == nil {
		if err := a.SessionQuery.DeleteSession(&tempSession); err != nil {
			return fmt.Errorf("signIn: %w", err)
		}
	}

	token, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("signIn: %w", err)
	}

	session.User = *user
	session.Token = token.String()

	return a.SessionQuery.CreateSession(session)
}

func (a *authService) SignUp(user *model.User) error {
	exist, err := a.UserQuery.IsExistUser(user)
	if err != nil {
		return fmt.Errorf("signUp: %w", err)
	}

	if exist {
		return fmt.Errorf("signUp: %w", model.ErrUserExists)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return fmt.Errorf("signUp: %w", err)
	}

	user.Password = string(hash)
	return a.UserQuery.CreateUser(user)
}

func (a *authService) SignOut(session *model.Session) error {
	return a.SessionQuery.DeleteSession(session)
}

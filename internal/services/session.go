package services

import (
	"forum/internal/model"
	"forum/internal/repository"
)

type SessionService interface {
	CreateSession(session *model.Session) error
	DeleteSession(session *model.Session) error
	GetSession(session *model.Session) error
}

type sessionService struct {
	repository.SessionQuery
}

func NewSessionService(dao repository.DAO) SessionService {
	return &sessionService{
		SessionQuery: dao.NewSessionQuery(),
	}
}

func (s *sessionService) CreateSession(session *model.Session) error {
	return s.SessionQuery.CreateSession(session)
}

func (s *sessionService) DeleteSession(session *model.Session) error {
	return s.SessionQuery.DeleteSession(session)
}

func (s *sessionService) GetSession(session *model.Session) error {
	return s.SessionQuery.GetSession(session)
}

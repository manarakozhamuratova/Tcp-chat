package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"forum/internal/model"
)

type SessionQuery interface {
	CreateSession(session *model.Session) error
	DeleteSession(session *model.Session) error
	GetSession(session *model.Session) error
	GetSessionByUserID(session *model.Session) error
}

type sessionQuery struct {
	db *sql.DB
}

func (s *sessionQuery) CreateSession(session *model.Session) error {
	sqlStmt := `INSERT INTO sessions(user_id, token, expiry) VALUES(?,?,?)`
	query, err := s.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("createSession: %w", err)
	}

	defer query.Close()

	result, err := query.Exec(session.User.ID, session.Token, session.Expiry)
	if err != nil {
		return fmt.Errorf("createSession: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("createSession: %w", err)
	}

	session.ID = id
	return nil
}

func (s *sessionQuery) DeleteSession(session *model.Session) error {
	sqlStmt := `DELETE FROM sessions WHERE token=?`
	query, err := s.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("deleteSession: %w", err)
	}

	defer query.Close()

	result, err := query.Exec(session.Token)
	if err != nil {
		return fmt.Errorf("deleteSession: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("deleteSession: %w", err)
	}

	if rowsAffected == 0 {
		return model.ErrDeleteFromDBFailed
	}

	return nil
}

func (s *sessionQuery) GetSession(session *model.Session) error {
	sqlStmt := `SELECT session_id, user_id, expiry FROM sessions WHERE token=?`
	query, err := s.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("getSession: %w", err)
	}

	defer query.Close()

	err = query.QueryRow(session.Token).Scan(&session.ID, &session.User.ID, &session.Expiry)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = model.ErrNoSession
		}
		return fmt.Errorf("getSession: %w", err)
	}

	return nil
}

func (s *sessionQuery) GetSessionByUserID(session *model.Session) error {
	sqlStmt := `SELECT session_id, token, expiry FROM sessions WHERE user_id=?`
	query, err := s.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("getSessionByUserID: %w", err)
	}

	defer query.Close()

	err = query.QueryRow(session.User.ID).Scan(&session.ID, &session.Token, &session.Expiry)
	if err != nil {
		return fmt.Errorf("getSessionByUserID: %w", err)
	}

	return nil
}

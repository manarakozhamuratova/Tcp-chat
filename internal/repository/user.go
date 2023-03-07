package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"forum/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type UserQuery interface {
	CreateUser(user *model.User) error
	DeleteUser(userID int64) error
	UserVerification(user *model.User) error
	IsExistUser(user *model.User) (bool, error)
	GetUserInfo(user *model.User) error
}

type userQuery struct {
	db *sql.DB
}

func (u *userQuery) CreateUser(user *model.User) error {
	sqlStmt := `INSERT INTO users(username, email, password) 
	VALUES(?, ?, ?)`

	query, err := u.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("createUser: %w", err)
	}

	defer query.Close()

	result, err := query.Exec(user.Username, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("createUser: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("createUser: %w", err)
	}

	user.ID = id
	return nil
}

func (u *userQuery) UserVerification(user *model.User) error {
	sqlStmt := `SELECT user_id, username, password FROM users WHERE email=?`
	query, err := u.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("userVerification: %w", err)
	}

	defer query.Close()

	tempPasswd := user.Password
	err = query.QueryRow(user.Email).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = model.ErrUserNotFound
		}
		return fmt.Errorf("userVerification: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tempPasswd))
	if err != nil {
		return fmt.Errorf("userVerification: %w", err)
	}

	return nil
}

func (u *userQuery) DeleteUser(userID int64) error {
	sqlStmt := `DELETE FROM users WHERE user_id=?`
	result, err := u.db.Exec(sqlStmt, userID)
	if err != nil {
		return fmt.Errorf("deleteUser: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("deleteUser: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("deleteUser: %w", model.ErrDeleteFromDBFailed)
	}

	return nil
}

func (u *userQuery) IsExistUser(user *model.User) (bool, error) {
	sqlStmt := `SELECT EXISTS(SELECT 1 FROM users WHERE username=? OR email=? LIMIT 1)`
	var exist bool

	err := u.db.QueryRow(sqlStmt, user.Username, user.Email).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("isExistUser: %w", err)
	}

	return exist, nil
}

func (u *userQuery) GetUserInfo(user *model.User) error {
	sqlStmt := `SELECT username, email, password FROM users WHERE user_id=?`
	query, err := u.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("getUserInfo: %w", err)
	}

	defer query.Close()

	err = query.QueryRow(user.ID).Scan(&user.Username, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = model.ErrUserNotFound
		}
		return fmt.Errorf("getUserInfo: %w", err)
	}

	return nil
}

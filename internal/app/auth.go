package app

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"forum/internal/app/validation"
	"forum/internal/model"

	"golang.org/x/crypto/bcrypt"
)

var SignUpMessage string

func (s *ServiceServer) SignIn(w http.ResponseWriter, r *http.Request) {
	_, err := s.getSession(r)
	if err != nil && !errors.Is(err, model.ErrNoSession) && !errors.Is(err, model.ErrUserNotFound) {
		log.Println("ERROR:\nSignIn:", err)
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	} else if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if r.Method == http.MethodGet {
		s.render(w, "signin", nil)
		return
	}
	s.PostSignIn(w, r)
	return
}

func (s *ServiceServer) SignUp(w http.ResponseWriter, r *http.Request) {
	_, err := s.getSession(r)
	if err != nil && !errors.Is(err, model.ErrNoSession) && !errors.Is(err, model.ErrUserNotFound) {
		log.Println("ERROR:\nSignUp:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	} else if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if r.Method == http.MethodGet {
		s.render(w, "signup", nil)
		return
	}
	s.PostSignUp(w, r)
	return
}

func (s *ServiceServer) PostSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("ERROR:\nSignIn With POST Method: ", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadGateway))
		return
	}

	user := model.User{Email: r.PostFormValue("email"), Password: r.PostFormValue("password")}
	if err := validation.ValidationFormSignIn(user.Email, user.Password); err != nil {
		if errors.Is(err, model.ErrMessageInvalid) {
			s.render(w, "signin", model.ErrMessageInvalid.Error())
			return
		}
		log.Println("ERROR:\nSignIn With POST Method: ", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	session := model.Session{Expiry: time.Now().Add(time.Minute * 15)}

	err := s.authService.SignIn(&user, &session)
	if err != nil {
		log.Println("ERROR:\nSignIn With POST Method: ", err)

		if errors.Is(err, model.ErrUserNotFound) || errors.Is(err, sql.ErrNoRows) || errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			message := "user email or password incorrect"
			s.render(w, "signin", message)
			return
		}
		log.Println("ERROR:\nSignIn with POST Method:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	cookie := http.Cookie{
		Name:     "authToken",
		Value:    session.Token,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   900,
		Expires:  session.Expiry,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *ServiceServer) PostSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("ERROR:\nSignUp With POST Method: ", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadGateway))
		return
	}

	user := model.User{Username: r.PostFormValue("username"), Email: r.PostFormValue("email"), Password: r.PostFormValue("password"), Password2: r.PostFormValue("password2")}
	if err := validation.ValidationFormSignUp(user.Username, user.Email, user.Password, user.Password2); err != nil {

		if errors.Is(err, model.ErrMessageInvalid) {
			s.render(w, "signup", model.ErrMessageInvalid.Error())
			return
		}
		log.Println("ERROR:\nSignUp With POST Method: ", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))

		return
	}

	err := s.authService.SignUp(&user)
	if err != nil {
		if errors.Is(err, model.ErrUserExists) {
			s.render(w, "signup", model.ErrUserExists.Error())
			return
		}
		log.Println("ERROR:\nSignUp With POST Method: ", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	SignUpMessage = ""
	http.Redirect(w, r, "/signIn", http.StatusFound)
}

func (s *ServiceServer) SignOut(w http.ResponseWriter, r *http.Request, session model.Session) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	err := s.authService.SignOut(&session)
	if err != nil {
		log.Println("ERROR:\nSignOut: ", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

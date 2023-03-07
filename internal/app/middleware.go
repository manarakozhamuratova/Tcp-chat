package app

import (
	"errors"
	"log"
	"net/http"

	"forum/internal/model"
)

func (s *ServiceServer) authMiddleware(next func(http.ResponseWriter, *http.Request, model.Session)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.getSession(r)
		if err != nil {
			if errors.Is(err, model.ErrNoSession) || errors.Is(err, model.ErrUserNotFound) {
				s.ErrorHandler(w, model.NewErrorWeb(http.StatusUnauthorized))
				return
			}
			log.Println("authMiddleware:", err)
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
			return
		}

		next(w, r, session)
	}
}

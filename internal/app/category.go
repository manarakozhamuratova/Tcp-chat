package app

import (
	"errors"
	"log"
	"net/http"

	"forum/internal/model"
)

func (s *ServiceServer) Category(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	data, err := s.getPostsOfCategory(r)
	if err != nil {
		log.Println("ERROR:\nCategory:", err)
		if errors.Is(err, model.ErrValueNotSet) {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	data.Status = true

	if _, err = s.getSession(r); err != nil {
		if !errors.Is(err, model.ErrNoSession) && !errors.Is(err, model.ErrUserNotFound) {
			log.Println("ERROR:\nCategory:", err)

			s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
			return
		}
		data.Status = false
	}

	s.render(w, "index", data)
}

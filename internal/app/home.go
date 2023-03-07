package app

import (
	"errors"
	"log"
	"net/http"

	"forum/internal/model"
)

func (s *ServiceServer) IndexWithSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	if r.URL.Path != "/" {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusNotFound))
		return
	}

	data, err := s.getAllPosts(r)
	if err != nil {
		log.Println("ERROR:\nIndexWithSession:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	data.Status = true

	if _, err = s.getSession(r); err != nil {
		if !errors.Is(err, model.ErrUserNotFound) && !errors.Is(err, model.ErrNoSession) {
			log.Println("ERROR:\nIndexWithSession:", err)

			s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
			return
		}

		data.Status = false
	}

	s.render(w, "index", data)
}

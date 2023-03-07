package app

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"forum/internal/model"
)

func (s *ServiceServer) PostLike(w http.ResponseWriter, r *http.Request, session model.Session) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	postID, err := s.setPostReaction(session.User, s.postService.PostLike, r)
	if err != nil {
		log.Println("ERROR:\npostLike:", err)
		if errors.Is(err, model.ErrPostNotFound) || errors.Is(err, model.ErrValueNotSet) {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, "/post?ID="+strconv.Itoa(postID), http.StatusFound)
}

func (s *ServiceServer) PostDislike(w http.ResponseWriter, r *http.Request, session model.Session) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	postID, err := s.setPostReaction(session.User, s.postService.PostDislike, r)
	if err != nil {
		log.Println("ERROR:\npostDislike:", err)
		if errors.Is(err, model.ErrPostNotFound) || errors.Is(err, model.ErrValueNotSet) {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, "/post?ID="+strconv.Itoa(postID), http.StatusFound)
}

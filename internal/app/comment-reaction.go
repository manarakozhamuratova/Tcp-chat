package app

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"forum/internal/model"
)

func (s *ServiceServer) CommentLike(w http.ResponseWriter, r *http.Request, session model.Session) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	postID, err := s.setCommentReaction(session.User, s.postService.CommentSetLike, r)
	if err != nil {
		log.Println("ERROR:\nCommentLike:", err)

		if errors.Is(err, model.ErrCommentNotFound) || errors.Is(err, model.ErrValueNotSet) || errors.Is(err, model.ErrPostNotFound) {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, "/post?ID="+strconv.Itoa(int(postID)), http.StatusFound)
}

func (s *ServiceServer) CommentDislike(w http.ResponseWriter, r *http.Request, session model.Session) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	postID, err := s.setCommentReaction(session.User, s.postService.CommentSetDislike, r)
	if err != nil {
		log.Println("ERROR:\nCommentDislike:", err)

		if errors.Is(err, model.ErrCommentNotFound) || errors.Is(err, model.ErrValueNotSet) || errors.Is(err, model.ErrPostNotFound) {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, "/post?ID="+strconv.Itoa(int(postID)), http.StatusFound)
}

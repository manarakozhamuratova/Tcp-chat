package app

import (
	"errors"
	"log"
	"net/http"

	"forum/internal/app/validation"
	"forum/internal/model"
)

func (s *ServiceServer) NewPost(w http.ResponseWriter, r *http.Request, session model.Session) {
	if r.Method == http.MethodGet {
		s.GetNewPost(w, r)
		return
	}
	s.PostNewPost(w, r, &session.User)
}

func (s *ServiceServer) GetNewPost(w http.ResponseWriter, r *http.Request) {
	allCategories, err := s.postService.GetAllCategories()
	if err != nil {
		log.Println("ERROR:\ngetNewPost:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	s.render(w, "create-post", model.Data{Status: true, Categories: allCategories})
}

func (s *ServiceServer) PostNewPost(w http.ResponseWriter, r *http.Request, user *model.User) {
	if r.Method != http.MethodPost {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("ERROR:\npostNewPost:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadGateway))
		return
	}

	post := model.Post{User: *user}

	allCategories, err := s.postService.GetAllCategories()
	if err != nil {
		log.Println("ERROR:\npostNewPost:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	err = validation.CheckInput(r, &post, allCategories)
	if err != nil {
		log.Println("ERROR:\npostNewPost:", err)

		if errors.Is(err, model.ErrMessageInvalid) {
			s.render(w, "create-post", model.Data{Status: false, Categories: allCategories})
		} else {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		}
		return
	}

	err = s.postService.CreatePost(&post)
	if err != nil {
		log.Println("ERROR:\npostNewPost:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *ServiceServer) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	session := true

	if _, err := s.getSession(r); err != nil {
		if !errors.Is(err, model.ErrNoSession) && !errors.Is(err, model.ErrUserNotFound) {
			log.Println("ERROR:\npost:", err)

			s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
			return
		}
		session = false
	}

	post, err := s.getPost(r)
	if err != nil {
		log.Println("ERROR:\npost:", err)

		if errors.Is(err, model.ErrPostNotFound) || errors.Is(err, model.ErrValueNotSet) {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	data := model.DataPost{Session: session, Post: post}

	s.render(w, "post", data)
}

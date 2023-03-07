package app

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"forum/internal/model"
)

func (s *ServiceServer) getSession(r *http.Request) (model.Session, error) {
	cookie, err := r.Cookie("authToken")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return model.Session{}, fmt.Errorf("getSession: %w", model.ErrNoSession)
		}
		return model.Session{}, fmt.Errorf("getSession: %w", err)
	}

	session := model.Session{Token: cookie.Value}

	err = s.sessionService.GetSession(&session)
	if err != nil {
		return model.Session{}, fmt.Errorf("getSession: %w", err)
	}

	if session.Expiry.Before(time.Now()) {
		err = s.sessionService.DeleteSession(&session)
		if err != nil {
			return model.Session{}, fmt.Errorf("getSession: %w", err)
		}
		return model.Session{}, fmt.Errorf("getSession: %w", model.ErrNoSession)
	}

	err = s.userService.GetUserInfo(&session.User)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			err = s.sessionService.DeleteSession(&session)
			if err != nil {
				return model.Session{}, fmt.Errorf("getSession: %w", err)
			}
		}

		return model.Session{}, fmt.Errorf("getSession: %w", err)
	}

	return session, nil
}

func (s *ServiceServer) getID(r *http.Request) (int, error) {
	if r.URL.Query().Get("ID") == "" {
		return 0, fmt.Errorf("getID: %w", model.ErrValueNotSet)
	}

	id, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		return 0, fmt.Errorf("getID: %w", model.ErrValueNotSet)
	}
	return id, nil
}

func (s *ServiceServer) getAllPosts(r *http.Request) (model.Data, error) {
	posts, err := s.postService.GetAllPosts()
	if err != nil {
		return model.Data{}, fmt.Errorf("getAllPosts: %w", err)
	}

	categories, err := s.postService.GetAllCategories()
	if err != nil {
		return model.Data{}, fmt.Errorf("getAllPosts: %w", err)
	}

	return model.Data{Categories: categories, Posts: posts}, nil
}

func (s *ServiceServer) getPostsOfCategory(r *http.Request) (model.Data, error) {
	id, err := s.getID(r)
	if err != nil {
		return model.Data{}, fmt.Errorf("getPostsOfCategory: %w", err)
	}

	category := model.Category{ID: int64(id)}
	posts, err := s.postService.GetPostsOfCategory(category)
	if err != nil {
		return model.Data{}, fmt.Errorf("getPostsOfCategory: %w", err)
	}

	categories, err := s.postService.GetAllCategories()
	if err != nil {
		return model.Data{}, fmt.Errorf("getPostsOfCategory: %w", err)
	}

	return model.Data{Categories: categories, Posts: posts}, nil
}

func (s *ServiceServer) getPost(r *http.Request) (model.Post, error) {
	postID, err := s.getID(r)
	if err != nil {
		return model.Post{}, fmt.Errorf("getPost: %w", err)
	}

	post := model.Post{ID: int64(postID)}
	err = s.postService.GetPost(&post)
	if err != nil {
		return model.Post{}, fmt.Errorf("getPost: %w", err)
	}

	return post, nil
}

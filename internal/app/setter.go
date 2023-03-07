package app

import (
	"fmt"
	"net/http"

	"forum/internal/model"
)

func (s *ServiceServer) setPostReaction(user model.User, f func(*model.PostReaction) error, r *http.Request) (int, error) {
	postID, err := s.getID(r)
	if err != nil {
		return 0, fmt.Errorf("setPostReaction: %w", err)
	}

	post := model.Post{ID: int64(postID)}
	err = f(&model.PostReaction{Post: post, User: user})
	if err != nil {
		return 0, fmt.Errorf("setPostReaction: %w", err)
	}

	return postID, nil
}

func (s *ServiceServer) setCommentReaction(user model.User, f func(*model.CommentReaction) error, r *http.Request) (int, error) {
	commentID, err := s.getID(r)
	if err != nil {
		return 0, fmt.Errorf("setCommentReaction: %w", err)
	}

	reaction := model.CommentReaction{Comment: model.Comment{ID: int64(commentID)}, User: user}
	err = f(&reaction)
	if err != nil {
		return 0, fmt.Errorf("setCommentReaction: %w", err)
	}

	err = s.postService.GetCommentInfo(&reaction.Comment)
	if err != nil {
		return 0, fmt.Errorf("setCommentReaction: %w", err)
	}

	return int(reaction.Comment.PostID), nil
}

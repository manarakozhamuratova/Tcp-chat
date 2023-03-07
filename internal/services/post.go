package services

import (
	"forum/internal/model"
	"forum/internal/repository"
)

type PostService interface {
	CreatePost(post *model.Post) error
	GetPost(post *model.Post) error
	CreateComment(comment *model.Comment) error
	GetPostComments(post *model.Post) error
	GetAllPosts() ([]model.Post, error)
	PostLike(reaction *model.PostReaction) error
	PostDislike(reaction *model.PostReaction) error
	CommentSetLike(reaction *model.CommentReaction) error
	CommentSetDislike(reaction *model.CommentReaction) error
	GetAllCategories() ([]model.Category, error)
	GetPostsOfCategory(category model.Category) ([]model.Post, error)
	GetCommentInfo(comment *model.Comment) error
	GetUserPosts(user model.User) ([]model.Post, error)
	GetRatedPosts(user model.User) ([]model.Post, error)
}

type postService struct {
	repository.PostQuery
	repository.CommentQuery
}

func NewPostService(dao repository.DAO) PostService {
	return &postService{
		PostQuery:    dao.NewPostQuery(),
		CommentQuery: dao.NewCommentQuery(),
	}
}

func (p *postService) CreatePost(post *model.Post) error {
	return p.PostQuery.CreatePost(post)
}

func (p *postService) GetPost(post *model.Post) error {
	err := p.PostQuery.GetPost(post)
	if err != nil {
		return err
	}

	return p.CommentQuery.GetPostComments(post)
}

func (p *postService) CreateComment(comment *model.Comment) error {
	return p.CommentQuery.CreateComment(comment)
}

func (p *postService) GetPostComments(post *model.Post) error {
	return p.CommentQuery.GetPostComments(post)
}

func (p *postService) GetAllPosts() ([]model.Post, error) {
	return p.PostQuery.GetAllPosts()
}

func (p *postService) PostLike(reaction *model.PostReaction) error {
	return p.PostQuery.PostSetLike(reaction)
}

func (p *postService) PostDislike(reaction *model.PostReaction) error {
	return p.PostQuery.PostSetDislike(reaction)
}

func (p *postService) CommentSetLike(reaction *model.CommentReaction) error {
	return p.CommentQuery.CommentSetLike(reaction)
}

func (p *postService) CommentSetDislike(reaction *model.CommentReaction) error {
	return p.CommentQuery.CommentSetDislike(reaction)
}

func (p *postService) GetAllCategories() ([]model.Category, error) {
	return p.PostQuery.GetAllCategories()
}

func (p *postService) GetPostsOfCategory(category model.Category) ([]model.Post, error) {
	return p.PostQuery.GetPostsOfCategory(category)
}

func (p *postService) GetCommentInfo(comment *model.Comment) error {
	return p.CommentQuery.GetCommentInfo(comment)
}

func (p *postService) GetUserPosts(user model.User) ([]model.Post, error) {
	return p.PostQuery.GetUserPosts(user)
}

func (p *postService) GetRatedPosts(user model.User) ([]model.Post, error) {
	return p.PostQuery.GetRatedPosts(user)
}

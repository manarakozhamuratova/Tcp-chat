package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"forum/internal/model"
)

type PostQuery interface {
	CreatePost(post *model.Post) error
	GetPost(post *model.Post) error
	GetAllPosts() ([]model.Post, error)
	PostSetLike(reaction *model.PostReaction) error
	PostSetDislike(reaction *model.PostReaction) error
	SetPostCategory(post *model.Post) error
	CreateCategory(category *model.Category) error
	GetAllCategories() ([]model.Category, error)
	GetPostsOfCategory(category model.Category) ([]model.Post, error)
	GetUserPosts(user model.User) ([]model.Post, error)
	GetRatedPosts(user model.User) ([]model.Post, error)
}

type postQuery struct {
	db *sql.DB
}

func (p *postQuery) CreatePost(post *model.Post) error {
	sqlStmt := `INSERT INTO posts (title, content, user_id, username)VALUES(?,?,?,?)`
	query, err := p.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("createPost: %w", err)
	}

	defer query.Close()

	result, err := query.Exec(post.Title, post.Content, post.User.ID, post.User.Username)
	if err != nil {
		return fmt.Errorf("createPost: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("createPost: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("createPost: %w", model.ErrInsertFailed)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("createPost: %w", err)
	}

	post.ID = id
	return p.SetPostCategory(post)
}

func (p *postQuery) GetPost(post *model.Post) error {
	sqlStmt := `SELECT title,content, user_id, username FROM posts WHERE post_id=?`
	query, err := p.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("getPost: %w", err)
	}

	defer query.Close()

	err = query.QueryRow(post.ID).Scan(&post.Title, &post.Content, &post.User.ID, &post.User.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = model.ErrPostNotFound
		}
		return fmt.Errorf("getPost: %w", err)
	}

	err = p.GetPostCategories(post)
	if err != nil {
		return fmt.Errorf("getPost: %w", err)
	}

	err = p.GetPostLikesDislikes(post)
	if err != nil {
		return fmt.Errorf("getPost: %w", err)
	}

	return nil
}

func (p *postQuery) GetAllPosts() ([]model.Post, error) {
	sqlStmt := `SELECT * FROM posts`
	rows, err := p.db.Query(sqlStmt)
	if err != nil {
		return []model.Post{}, fmt.Errorf("getAllPosts: %w", err)
	}

	defer rows.Close()

	posts := []model.Post{}
	for rows.Next() {
		post := model.Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.User.ID, &post.User.Username)
		if err != nil {
			return []model.Post{}, fmt.Errorf("getAllPosts: %w", err)
		}

		err = p.GetPostCategories(&post)
		if err != nil {
			return []model.Post{}, fmt.Errorf("getAllPosts: %w", err)
		}

		err = p.GetPostLikesDislikes(&post)
		if err != nil {
			return []model.Post{}, fmt.Errorf("getAllPosts: %w", err)
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func (p *postQuery) GetPostsOfCategory(category model.Category) ([]model.Post, error) {
	sqlStmt := `SELECT posts.post_id, posts.title, posts.content, posts.user_id, posts.username FROM posts
	INNER JOIN post_categories ON posts.post_id = post_categories.post_id
	INNER JOIN categories ON post_categories.category_id = categories.category_id
	WHERE categories.category_id = ?`

	query, err := p.db.Prepare(sqlStmt)
	if err != nil {
		return []model.Post{}, fmt.Errorf("getPostsOfCategory: %w", err)
	}

	defer query.Close()

	posts := []model.Post{}
	rows, err := query.Query(category.ID)
	if err != nil {
		return []model.Post{}, fmt.Errorf("getPostsOfCategory: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		post := model.Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.User.ID, &post.User.Username)
		if err != nil {
			return []model.Post{}, fmt.Errorf("getPostsOfCategory: %w", err)
		}

		err = p.GetPostCategories(&post)
		if err != nil {
			return []model.Post{}, fmt.Errorf("getPostsOfCategory: %w", err)
		}

		err = p.GetPostLikesDislikes(&post)
		if err != nil {
			return []model.Post{}, fmt.Errorf("getPostsOfCategory: %w", err)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (p *postQuery) GetUserPosts(user model.User) ([]model.Post, error) {
	sqlStmt := `SELECT * FROM posts WHERE user_id=?`
	rows, err := p.db.Query(sqlStmt, user.ID)
	if err != nil {
		return []model.Post{}, fmt.Errorf("getUserPosts: %w", err)
	}

	defer rows.Close()

	posts := []model.Post{}
	for rows.Next() {
		post := model.Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.User.ID, &post.User.Username)
		if err != nil {
			return []model.Post{}, fmt.Errorf("getUserPosts: %w", err)
		}

		err = p.GetPostCategories(&post)
		if err != nil {
			return []model.Post{}, fmt.Errorf("getUserPosts: %w", err)
		}

		err = p.GetPostLikesDislikes(&post)
		if err != nil {
			return []model.Post{}, fmt.Errorf("getUserPosts: %w", err)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (p *postQuery) GetRatedPosts(user model.User) ([]model.Post, error) {
	sqlStmt := `SELECT posts.post_id, posts.title, posts.content, posts.user_id, posts.username FROM posts
	INNER JOIN posts_likes_dislikes ON posts.post_id = posts_likes_dislikes.post_id
	WHERE posts_likes_dislikes.user_id=? AND (posts_likes_dislikes.like=1 OR posts_likes_dislikes.dislike=1)`
	query, err := p.db.Prepare(sqlStmt)
	if err != nil {
		return []model.Post{}, fmt.Errorf("getRatedPosts: %w", err)
	}

	defer query.Close()

	posts := []model.Post{}
	rows, err := query.Query(user.ID)
	if err != nil {
		return []model.Post{}, fmt.Errorf("getRatedPosts: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		post := model.Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.User.ID, &post.User.Username)
		if err != nil {
			return []model.Post{}, fmt.Errorf("getRatedPosts: %w", err)
		}

		err = p.GetPostCategories(&post)
		if err != nil {
			return []model.Post{}, fmt.Errorf("getRatedPosts: %w", err)
		}

		err = p.GetPostLikesDislikes(&post)
		if err != nil {
			return []model.Post{}, fmt.Errorf("getRatedPosts: %w", err)
		}

		posts = append(posts, post)
	}

	return posts, nil
}

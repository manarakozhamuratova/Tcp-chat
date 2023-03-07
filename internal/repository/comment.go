package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"forum/internal/model"
)

type CommentQuery interface {
	CreateComment(comment *model.Comment) error
	GetPostComments(post *model.Post) error
	CommentSetLike(reaction *model.CommentReaction) error
	CommentSetDislike(reaction *model.CommentReaction) error
	GetCommentInfo(comment *model.Comment) error
}

type commentQuery struct {
	db *sql.DB
}

func (c *commentQuery) CreateComment(comment *model.Comment) error {
	sqlStmt := `INSERT INTO comments (post_id, user_id, username, message) 
	SELECT post_id, ?, ?, ?
	FROM posts
	WHERE EXISTS (SELECT * FROM posts WHERE post_id=?) AND post_id=?`
	query, err := c.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("createComment: %w", err)
	}

	defer query.Close()

	result, err := query.Exec(comment.UserID, comment.Username, comment.Message, comment.PostID, comment.PostID)
	if err != nil {
		return fmt.Errorf("createComment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("createComment: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("createComment: %w", model.ErrPostNotFound)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("createComment: %w", err)
	}

	comment.ID = id
	return nil
}

func (c *commentQuery) GetPostComments(post *model.Post) error {
	sqlStmt := `SELECT comment_id, user_id, username, message FROM comments WHERE post_id=?`
	query, err := c.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("getPostComments: %w", err)
	}

	defer query.Close()

	rows, err := query.Query(post.ID)
	if err != nil {
		return fmt.Errorf("getPostComments: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		comment := model.Comment{}
		err = rows.Scan(&comment.ID, &comment.UserID, &comment.Username, &comment.Message)
		if err != nil {
			return fmt.Errorf("getPostComments: %w", err)
		}

		err = c.getCommentLikesDislikes(&comment)
		if err != nil {
			return fmt.Errorf("getPostComments: %w", err)
		}

		post.Comments = append(post.Comments, comment)
	}

	return nil
}

func (c *commentQuery) GetCommentInfo(comment *model.Comment) error {
	sqlStmt := `SELECT * FROM comments WHERE comment_id=?`
	err := c.db.QueryRow(sqlStmt, comment.ID).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Username, &comment.Message)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = model.ErrCommentNotFound
		}
		return fmt.Errorf("getCommentInfo: %w", err)
	}

	return nil
}

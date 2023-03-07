package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"forum/internal/model"
)

func (c *commentQuery) createReactionToComment(reaction *model.CommentReaction) error {
	sqlStmt := `INSERT INTO comments_likes_dislikes (comment_id, user_id, like, dislike) 
	SELECT comment_id, ?, 0, 0
	FROM comments
	WHERE EXISTS (SELECT * FROM posts WHERE comment_id=?)`
	query, err := c.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("createReactionToComment: %w", err)
	}

	defer query.Close()

	res, err := query.Exec(reaction.User.ID, reaction.Comment.ID)
	if err != nil {
		return fmt.Errorf("createReactionToComment: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("createReactionToComment: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("createReactionToComment: %w", model.ErrCommentNotFound)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("createReactionToComment: %w", err)
	}

	reaction.ID = id
	return nil
}

func (c *commentQuery) CommentSetLike(reaction *model.CommentReaction) error {
	var sqlStmt string
	err := c.getUserReactionToComment(reaction)
	if err != nil {
		return fmt.Errorf("commentSetLike: %w", err)
	}

	if reaction.Like == reaction.Dislike {
		sqlStmt = `UPDATE comments_likes_dislikes SET like=1 WHERE Id=?`
		err = c.updateCommentReaction(sqlStmt, c.db, reaction)
	} else if reaction.Like == 0 {
		sqlStmt = `UPDATE comments_likes_dislikes SET like=1, dislike=0 WHERE Id=?`
		err = c.updateCommentReaction(sqlStmt, c.db, reaction)
	} else {
		sqlStmt = `UPDATE comments_likes_dislikes SET like=0 WHERE Id=?`
		err = c.updateCommentReaction(sqlStmt, c.db, reaction)
	}
	if err != nil {
		return fmt.Errorf("commentSetLike: %w", err)
	}

	sqlStmt = `SELECT like, dislike FROM comments_likes_dislikes WHERE Id=?`
	err = c.db.QueryRow(sqlStmt, reaction.ID).Scan(&reaction.Like, &reaction.Dislike)
	if err != nil {
		return fmt.Errorf("commentSetLike: %w", err)
	}
	return nil
}

func (c *commentQuery) CommentSetDislike(reaction *model.CommentReaction) error {
	var sqlStmt string
	err := c.getUserReactionToComment(reaction)
	if err != nil {
		return fmt.Errorf("commentSetDislike: %w", err)
	}

	if reaction.Like == reaction.Dislike {
		sqlStmt = `UPDATE comments_likes_dislikes SET dislike=1 WHERE Id=?`
		err = c.updateCommentReaction(sqlStmt, c.db, reaction)
	} else if reaction.Dislike == 0 {
		sqlStmt = `UPDATE comments_likes_dislikes SET like=0, dislike=1 WHERE Id=?`
		err = c.updateCommentReaction(sqlStmt, c.db, reaction)
	} else {
		sqlStmt = `UPDATE comments_likes_dislikes SET dislike=0 WHERE Id=?`
		err = c.updateCommentReaction(sqlStmt, c.db, reaction)
	}

	if err != nil {
		return fmt.Errorf("commentSetDislike: %w", err)
	}

	sqlStmt = `SELECT like, dislike FROM comments_likes_dislikes WHERE Id=?`
	err = c.db.QueryRow(sqlStmt, reaction.ID).Scan(&reaction.Like, &reaction.Dislike)
	if err != nil {
		return fmt.Errorf("commentSetDislike: %w", err)
	}

	return nil
}

func (c *commentQuery) updateCommentReaction(sqlStmt string, db *sql.DB, reaction *model.CommentReaction) error {
	result, err := db.Exec(sqlStmt, reaction.ID)
	if err != nil {
		return fmt.Errorf("updateCommentReaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("updateCommentReaction: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("updateCommentReaction: %w", model.ErrUpdateFailed)
	}
	return nil
}

func (c *commentQuery) getUserReactionToComment(reaction *model.CommentReaction) error {
	sqlStmt := `SELECT id, like, dislike FROM comments_likes_dislikes WHERE comment_id=? and user_id=?`
	query, err := c.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("getUserReactionToComment: %w", err)
	}

	defer query.Close()

	err = query.QueryRow(reaction.Comment.ID, reaction.User.ID).Scan(&reaction.ID, &reaction.Like, &reaction.Dislike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.createReactionToComment(reaction)
		}

		return fmt.Errorf("getUserReactionToComment: %w", err)
	}
	return nil
}

func (c *commentQuery) getCommentLikesDislikes(comment *model.Comment) error {
	sqlStmt := `SELECT COALESCE(SUM(like), 0), COALESCE(SUM(dislike), 0) FROM comments_likes_dislikes WHERE comment_id=?`
	query, err := c.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("getCommentLikesDislikes: %w", err)
	}

	defer query.Close()

	err = query.QueryRow(comment.ID).Scan(&comment.Like, &comment.Dislike)
	if err != nil {
		return fmt.Errorf("getCommentLikesDislikes: %w", err)
	}

	return nil
}

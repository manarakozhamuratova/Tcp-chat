package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"forum/internal/model"
)

func (pr *postQuery) PostSetLike(reaction *model.PostReaction) error {
	var sqlStmt string

	err := pr.getUserReactionToPost(reaction)
	if err != nil {
		return fmt.Errorf("postSetLike: %w", err)
	}

	if reaction.Like == reaction.Dislike {
		sqlStmt = `UPDATE posts_likes_dislikes SET like=1 WHERE Id=?`
		err = pr.updatePostReaction(sqlStmt, pr.db, reaction)
	} else if reaction.Like == 0 {
		sqlStmt = `UPDATE posts_likes_dislikes SET like=1, dislike=0 WHERE Id=?`
		err = pr.updatePostReaction(sqlStmt, pr.db, reaction)
	} else {
		sqlStmt = `UPDATE posts_likes_dislikes SET like=0 WHERE Id=?`
		err = pr.updatePostReaction(sqlStmt, pr.db, reaction)
	}

	if err != nil {
		return fmt.Errorf("postSetLike: %w", err)
	}

	sqlStmt = `SELECT like, dislike FROM posts_likes_dislikes WHERE Id=?`
	err = pr.db.QueryRow(sqlStmt, reaction.ID).Scan(&reaction.Like, &reaction.Dislike)
	if err != nil {
		return fmt.Errorf("postSetLike: %w", err)
	}

	return nil
}

func (pr *postQuery) PostSetDislike(reaction *model.PostReaction) error {
	var sqlStmt string

	err := pr.getUserReactionToPost(reaction)
	if err != nil {
		return fmt.Errorf("postSetDislike: %w", err)
	}

	if reaction.Like == reaction.Dislike {
		sqlStmt = `UPDATE posts_likes_dislikes SET dislike=1 WHERE Id=?`
		err = pr.updatePostReaction(sqlStmt, pr.db, reaction)
	} else if reaction.Dislike == 0 {
		sqlStmt = `UPDATE posts_likes_dislikes SET like=0, dislike=1 WHERE Id=?`
		err = pr.updatePostReaction(sqlStmt, pr.db, reaction)
	} else {
		sqlStmt = `UPDATE posts_likes_dislikes SET dislike=0 WHERE Id=?`
		err = pr.updatePostReaction(sqlStmt, pr.db, reaction)
	}

	if err != nil {
		return fmt.Errorf("postSetDislike: %w", err)
	}

	sqlStmt = `SELECT like, dislike FROM posts_likes_dislikes WHERE Id=?`
	err = pr.db.QueryRow(sqlStmt, reaction.ID).Scan(&reaction.Like, &reaction.Dislike)
	if err != nil {
		return fmt.Errorf("postSetDislike: %w", err)
	}

	return nil
}

func (pr *postQuery) createReactionToPost(reaction *model.PostReaction) error {
	sqlStmt := `INSERT INTO posts_likes_dislikes (post_id, user_id, like, dislike) 
	SELECT post_id, ?, 0, 0
	FROM posts
	WHERE EXISTS (SELECT * FROM posts WHERE post_id=?)`
	query, err := pr.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("createReactionToPost: %w", err)
	}

	defer query.Close()

	res, err := query.Exec(reaction.User.ID, reaction.Post.ID)
	if err != nil {
		return fmt.Errorf("createReactionToPost: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("createReactionToPost: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("createReactionToPost: %w", model.ErrPostNotFound)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("createReactionToPost: %w", err)
	}

	reaction.ID = id
	return nil
}

func (pr *postQuery) updatePostReaction(sqlStmt string, db *sql.DB, reaction *model.PostReaction) error {
	result, err := db.Exec(sqlStmt, reaction.ID)
	if err != nil {
		return fmt.Errorf("updatePostReaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("updatePostReaction: %w", model.ErrUpdateFailed)
	}

	return nil
}

func (pr *postQuery) getUserReactionToPost(reaction *model.PostReaction) error {
	sqlStmt := `SELECT id, like, dislike FROM posts_likes_dislikes WHERE post_id=? and user_id=?`
	query, err := pr.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("getUserReactionToPost: %w", err)
	}

	defer query.Close()

	err = query.QueryRow(reaction.Post.ID, reaction.User.ID).Scan(&reaction.ID, &reaction.Like, &reaction.Dislike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pr.createReactionToPost(reaction)
		}

		return fmt.Errorf("getUserReactionToPost: %w", err)
	}
	return nil
}

func (p *postQuery) GetPostLikesDislikes(post *model.Post) error {
	sqlStmt := `SELECT COALESCE(SUM(like), 0), COALESCE(SUM(dislike), 0) FROM posts_likes_dislikes WHERE post_id=?`
	query, err := p.db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("getPostLikesDislikes: %w", err)
	}

	defer query.Close()

	err = query.QueryRow(post.ID).Scan(&post.Like, &post.Dislike)
	if err != nil {
		return fmt.Errorf("getPostLikesDislikes: %w", err)
	}
	return nil
}

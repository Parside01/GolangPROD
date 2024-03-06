package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	createReactionsTable = `CREATE TABLE IF NOT EXISTS reactions (
		post_id VARCHAR(255) NOT NULL,
		user_id VARCHAR(255) NOT NULL,
		reaction_type INT NOT NULL,
		createdAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT post_id_fk FOREIGN KEY (post_id) REFERENCES posts(id),
		CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users(id));`
	addReactionsql = `INSERT INTO reactions (post_id, user_id, reaction_type) VALUES ($1, $2, $3)`

	getLastUserPostReaction = `SELECT reaction_type
								FROM reactions
								WHERE user_id = $1
								AND post_id = $2
								ORDER BY createdAt DESC
								LIMIT 1;`
	updateUserReaction = `UPDATE reactions
	SET reaction_type = $1
	WHERE post_id = $2
	AND user_id = $3;`
	setPostLike = `UPDATE posts
	SET likesCount = likesCount + 1
	WHERE id = $1;`
	setPostDislike = `UPDATE posts
	SET dislikesCount = dislikesCount + 1
	WHERE id = $1;`
	incPostLike      = `UPDATE posts SET likesCount = likesCount - 1 WHERE id = $2;`
	incPostDislike   = `UPDATE posts SET dislikesCount = dislikesCount - 1 WHERE id = $2;`
	getreactioncount = `SELECT COUNT(*) FROM reactions WHERE post_id = $1 AND reaction_type = $2;`

	getDislike = `SELECT dislikesCount FROM posts WHERE id = $1`
	getLikes   = `SELECT likesCount FROM posts WHERE id = $1`
)

const (
	reactionLike    = 1
	reactionDislike = -1
)

type ReactionType int

const (
	reaction_dislike ReactionType = reactionLike
	reaction_like    ReactionType = reactionDislike
)

func (s *PostgresDB) setupReactionTable() error {
	_, err := s.db.Exec(createReactionsTable)
	return err
}

func (s *PostgresDB) addReaction(ctx context.Context, postID, userID string, reactionType ReactionType) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var existingReaction ReactionType
	err = tx.QueryRowContext(ctx, getLastUserPostReaction, userID, postID).Scan(&existingReaction)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to get existing reaction: %w", err)
	}

	if existingReaction == reactionType {
		return nil // User already reacted with the same type
	}

	_, err = tx.ExecContext(ctx, updateUserReaction, reactionType, postID, userID)
	if err != nil {
		return fmt.Errorf("failed to update reaction: %w", err)
	}

	var updateLikes bool
	switch reactionType {
	case reaction_like:
		updateLikes = true
	case reaction_dislike:
		updateLikes = false
	default:
		return errors.New("invalid reaction type")
	}

	var updateStmt *sql.Stmt
	if updateLikes {
		updateStmt, err = tx.PrepareContext(ctx, setPostLike)
	} else {
		updateStmt, err = tx.PrepareContext(ctx, setPostDislike)
	}
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer updateStmt.Close()

	_, err = updateStmt.ExecContext(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to update post count: %w", err)
	}

	return nil
}

func (s *PostgresDB) getLikesCount(postID string) (int, error) {
	var count int
	err := s.db.QueryRow(getLikes, postID).Scan(&count)
	return count, err
}

func (s *PostgresDB) getDislikesCount(postID string) (int, error) {
	var count int
	err := s.db.QueryRow(getDislike, postID).Scan(&count)
	return count, err
}

package store

import (
	"context"
	"database/sql"
)

type Follower struct {
	FollowerID     int64  `json:"follower_id"`
	FollowedUserID int64  `json:"followed_id"`
	CreatedAt      string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, FollowedUserID, FollowerID int64) error {
	query := `
		INSERT INTO followers (followed_id, follower_id)
		VALUES ($1, $2)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, FollowedUserID, FollowerID)
	return err
}

func (s *FollowerStore) UnFollow(ctx context.Context, FollowedUserID, FollowerID int64) error {
	query := `
		DELETE FROM followers 
		WHERE follower_id = $1 AND followed_id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, FollowerID, FollowedUserID)
	return err
}


func (s *FollowerStore) FollowerCount(ctx context.Context, userID int64) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM followers
		WHERE followed_id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var count int64
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
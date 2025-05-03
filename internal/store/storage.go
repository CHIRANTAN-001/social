package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrPostNotFound      = errors.New("not found")
	QueryTimeoutDuration = 5 * time.Second
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
		GetUserFeed(context.Context, int64,PaginationFeedQuery) ([]PostWithMetaData, error)
	}
	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, int64) (*User, error)
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
	}
	Followers interface {
		Follow(ctx context.Context, FollowedUserID, FollowerID int64) error
		UnFollow(ctx context.Context, FollowedUserID, FollowerID int64) error
		FollowerCount(ctx context.Context, userID int64) (int64, error)
	}
}

// NewStorage initializes and returns a new Storage instance with connected
// PostStore, UserStore, and LikesStore using the provided database connection.

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}

package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/CHIRANTAN-001/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user, "User fetched successfully"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type FollowerUser struct {
	FollowerUserID int64 `json:"follower_id"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserFromCtx(r)

	payload := &FollowerUser{}
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if followedUser.ID == payload.FollowerUserID {
		app.badRequestResponse(w, r, errors.New("cannot follow yourself"))
		return
	}

	if err := app.store.Followers.Follow(ctx, followedUser.ID, payload.FollowerUserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil, "User followed successfully"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowedUser := getUserFromCtx(r)

	payload := &FollowerUser{}
	if err := readJSON(w, r, payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if unfollowedUser.ID == payload.FollowerUserID {
		app.badRequestResponse(w, r, errors.New("cannot unfollow yourself"))
		return
	}

	if err := app.store.Followers.UnFollow(ctx, unfollowedUser.ID, payload.FollowerUserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil, "User unfollowed successfully"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getUserFollowersCountHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	followers, err := app.store.Followers.FollowerCount(r.Context(), user.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// followerCount := struct {
	// 	FollowersCount int64 `json:"followers_count"`
	// }{
	// 	FollowersCount: followers,
	// }

	followerCount := map[string] int64{
		"followers_count": followers,
	}

	if err := app.jsonResponse(w, http.StatusOK, followerCount, "Followers fetched successfully"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, userId)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrPostNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}

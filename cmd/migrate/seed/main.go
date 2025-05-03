package main

import (
	"github.com/CHIRANTAN-001/social/internal/db"
	"github.com/CHIRANTAN-001/social/internal/store"
)

func main() {
	conn, err := db.New("postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable", 25, 25, "15m")
	if err != nil {
		panic(err)
	}

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}

package main

import (
	"context"
	"fmt"

	"github.com/MikkelvtK/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) commandHandler {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error fetching user: %v", err)
		}

		return handler(s, cmd, user)
	}
}

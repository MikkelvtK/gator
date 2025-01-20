package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MikkelvtK/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login: username is required")
	}

	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		return err
	}

	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return err
	}

	fmt.Printf("new username has been set: %s\n", s.cfg.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("username is required")
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      cmd.args[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("new username has been registered: %s\n", user.Name)
	return nil
}

func handlerReset(s *state, _ command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting users: %v", err)
	}

	fmt.Println("the database was successfully reset")
	return nil
}

func handlerUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching users: %v", err)
	}

	for _, u := range users {
		row := fmt.Sprintf("* %s", u.Name)
		if u.Name == s.cfg.CurrentUserName {
			row = fmt.Sprintf("* %s (current)", u.Name)
		}

		fmt.Println(row)
	}

	return nil
}

func handlerAgg(_ *state, _ command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

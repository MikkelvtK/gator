package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/MikkelvtK/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return numArgError(1, "username")
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
		return numArgError(1, "username")
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

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return numArgError(1, "time_between_reqs")
	}

	interval, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		if err = scrapeFeeds(s); err != nil {
			return err
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return numArgError(2, "name", "url")
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("error creating feed: %v", err)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("error creating feed follow: %v", err)
	}

	fmt.Printf("%s created the feed: %s\n", user.Name, feed.Name)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("no feeds found")
		return nil
	}

	fmt.Printf("Found %d feeds\n", len(feeds))

	for i, f := range feeds {
		fmt.Printf("Feed %d\n", i+1)
		fmt.Printf("  ID: %s\n", f.ID)
		fmt.Printf("  Feed Name: %s\n", f.FeedName)
		fmt.Printf("  User Name: %s\n", f.UserName)
		fmt.Printf("  URL: %s\n", f.Url)
		fmt.Printf("  Created At: %s\n", f.CreatedAt)
		fmt.Printf("  Updated At: %s\n\n", f.UpdatedAt)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected 1 argument: url")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("error creating feed follow: %v", err)
	}

	fmt.Printf("%s followed the feed: %s\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, _ command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error fetching feeds: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("no feeds found")
		return nil
	}

	fmt.Printf("found %d feeds for %s\n", len(feeds), user.Name)

	for _, f := range feeds {
		fmt.Printf("* %s\n", f.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return numArgError(1, "url")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	deleteFeedFollowParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), deleteFeedFollowParams)
	if err != nil {
		return fmt.Errorf("error deleting feed follow: %v", err)
	}

	fmt.Printf("successfully unfollowed %s\n", cmd.args[0])
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2

	if len(cmd.args) == 1 {
		arg, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}

		limit = arg
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}

	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}

	for _, p := range posts {
		fmt.Printf("Title: %s\n", p.Title)
		fmt.Println("==============================================")
		fmt.Printf("%s\n\n", p.Description)
	}

	return nil
}

package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/MikkelvtK/gator/internal/config"
	"github.com/MikkelvtK/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: gator [command]\n")
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v\n", err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("err opening db connection: %v\n", err)
	}

	s := &state{
		cfg: cfg,
		db:  database.New(db),
	}

	cmds := &commands{
		commands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(s, cmd)
	if err != nil {
		log.Fatalf("error running command: %v\n", err)
	}
}

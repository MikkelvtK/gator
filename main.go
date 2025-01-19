package main

import (
	"log"
	"os"

	"github.com/MikkelvtK/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: gator [command]\n")
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v\n", err)
	}

	s := &state{
		cfg: cfg,
	}

	cmds := &commands{
		commands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(s, cmd)
	if err != nil {
		log.Fatalf("error running command: %v\n", err)
	}
}

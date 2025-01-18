package main

import (
	"fmt"
	"os"

	"github.com/MikkelvtK/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config: %v", err)
		os.Exit(1)
	}

	err = cfg.SetUser("Michiel")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(cfg)
}

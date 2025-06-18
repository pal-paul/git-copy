package main

import (
	"log"

	"github.com/pal-paul/git-copy/internal/gitcopy"
)

func main() {
	// Initialize environment variables
	if err := gitcopy.InitializeEnvironment(); err != nil {
		log.Fatal(err)
	}

	gitcopy.RunApplication()
}

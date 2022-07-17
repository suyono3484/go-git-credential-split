package main

import (
	"github.com/suyono3484/go-git-credential-split/internal/split"
	"os"
)

func main() {
	var (
		err error
	)

	if err = split.Split(); err != nil {
		os.Exit(2)
	}
}

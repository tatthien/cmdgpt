package main

import (
	"log"
	"os"

	"github.com/tatthien/cmdgpt/cmd"
)

func main() {
	if err := cmd.Exec(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

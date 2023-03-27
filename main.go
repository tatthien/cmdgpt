package main

import (
	"log"

	"github.com/tatthien/cmdgpt/cmd"
)

func main() {
	if err := cmd.Exec(); err != nil {
		log.Fatal(err)
	}
}

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
	// c := exec.Command("du", "-sh")
	// args := []string{"-sh", "."}
	// fmt.Println(c.String())
	// c.Args = append(c.Args, args...)
	// o, _ := c.Output()
	// fmt.Println(string(o))
}

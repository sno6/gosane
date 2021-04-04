package main

import (
	"log"

	cmd "github.com/sno6/gosane/cmd/gosane"
)

func main() {
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

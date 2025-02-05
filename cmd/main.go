package main

import (
	"fmt"
	"os"

	"github.com/TezzBhandari/mgs/pkg/server"
)

func main() {
	server := server.NewServer(":7070")
	if err := server.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

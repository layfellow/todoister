package main

import (
	"log"

	"github.com/layfellow/todoister/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, "./doc")
	if err != nil {
		log.Fatal(err)
	}
}

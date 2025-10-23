package main

import (
	"fmt"
	"os"

	"example.com/templates"
)

func main() {
	templates.InitTemplates()
	var page = templates.PageInh1
	wr := os.Stdout

	if err := page.Execute(wr, nil); err != nil {
		fmt.Print(err)
	}
}

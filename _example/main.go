package main

import (
	"fmt"
	"os"

	"example.com/templates"
)

func main() {
	tmpls := templates.New(nil)
	var page = templates.PagePag1
	w := os.Stdout

	if err := tmpls.Execute(w, page, nil); err != nil {
		fmt.Print(err)
	}
}

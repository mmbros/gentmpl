package main

/*
Writing Go Applications with Reusable Logic
https://npf.io/2016/10/reusable-commands/

https://github.com/natefinch/gorram
*/

import (
	"os"

	"github.com/mmbros/gentmpl/cli"
)

func main() {
	os.Exit(cli.Run())
}

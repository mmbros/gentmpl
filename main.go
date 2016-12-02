// gentmpl
package main

/*
Writing Go Applications with Reusable Logic
https://npf.io/2016/10/reusable-commands/

https://github.com/natefinch/gorram
*/

import (
	"fmt"
	"os"

	"github.com/mmbros/gentmpl/cli"
)

func main() {
	cfg, err := gentmpl.LoadConfig()
	if err == nil {
		err = cfg.WriteModule()
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

/*func main{
    os.Exit(cli.Run())
}*/

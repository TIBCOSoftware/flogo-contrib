package main

import (
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-contrib/trigger/cli"
)

func main() {
	result, _ := cli.Invoke()
	fmt.Fprintf(os.Stdout, "%s", result)
	os.Exit(0)
}

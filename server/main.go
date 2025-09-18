package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"

	"goi_example/server/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

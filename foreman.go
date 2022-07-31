package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.NewFlagSet("start", flag.ExitOnError)

	if len(os.Args) < 2 {
		// To-Do
		fmt.Println("usage")
		os.Exit(1)
	}
	
	switch SubCommand(os.Args[1]) {
	case startCmd:
		foreman := initForeman()
		foreman.runServices()
	default:
		// To-Do
		fmt.Println("usage")
		os.Exit(1)
	}
}
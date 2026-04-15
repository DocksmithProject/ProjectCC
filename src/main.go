package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: docksmith <command> [arguments]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "build":
		runBuild()
	case "run":
		runContainer()
	case "images":
		fmt.Println("Listing images...")
	case "rmi":
		fmt.Println("Removing image...")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
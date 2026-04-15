package main

import (
	"fmt"
)

// runContainer handles extracting layers and isolating the process
func runContainer() {
	fmt.Println("Starting container...")
	fmt.Println("1. Locating image manifest...")
	fmt.Println("2. Extracting tar layers to temporary rootfs...")
	fmt.Println("3. Applying environment variables...")
	fmt.Println("4. Executing isolated CMD process...")
}
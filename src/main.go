package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Make sure the user typed a command
	if len(os.Args) < 2 {
		fmt.Println("Usage: docksmith <command> [args]")
		return
	}

	command := os.Args[1]

	// The "Traffic Cop" router
	switch command {
	case "build":
		runBuild() // Calls the function in build.go
	case "run":
		runContainer() // Calls the function in run.go
	case "rmi":
		// Make sure they typed an image name after 'rmi'
		if len(os.Args) < 3 {
			fmt.Println("❌ Error: Please provide an image name (e.g., ./docksmith rmi my-sample-app)")
			return
		}
		removeImage(os.Args[2])
	case "images":
		listImages()
	default:
		fmt.Printf("❌ Unknown command: %s\n", command)
	}
}

// removeImage deletes the JSON manifest from the ~/.docksmith/images folder
func removeImage(imageName string) {
	homeDir, _ := os.UserHomeDir()
	manifestPath := filepath.Join(homeDir, ".docksmith", "images", imageName+".json")

	// Check if the image actually exists first
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		fmt.Printf("❌ Error: Image '%s' not found in your system.\n", imageName)
		return
	}

	// Delete the file
	err := os.Remove(manifestPath)
	if err != nil {
		fmt.Printf("❌ Error deleting image: %v\n", err)
		return
	}

	fmt.Printf("🗑️  Successfully removed image: %s\n", imageName)
}

func listImages() {
	homeDir, _ := os.UserHomeDir()
	imagesPath := filepath.Join(homeDir, ".docksmith", "images")

	files, err := os.ReadDir(imagesPath)
	if err != nil {
		fmt.Println("❌ Error reading images directory:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No images found.")
		return
	}

	fmt.Println("📦 Available Images:")
	for _, file := range files {
		fmt.Println("-", file.Name())
	}
}
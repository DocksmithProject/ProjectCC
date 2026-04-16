package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func runContainer() {
	if len(os.Args) < 3 {
		fmt.Println("❌ Error: Please provide an image name (e.g., ./docksmith run my-sample-app)")
		return
	}

	imageName := os.Args[2]
	customEnv := ""

	// Check if the user passed a -e flag to change the env variable
	// Example: ./docksmith run my-sample-app -e USER_NAME=Professor
	if len(os.Args) >= 5 && os.Args[3] == "-e" {
		customEnv = os.Args[4]
	}

	homeDir, _ := os.UserHomeDir()
	manifestPath := filepath.Join(homeDir, ".docksmith", "images", imageName+".json")

	// 1. Read the image manifest
	manifestBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		fmt.Printf("❌ Error: Could not find image '%s'. Did you build it?\n", imageName)
		return
	}

	var manifest ImageManifest
	json.Unmarshal(manifestBytes, &manifest)

	fmt.Printf("🚀 Starting container from image: %s\n", manifest.Name)

	// 2. Setup Environment Variables
	envVars := manifest.Env
	if customEnv != "" {
		fmt.Printf("🔧 Changing Env Variable at runtime: %s\n", customEnv)
		// Add the new variable (it will override the old one in the process)
		envVars = append(envVars, customEnv) 
	} else {
		fmt.Println("ℹ️  Using default environment variables from Docksmithfile.")
	}

	fmt.Printf("🌍 Container Environment: %v\n", envVars)
	fmt.Printf("⚙️  Executing Command: %v\n", manifest.Cmd)

	// 3. Run the Container Command
	if len(manifest.Cmd) > 0 {
		cmd := exec.Command(manifest.Cmd[0], manifest.Cmd[1:]...)
		cmd.Dir = "./sample-app" // Run it inside the sample-app folder
		
		// Inject the environment variables into the running process
		cmd.Env = append(os.Environ(), envVars...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Println("\n--- 📦 Container Output ---")
		err := cmd.Run()
		if err != nil {
			fmt.Println("\n❌ Container execution failed (Do you have Python installed?):", err)
		}
		fmt.Println("---------------------------")
	}
}
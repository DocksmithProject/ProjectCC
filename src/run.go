package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runContainer() {
	if len(os.Args) < 3 {
		fmt.Println("❌ Error: Please provide an image name (e.g., ./docksmith run my-sample-app)")
		return
	}

	var imageName string
	var customEnv string

	// Supports:
	// ./docksmith run my-sample-app
	// ./docksmith run my-sample-app -e USER_NAME=Professor
	// ./docksmith run -e USER_NAME=Professor my-sample-app
	if os.Args[2] == "-e" {
		if len(os.Args) < 5 {
			fmt.Println("❌ Error: Usage: ./docksmith run -e KEY=VALUE <image-name>")
			return
		}
		customEnv = os.Args[3]
		imageName = os.Args[4]
	} else {
		imageName = os.Args[2]
		if len(os.Args) >= 5 && os.Args[3] == "-e" {
			customEnv = os.Args[4]
		}
	}

	// Strip tag if user passes image:tag
	if strings.Contains(imageName, ":") {
		imageName = strings.Split(imageName, ":")[0]
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
	err = json.Unmarshal(manifestBytes, &manifest)
	if err != nil {
		fmt.Println("❌ Error: Failed to parse image manifest:", err)
		return
	}

	fmt.Printf("🚀 Starting container from image: %s\n", manifest.Name)

	// 2. Setup Environment Variables
	envVars := append([]string{}, manifest.Env...)

	if customEnv != "" {
		fmt.Printf("🔧 Changing Env Variable at runtime: %s\n", customEnv)

		overrideKey := strings.SplitN(customEnv, "=", 2)[0]
		replaced := false

		for i, env := range envVars {
			key := strings.SplitN(env, "=", 2)[0]
			if key == overrideKey {
				envVars[i] = customEnv
				replaced = true
				break
			}
		}

		if !replaced {
			envVars = append(envVars, customEnv)
		}
	} else {
		fmt.Println("ℹ️  Using default environment variables from Docksmithfile.")
	}

	fmt.Printf("🌍 Container Environment: %v\n", envVars)
	fmt.Printf("⚙️  Executing Command: %v\n", manifest.Cmd)

	// 3. Run the Container Command
	if len(manifest.Cmd) == 0 {
		fmt.Println("❌ Error: No command found in image manifest.")
		return
	}

	cmd := exec.Command(manifest.Cmd[0], manifest.Cmd[1:]...)
	cmd.Dir = "./sample-app" // Run it inside the sample-app folder

	// Inject environment variables into the running process
	cmd.Env = append(os.Environ(), envVars...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("\n--- 📦 Container Output ---")
	err = cmd.Run()
	if err != nil {
		fmt.Println("\n❌ Container execution failed (Do you have Python installed?):", err)
	}
	fmt.Println("---------------------------")
}
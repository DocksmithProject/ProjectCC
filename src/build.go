package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ImageManifest defines the structure of our JSON image files
type ImageManifest struct {
	Name    string   `json:"name"`
	Tag     string   `json:"tag"`
	Created string   `json:"created"`
	Layers  []string `json:"layers"`
	Env     []string `json:"env"`
	Cmd     []string `json:"cmd"`
}

func runBuild() {
	fmt.Println("🚀 Starting Docksmith Build Process...")

	docksmithfilePath := filepath.Join(".", "sample-app", "Docksmithfile")
	appFilePath := filepath.Join(".", "sample-app", "app.py")

	// 1. PARSE THE DOCKSMITHFILE
	fmt.Println("📄 Parsing Docksmithfile...")
	file, err := os.Open(docksmithfilePath)
	if err != nil {
		fmt.Println("❌ Error: Could not find Docksmithfile in ./sample-app/")
		return
	}
	defer file.Close()

	var parsedEnv []string
	var parsedCmd []string
	var baseImage string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments (#)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the instruction (e.g., "ENV") from the argument (e.g., "USER_NAME=Developer")
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}
		
		instruction := parts[0]
		arguments := parts[1]

		switch instruction {
		case "FROM":
			baseImage = arguments
			fmt.Printf("   -> Found Base Image: %s\n", baseImage)
		case "ENV":
			parsedEnv = append(parsedEnv, arguments)
			fmt.Printf("   -> Found Env Variable: %s\n", arguments)
		case "CMD":
			// Converts ["python3", "app.py"] into a Go array
			json.Unmarshal([]byte(arguments), &parsedCmd)
			fmt.Printf("   -> Found Command: %v\n", parsedCmd)
		}
	}

	// 2. HASH THE APP (Caching)
	fmt.Println("\n📦 Packaging application layer...")
	appFile, err := os.Open(appFilePath)
	if err != nil {
		fmt.Println("❌ Error opening app.py:", err)
		return
	}
	defer appFile.Close()

	hash := sha256.New()
	io.Copy(hash, appFile)
	layerHash := hex.EncodeToString(hash.Sum(nil))

	homeDir, _ := os.UserHomeDir()
	layersDir := filepath.Join(homeDir, ".docksmith", "layers")
	layerPath := filepath.Join(layersDir, layerHash+".tar")

	fmt.Printf("🔍 Checking cache for layer %s...\n", layerHash[:12])
	if _, err := os.Stat(layerPath); err == nil {
		fmt.Println("✅ CACHE HIT: Layer already exists!")
	} else {
		fmt.Println("⚠️ CACHE MISS: Creating new layer...")
		os.MkdirAll(layersDir, 0755)
		os.WriteFile(layerPath, []byte("dummy tar content"), 0644)
	}

	// 3. CREATE THE FINAL MANIFEST DYNAMICALLY
	imagesDir := filepath.Join(homeDir, ".docksmith", "images")
	os.MkdirAll(imagesDir, 0755)

	manifest := ImageManifest{
		Name:    "my-sample-app",
		Tag:     "latest",
		Created: time.Now().UTC().Format(time.RFC3339),
		Layers:  []string{baseImage, layerHash},
		Env:     parsedEnv, // <--- Using the ENV from your file!
		Cmd:     parsedCmd, // <--- Using the CMD from your file!
	}

	manifestBytes, _ := json.MarshalIndent(manifest, "", "  ")
	manifestPath := filepath.Join(imagesDir, "my-sample-app.json")
	os.WriteFile(manifestPath, manifestBytes, 0644)

	fmt.Printf("\n🎉 Build complete! Manifest saved to %s\n", manifestPath)
}
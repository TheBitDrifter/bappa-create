package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed templates
var templateFS embed.FS

// Define file types that should be copied as-is (not processed for replacements)
var binaryExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".wav":  true,
	".mp3":  true,
	".ogg":  true,
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	projectName := os.Args[1]
	template := "topdown" // Default template

	// Check if user specified a template
	if len(os.Args) >= 4 && os.Args[2] == "--template" {
		template = os.Args[3]
	}

	createProject(projectName, template)
}

func printUsage() {
	fmt.Println("Bappa Game Template Generator")
	fmt.Println("===============================")
	fmt.Println("Usage: bappa-create username/project-name [--template <template-name>]")
	fmt.Println()
	fmt.Println("This tool creates a new Bappa game project with the specified name.")
	fmt.Println("The username/ prefix is used to create the proper Go module path.")
	fmt.Println("Example: bappa-create johndoe/my-awesome-game --template platformer")
	fmt.Println()
	fmt.Println("Available templates:")
	fmt.Println("  platformer      - A simple platformer game")
	fmt.Println("  platformer-split - A platformer game with split-screen co-op support")
	fmt.Println("  topdown         - A top-down perspective game (default)")
	fmt.Println("  topdown-split   - A top-down game with split-screen co-op support")
	fmt.Println("  sandbox         - An open sandbox game environment")
}

func createProject(projectName, templateName string) {
	// Check if the project name contains a username
	parts := strings.Split(projectName, "/")
	if len(parts) != 2 {
		fmt.Println("Error: Project name must be in the format 'username/project-name'")
		os.Exit(1)
	}

	username := parts[0]
	projectNameOnly := parts[1]

	// Clean the project name for Go module naming
	moduleName := strings.ToLower(strings.ReplaceAll(projectNameOnly, " ", "-"))

	// Create project directory (just the project part, not with username)
	err := os.MkdirAll(projectNameOnly, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Creating new Bappa game project: %s (using %s template)\n", projectNameOnly, templateName)

	// The template import path that will be replaced
	templateImportPath := fmt.Sprintf("github.com/TheBitDrifter/bappacreate/templates/%s", templateName)

	// Make sure the template exists
	templatePath := filepath.Join("templates", templateName)
	_, err = fs.Stat(templateFS, templatePath)
	if err != nil {
		fmt.Printf("Error: Template '%s' not found\n", templateName)
		os.Exit(1)
	}

	// Create additional directories that might not be in the templates
	os.MkdirAll(filepath.Join(projectNameOnly, "assets", "images"), 0755)
	os.MkdirAll(filepath.Join(projectNameOnly, "assets", "sounds"), 0755)

	// Create additional directories for split-screen co-op templates
	if strings.HasSuffix(templateName, "-split") {
		os.MkdirAll(filepath.Join(projectNameOnly, "players"), 0755)
		os.MkdirAll(filepath.Join(projectNameOnly, "splitscreen"), 0755)
	}

	// Process all template files
	err = fs.WalkDir(templateFS, templatePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root template directory itself
		if path == templatePath {
			return nil
		}

		// Create relative path from the template directory
		relPath, err := filepath.Rel(templatePath, path)
		if err != nil {
			return err
		}

		// Create the target path in the new project using projectNameOnly
		targetPath := filepath.Join(projectNameOnly, relPath)

		// Create directories
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Handle the file appropriately
		ext := strings.ToLower(filepath.Ext(path))

		if binaryExtensions[ext] {
			// Copy binary files directly
			return copyBinaryFile(path, targetPath)
		} else {
			// Process text files with replacement
			return processTextFile(path, targetPath, username, moduleName, templateImportPath)
		}
	})
	if err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		os.Exit(1)
	}

	// Calculate the full module path with username
	modulePath := fmt.Sprintf("github.com/%s/%s", username, moduleName)

	// Initialize Go module and dependencies
	initGoModule(projectNameOnly, modulePath)

	fmt.Printf("\nSuccessfully created Bappa game project: %s\n", projectNameOnly)
	fmt.Printf("\nTo run your game:\n")
	fmt.Printf("  cd %s\n", projectNameOnly)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  go run .\n")
}

func copyBinaryFile(sourcePath, targetPath string) error {
	fmt.Println("Copying asset:", targetPath)

	// Read binary file
	data, err := templateFS.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	// Write to target path
	return os.WriteFile(targetPath, data, 0644)
}

func processTextFile(sourcePath, targetPath, username, moduleName, templateImportPath string) error {
	fmt.Println("Creating:", targetPath)

	// Calculate the new import path including username
	newImportPath := fmt.Sprintf("github.com/%s/%s", username, moduleName)

	// Read file content
	content, err := templateFS.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	fileContent := string(content)

	// Replace the full template import path with the new module path including username
	fileContent = strings.ReplaceAll(fileContent, templateImportPath, newImportPath)

	// Special handling for go.mod file
	if filepath.Base(targetPath) == "go.mod" {
		lines := strings.Split(fileContent, "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, "module ") {
				lines[i] = "module " + newImportPath
				break
			}
		}
		fileContent = strings.Join(lines, "\n")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	// Write processed content
	return os.WriteFile(targetPath, []byte(fileContent), 0644)
}

func initGoModule(projectName, modulePath string) {
	fmt.Println("\nSetting up Go module...")

	// Change to project directory
	currentDir, _ := os.Getwd()
	projectDir := filepath.Join(currentDir, projectName)
	err := os.Chdir(projectDir)
	if err != nil {
		fmt.Printf("Error changing to project directory: %v\n", err)
		return
	}
	defer os.Chdir(currentDir)

	// Initialize go.mod with the correct module name
	runCommand("go", "mod", "init", modulePath)

	// Get required dependencies
	runCommand("go", "get", "github.com/TheBitDrifter/coldbrew@latest")
	runCommand("go", "get", "github.com/TheBitDrifter/blueprint@latest")
	runCommand("go", "get", "github.com/TheBitDrifter/warehouse@latest")

	// The go mod tidy step is removed to let the user handle it themselves
}

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running: %s %s\n", command, strings.Join(args, " "))

	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: Command failed: %v\n", err)
	}
}

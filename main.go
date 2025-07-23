package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
)

func getRoot() error {
	if os.Geteuid() != 0 {
		fmt.Println("Required root access to install crates.")
		if runtime.GOOS == "windows" {
			fmt.Println("Can't rerun with admin privileges on Windows, exiting...")
		} else {
			fmt.Println("Rerunning with sudo...")
			cmd := exec.Command("sudo", os.Args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			err := cmd.Run()
			if err != nil {
				return fmt.Errorf("error rerunning with sudo: %v", err)
			}
			os.Exit(0)
		}
	}
	return nil
}

type Crate struct {
	ProjectName string
	BinaryName  string
	BinaryFile  []byte
	SourceURL   string
}

func (c Crate) Save(savePath string) error {
	if _, err := os.Stat(savePath); !os.IsNotExist(err) {
		return fmt.Errorf("file already exists at %s", savePath)
	}
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	crateBytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal crate: %w", err)
	}
	_, err = file.Write(crateBytes)
	if err != nil {
		return fmt.Errorf("failed to write crate to file: %w", err)
	}

	return nil
}

func LoadCrate(filePath string) (Crate, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Crate{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var crate Crate
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&crate)
	if err != nil {
		return Crate{}, fmt.Errorf("failed to decode crate: %w", err)
	}

	return crate, nil
}

func (c Crate) UnpackBin(prefix string) {
	var filePath string
	filePath = path.Join(prefix, c.BinaryName)

	err := os.WriteFile(filePath, c.BinaryFile, 0755)
	if err != nil {
		log.Fatalf("Error writing binary file: %v", err)
	}
}

func (c Crate) Install() error {
	if runtime.GOOS == "windows" {
		log.Fatalln("Installing on windows not available yet")
	}
	if runtime.GOOS == "linux" {
		filePath := path.Join("/bin", c.BinaryName)
		err := os.WriteFile(filePath, c.BinaryFile, 0755)
		if err != nil {
			return fmt.Errorf("failed to write binary file: %w", err)
		}
	}
	return nil
}

func (c Crate) Uninstall() error {
	filePath := path.Join("/bin", c.BinaryName)
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to remove binary file: %w", err)
	}
	return nil
}

func buildCrate(projectName string, binFile string, sourceURL string) Crate {
	file, err := os.Open(binFile)
	if err != nil {
		log.Fatalf("Error opening binary file: %v", err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading binary file: %v", err)
	}

	crate := Crate{
		ProjectName: projectName,
		BinaryName:  path.Base(binFile),
		BinaryFile:  content,
		SourceURL:   sourceURL,
	}

	return crate
}

func main() {
	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	} else {
		command = ""
	}
	var subCommand string
	if len(os.Args) > 2 {
		subCommand = os.Args[2]
	} else {
		subCommand = ""
	}
	switch command {
	case "build", "Build", "B", "b":
		fmt.Println("crating the project...")
		projectName := ""
		binFile := ""
		sourceURL := ""

		fmt.Print("Enter the project name: ")
		_, err := fmt.Scanln(&projectName)
		if err != nil {
			log.Fatalln("Error reading project name:", err)
		}

		fmt.Print("Enter the binary file name: ")
		_, err = fmt.Scanln(&binFile)
		if err != nil {
			log.Fatalln("Error reading binary file name:", err)
		}

		fmt.Print("Enter the source URL (optional, press Enter to skip): ")
		_, err = fmt.Scanln(&sourceURL)
		if err != nil && errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			log.Fatalln("Error reading source URL:", err)
		}

		crate := buildCrate(projectName, binFile, sourceURL)
		savePath := projectName + ".json"
		err = crate.Save(savePath)
		if err != nil {
			log.Fatalln("Error saving crate:", err)
		}
	case "install", "Install", "I", "i":
		err := getRoot()
		if err != nil {
			log.Fatalf("Error getting root: %v", err)
		}
		crate, err := LoadCrate(subCommand)
		if err != nil {
			log.Fatalf("Error loading crate: %v", err)
		}
		if crate.BinaryFile == nil || len(crate.BinaryFile) == 0 {
			log.Fatalf("No binary file found in crate: %s", crate.ProjectName)
		}
		fmt.Println("Crate loaded successfully:", crate.ProjectName)
		fmt.Println("Source URL: ", crate.SourceURL)
		fmt.Println("Binary name: ", crate.BinaryName)
		err = crate.Install()
		if err != nil {
			log.Fatalf("Error installing crate: %v", err)
		}
		fmt.Println("Crate installed successfully.")
	case "uninstall", "Uninstall", "U", "u":
		err := getRoot()
		if err != nil {
			log.Fatalf("Error getting root: %v", err)
		}
		crate, err := LoadCrate(subCommand)
		if err != nil {
			log.Fatalf("Error loading crate: %v", err)
		}
		err = crate.Uninstall()
		if err != nil {
			log.Fatalf("Error uninstalling crate: %v", err)
		}
		fmt.Println("Crate uninstalled successfully.")
	case "get-bin", "Get-bin", "g", "G":
		crate, err := LoadCrate(subCommand)
		if err != nil {
			log.Fatalf("Error loading crate: %v", err)
		}
		if crate.BinaryFile == nil || len(crate.BinaryFile) == 0 {
			log.Fatalf("No binary file found in crate: %s", crate.ProjectName)
		}

		fmt.Print("Enter the prefix path to unpack the binary: ")
		var prefix string
		_, err = fmt.Scanln(&prefix)
		if err != nil {
			log.Fatalf("Error reading prefix path: %v", err)
		}
		crate.UnpackBin(prefix)
		fmt.Println("Binary unpacked successfully at:", prefix)
	case "pull", "Pull", "P", "p":
		if len(subCommand) == 0 {
			log.Fatalln("Please provide a source URL to pull the crate.")
		}
		fmt.Println("Pulling crate from source URL:", subCommand)
		resp, err := http.Get(subCommand)
		if err != nil {
			log.Fatalf("Error pulling crate from URL: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Failed to pull crate, status code: %d", resp.StatusCode)
		}
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}
		var crate Crate
		err = json.Unmarshal(respBytes, &crate)
		if err != nil {
			log.Fatalf("Error unmarshalling crate: %v", err)
		}
		fmt.Println("Crate pulled successfully:", crate.ProjectName)
		fmt.Println("Do you want to install or save the crate? (install/save)")
		var action string
		_, err = fmt.Scanln(&action)
		if err != nil {
			log.Fatalf("Error reading action: %v", err)
		}
		switch action {
		case "install":
			err = crate.Install()
			if err != nil {
				log.Fatalf("Error installing pulled crate: %v", err)
			}
		case "save":
			savePath := crate.ProjectName + ".json"
			err = crate.Save(savePath)
			if err != nil {
				log.Fatalf("Error saving pulled crate: %v", err)
			}
			fmt.Println("Crate saved successfully at:", savePath)
		}
	default:
		fmt.Println("Usage: go run main.go [command] [options]")
		fmt.Println("Commands:")
		fmt.Println("  build  |Build  |b|B             - Build a new crate")
		fmt.Println("  install|Install|i|I  <filename> - Install a crate")
		fmt.Println("  get-bin|Get-bin|g|G  <filename> - Unpack a crate binary to a specified prefix path")
		fmt.Println("  pull   |Pull   |p|P             - Pull a crate from a source URL")
	}

}

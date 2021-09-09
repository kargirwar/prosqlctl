//go:build windows
// +build windows

package adapter

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	utils "github.com/kargirwar/prosqlctl/utils"
)

func CopyAgent() {
	root := os.Getenv("programfiles")
	path := filepath.Join(root, "ProsqlAgent")
	fmt.Println("Creating " + path)
	err := os.MkdirAll(path, os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Copying agent to " + path)
	agent := filepath.Join(utils.GetCwd(), "prosql-agent", "prosql-agent.exe")

	//copy executable to /usr/local/bin
	cmd := exec.Command("cmd", "/c", "copy", agent, path)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func StartAgent() {
	root := os.Getenv("programfiles")
	agent := filepath.Join(root, "ProsqlAgent", "prosql-agent.exe")

	fmt.Printf("Installing agent.. ")
	cmd := exec.Command("nssm.exe", "install", "prosql-agent", agent)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")

	fmt.Printf("Setting app directory.. ")
	appdir := filepath.Join(root, "ProsqlAgent")
	cmd = exec.Command("nssm.exe", "set", "prosql-agent", "AppDirectory", appdir)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")

	fmt.Printf("Starting agent.. ")
	cmd = exec.Command("nssm.exe", "start", "prosql-agent")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")
}

func DelAgent() {
	//noop
}

func StopAgent() {
	cmd := exec.Command("nssm.exe", "stop", "prosql-agent")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command("nssm.exe", "remove", "prosql-agent", "confirm")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

//go:build linux
// +build linux

package adapter

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"

	utils "github.com/kargirwar/prosqlctl/utils"
)

const RELEASE_ARCHIVE = "release.zip"
const UNIT = `
[Unit]
Description=prosql-agent for prosql.io
[Install]
WantedBy=multi-user.target
[Service]
Type=simple
User={{.User}}
ExecStart={{.Program}}
WorkingDirectory={{.WorkingDir}}
Restart=always
RestartSec=5
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=%n
`

func StartAgent() {
	//create unit file and use systemctl to start agent
	fmt.Println("Creating unit ...")
	data := struct {
		User       string
		Program    string
		WorkingDir string
	}{
		Program:    "prosql-agent",
		WorkingDir: utils.GetCwd(),
		User:       os.Getenv("SUDO_USER"),
	}

	unit := fmt.Sprintf("/etc/systemd/system/prosql-agent.service")
	f, err := os.Create(unit)

	t := template.Must(template.New("unit").Parse(UNIT))
	err = t.Execute(f, data)
	if err != nil {
		log.Fatalf("Unable to create unit file: %s", err)
	}
	fmt.Println("Reloading serivces...")

	cmd := exec.Command("systemctl", "daemon-reload")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Enabling prosql-agent...")
	cmd = exec.Command("systemctl", "enable", "prosql-agent.service")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting prosql-agent...")
	cmd = exec.Command("systemctl", "start", "prosql-agent.service")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func DownloadAgent() {
	release := utils.GetLatestRelease()

	//Download and extract
	fmt.Printf("Downloading release %s .. ", release.Version)
	utils.DownloadFile(RELEASE_ARCHIVE, release.Linux)
	fmt.Println("Done.")

	fmt.Printf("Extracting files ..")
	utils.Unzip(RELEASE_ARCHIVE, utils.GetCwd())
	fmt.Println("Done.")
}

func Cleanup() {
	fmt.Printf("Cleaning up ..")

	err := os.RemoveAll("prosql-agent")
	if err != nil {
		log.Fatal(err)
	}

	//delete archive
	err = os.Remove(RELEASE_ARCHIVE)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done.")
}

func CopyAgent() {
	fmt.Println("Copying agent to /usr/local/bin ...")
	agent := utils.GetCwd() + "/prosql-agent/prosql-agent"

	//copy executable to /usr/local/bin
	cmd := exec.Command("cp", agent, "/usr/local/bin")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func DelAgent() {
	fmt.Println("Deleting agent from /usr/local/bin ...")
	//copy executable to /usr/local/bin
	cmd := exec.Command("rm", "-f", "/usr/local/bin/prosql-agent")
	err := cmd.Run()

	if err != nil {
		//can't do much about error here
		log.Println(err)
	}
}

func StopAgent() {
	fmt.Println("Stopping agent ...")
	cmd := exec.Command("systemctl", "stop", "prosql-agent.service")
	err := cmd.Run()
	if err != nil {
		//can't do much about error here
		log.Println(err)
	}

	fmt.Println("Deleting unit ...")
	cmd = exec.Command("rm", "-f", "/etc/systemd/system/prosql-agent.service")
	err = cmd.Run()

	if err != nil {
		//can't do much about error here
		log.Println(err)
	}
}

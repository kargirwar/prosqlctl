//go:build mac
// +build mac

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
const BINARY = "prosql-agent"
const LABEL = "io.prosql.agent"
const PLIST = `<?xml version='1.0' encoding='UTF-8'?>
<!DOCTYPE plist PUBLIC \"-//Apple Computer//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\" >
<plist version='1.0'>
  <dict>
    <key>Label</key><string>{{.Label}}</string>
    <key>Program</key><string>{{.Program}}</string>
	<key>StandardOutPath</key><string>/tmp/{{.Label}}.out.log</string>
	<key>StandardErrorPath</key><string>/tmp/{{.Label}}.err.log</string>
    <key>KeepAlive</key><true/>
    <key>RunAtLoad</key><true/>
  </dict>
</plist>
`
const ROOT_DIR = "~/.prosql-agent/bin"

func DownloadAgent() {
	release := utils.GetLatestRelease()

	//Download and extract
	fmt.Printf("Downloading release %s .. ", release.Version)
	utils.DownloadFile(RELEASE_ARCHIVE, release.Mac)
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

func installAgent() {
	CopyAgent()
	StartAgent()
}

func unInstallAgent() {
	DelAgent()
	StopAgent()
}

func StartAgent() {
	//create plist and use launchctrl to start service
	fmt.Println("Creating plist ...")
	data := struct {
		Label   string
		Program string
	}{
		Label:   LABEL,
		//Program: "/usr/local/prosql-agent/bin/prosql-agent",
		Program: getRootDir() + "/prosql-agent",
	}

	plist := fmt.Sprintf("%s/Library/LaunchAgents/%s.plist", os.Getenv("HOME"), data.Label)
	f, err := os.Create(plist)

	t := template.Must(template.New("launchdConfig").Parse(PLIST))
	err = t.Execute(f, data)
	if err != nil {
		log.Fatalf("Template generation failed: %s", err)
	}
	fmt.Println("Starting ...")

	cmd := exec.Command("launchctl", "load", plist)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func CopyAgent() {
	err := os.MkdirAll(getRootDir(), 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	fmt.Println("Copying agent to " + getRootDir())
	program := utils.GetCwd() + "/prosql-agent/prosql-agent"

	//copy executable to /usr/local/bin
	cmd := exec.Command("cp", "-v", program, getRootDir())
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func DelAgent() {
	fmt.Println("Deleting agent from " + getRootDir())
	//copy executable to /usr/local/bin
	cmd := exec.Command("rm", "-rf", getRootDir())
	err := cmd.Run()

	if err != nil {
		//can't do much about error here
		log.Println(err)
	}
}

func StopAgent() {
	plist := fmt.Sprintf("%s/Library/LaunchAgents/%s.plist", os.Getenv("HOME"), LABEL)

	fmt.Println("Stopping agent ...")
	cmd := exec.Command("launchctl", "unload", plist)
	err := cmd.Run()
	if err != nil {
		//can't do much about error here
		log.Println(err)
	}

	fmt.Println("Deleting plist ...")
	cmd = exec.Command("rm", "-f", plist)
	err = cmd.Run()

	if err != nil {
		//can't do much about error here
		log.Println(err)
	}
}

func getRootDir() string {
	return os.Getenv("HOME") + "/.prosql-agent"
}

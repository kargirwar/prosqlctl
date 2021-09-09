package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/kargirwar/prosqlctl/adapter"
	"github.com/kargirwar/prosqlctl/utils"
)

var VERSION = "VERSION"

const MAX_TRIES = 3

func main() {
	var install = flag.Bool("install", false, "Install agent on your system")
	var help = flag.Bool("help", false, "Show help message")
	var uninstall = flag.Bool("uninstall", false, "Uninstall prosql-agent from your system")
	var update = flag.Bool("update", false, "Update prosql-agent")
	var status = flag.Bool("status", false, "Show prosql-agent status")
	var version = flag.Bool("version", false, "Show prosqlctl version")
	flag.Parse()

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "prosqlctl usage:\n")
		flag.PrintDefaults()
	}

	if *help {
		flag.Usage()
		return
	}

	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	if *version {
		fmt.Println("prosqlctl version: " + VERSION)
		return
	}

	if *install {
		installAgent()
		return
	}

	if *uninstall {
		unInstallAgent()
		return
	}

	if *update {
		updateAgent()
		return
	}

	if *status {
		res := utils.GetStatus()
		if res.Status == "ok" {
			fmt.Println("prosql-agent is installed and running")
			return
		}

		fmt.Println(res.Msg)
	}
}

func installAgent() {
	res := utils.GetStatus()
	if res.Status == "ok" {
		fmt.Println("prosql-agent is already installed. You may want to use -update")
		return
	}

	adapter.DownloadAgent()
	adapter.CopyAgent()
	adapter.StartAgent()
	adapter.Cleanup()
	fmt.Println("Installed successfully!")
}

func unInstallAgent() {
	adapter.DelAgent()
	adapter.StopAgent()
	fmt.Println("Uninstalled prosql-agent")
}

func updateAgent() {
	res := utils.GetStatus()
	if res.Status == "ok" {
		unInstallAgent()
	}

	installAgent()

	i := 0
	for {
		//there could be a delay in starting the service
		//so we wait for some time and try again if required
		time.Sleep(1 * time.Second)

		res = utils.GetStatus()
		if res.Status == "ok" {
			version := utils.GetValue(res, "version")
			fmt.Printf("\nUpdated prosql-agent to %s\n", version)
			return
		}

		i++
		if i == MAX_TRIES {
			break
		}
	}

	fmt.Println(res.Msg)
}

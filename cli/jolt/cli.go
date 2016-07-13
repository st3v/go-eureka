package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}

func main() {
	app := cli.NewApp()
	app.Name = "jolt"
	app.Usage = "Command-line client for Netflix Eureka"

	app.Commands = []cli.Command{
		registerCmd,
		deregisterCmd,
		heartbeatCmd,
		instancesCmd,
	}

	app.Run(os.Args)
}

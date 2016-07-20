package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"

	"github.com/st3v/go-eureka"
)

var getStatus = func(c *cli.Context, required bool) eureka.Status {
	name := c.Args().First()

	if name == "" && required {
		fmt.Fprintln(c.App.Writer, "must specify status")
		os.Exit(1)
	}

	if name == "" {
		return eureka.StatusUnknown
	}

	status, err := eureka.ParseStatus(c.Args().First())
	if err != nil {
		fmt.Fprintln(c.App.Writer, err)
		os.Exit(1)
	}

	return status
}

var overrideCmd = cli.Command{
	Name:  "override-status",
	Usage: "override the status of a registered instance",

	Flags: []cli.Flag{
		instanceFlag,
		endpointsFlag,
	},

	Action: func(c *cli.Context) {
		instance := getInstance(c, "override-status")
		endpoints := getEndpoints(c, "override-status")
		status := getStatus(c, true)

		log.Printf("Overriding status for instance '%s' of application '%s'... \n", instance.Id, instance.AppName)
		client := eureka.NewClient(endpoints)
		if err := client.StatusOverride(instance, status); err != nil {
			log.Fatalf("Error overriding status: %s", err)
		}

		log.Println("Success")
	},
}

var removeOverrideCmd = cli.Command{
	Name:  "remove-override",
	Usage: "remove status override from a registered instance",

	Flags: []cli.Flag{
		instanceFlag,
		endpointsFlag,
	},

	Action: func(c *cli.Context) {
		instance := getInstance(c, "remove-override")
		endpoints := getEndpoints(c, "remove-override")
		status := getStatus(c, false)

		log.Printf("Remove overridden status from instance '%s' of application '%s'... \n", instance.Id, instance.AppName)
		client := eureka.NewClient(endpoints)
		if err := client.RemoveStatusOverride(instance, status); err != nil {
			log.Fatalf("Error overriding status: %s", err)
		}

		log.Println("Success")
	},
}

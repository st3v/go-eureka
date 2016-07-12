package main

import (
	"log"

	"github.com/codegangsta/cli"

	"github.com/st3v/jolt"
)

var deregisterCmd = cli.Command{
	Name:  "deregister",
	Usage: "deregister an instance from Eureka",

	Flags: []cli.Flag{
		instanceFlag,
		endpointsFlag,
	},

	Action: func(c *cli.Context) {
		instance := getInstance(c, "deregister")
		endpoints := getEndpoints(c, "deregister")

		log.Printf("Deregistering instance '%s' for application '%s'... \n", instance.Id, instance.App)
		client := jolt.NewClient(endpoints)
		if err := client.Deregister(instance); err != nil {
			log.Fatalf("Error deregistering instance with Eureka: %s", err)
		}

		log.Println("Success")
	},
}

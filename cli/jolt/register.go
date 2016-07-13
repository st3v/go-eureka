package main

import (
	"log"

	"github.com/codegangsta/cli"

	"github.com/st3v/jolt"
)

var registerCmd = cli.Command{
	Name:  "register",
	Usage: "register an instance with Eureka",

	Flags: []cli.Flag{
		instanceFlag,
		endpointsFlag,
	},

	Action: func(c *cli.Context) {
		instance := getInstance(c, "register")
		endpoints := getEndpoints(c, "register")

		log.Printf("Registering instance '%s' for application '%s'... \n", instance.Id, instance.AppName)
		client := jolt.NewClient(endpoints)
		if err := client.Register(instance); err != nil {
			log.Fatalf("Error registering instance with Eureka: %s", err)
		}

		log.Println("Success")
	},
}

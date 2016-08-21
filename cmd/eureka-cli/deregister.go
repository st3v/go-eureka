package main

import (
	"log"

	"github.com/codegangsta/cli"

	"github.com/st3v/go-eureka"
)

var deregisterCmd = cli.Command{
	Name:  "deregister",
	Usage: "deregister an instance from Eureka",

	Flags: []cli.Flag{
		instanceFlag,
		endpointsFlag,
	},

	Action: func(c *cli.Context) error {
		instance := getInstance(c, "deregister")
		endpoints := getEndpoints(c, "deregister")

		log.Printf("Deregistering instance '%s' for application '%s'... \n", instance.ID, instance.AppName)
		client := eureka.NewClient(endpoints)
		if err := client.Deregister(instance); err != nil {
			log.Printf("Error deregistering instance with Eureka: %s\n", err)
			return err
		}

		log.Println("Success")
		return nil
	},
}

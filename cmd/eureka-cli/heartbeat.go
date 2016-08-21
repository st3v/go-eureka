package main

import (
	"log"

	"github.com/codegangsta/cli"

	"github.com/st3v/go-eureka"
)

var heartbeatCmd = cli.Command{
	Name:  "heartbeat",
	Usage: "send heartbeat to a registered instance",

	Flags: []cli.Flag{
		instanceFlag,
		endpointsFlag,
	},

	Action: func(c *cli.Context) error {
		instance := getInstance(c, "heartbeat")
		endpoints := getEndpoints(c, "heartbeat")

		log.Printf("Sending heartbeat for instance '%s' of application '%s'... \n", instance.ID, instance.AppName)
		client := eureka.NewClient(endpoints)
		if err := client.Heartbeat(instance); err != nil {
			log.Printf("Error sending heartbeat: %s\n", err)
			return err
		}

		log.Println("Success")
		return nil
	},
}

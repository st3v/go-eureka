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

	Action: func(c *cli.Context) {
		instance := getInstance(c, "heartbeat")
		endpoints := getEndpoints(c, "heartbeat")

		log.Printf("Sending heartbeat for instance '%s' of application '%s'... \n", instance.Id, instance.AppName)
		client := eureka.NewClient(endpoints)
		if err := client.Heartbeat(instance); err != nil {
			log.Fatalf("Error sending heartbeat: %s", err)
		}

		log.Println("Success")
	},
}

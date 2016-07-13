package main

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/st3v/jolt"
)

var instancesCmd = cli.Command{
	Name:  "instances",
	Usage: "retrieve registered instances from Eureka",

	Flags: []cli.Flag{
		endpointsFlag,
		appNameFlag,
		instanceIdFlag,
	},

	Action: func(c *cli.Context) {
		endpoints := getEndpoints(c, "heartbeat")
		client := jolt.NewClient(endpoints)

		instances := make([]jolt.Instance, 1)

		appName := c.String("app")
		instanceId := c.String("instance")

		switch {
		case instanceId != "":
			if appName == "" {
				cli.ShowCommandHelp(c, "instances")
				log.Fatalln("--app flag required")
			}

			log.Printf("Retrieving instances for application '%s' and instance id '%s'...", appName, instanceId)

			instance, err := client.Instance(appName, instanceId)
			if err != nil {
				log.Fatalf("Error retrieving instances: %s", err)
			}

			instances = append(instances, instance)
		case appName != "":
			log.Printf("Retrieving instances for application '%s'...", appName)

			app, err := client.App(appName)
			if err != nil {
				log.Fatalf("Error retrieving instances: %s", err)
			}

			instances = append(instances, app.Instances...)
		default:
			log.Println("Retrieving instances for all registered applications ...")

			apps, err := client.Apps()
			if err != nil {
				log.Fatalf("Error retrieving instances: %s", err)
			}

			if len(apps) == 0 {
				fmt.Println("No registered apps")
				return
			}

			for _, app := range apps {
				instances = append(instances, app.Instances...)
			}
		}

		for _, i := range instances {
			if data, err := xml.MarshalIndent(i, "", "  "); err == nil {
				log.Println(string(data))
			}
		}
	},
}

package main

import (
	"encoding/xml"
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

		var instances []jolt.Instance

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
				log.Fatalf("Error retrieving instance: %s", err)
			}

			instances = append(instances, instance)
		case appName != "":
			log.Printf("Retrieving instances for application '%s'...", appName)

			app, err := client.App(appName)
			if err != nil {
				log.Fatalf("Error retrieving application: %s", err)
			}

			instances = append(instances, app.Instances...)
		default:
			log.Println("Retrieving instances for all registered applications ...")

			apps, err := client.Apps()
			if err != nil {
				log.Fatalf("Error retrieving applications: ", err)
			}

			for _, app := range apps {
				instances = append(instances, app.Instances...)
			}
		}

		output := struct {
			XMLName   xml.Name        `xml:"instances"`
			Instances []jolt.Instance `xml:"instance"`
		}{
			Instances: instances,
		}

		data, err := xml.MarshalIndent(output, "", "  ")
		if err != nil {
			log.Fatalf("Error rendering output: %s", err)
		}

		log.Println(string(data))
	},
}

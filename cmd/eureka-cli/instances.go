package main

import (
	"encoding/xml"
	"log"

	"github.com/codegangsta/cli"

	"github.com/st3v/go-eureka"
)

var instancesCmd = cli.Command{
	Name:  "instances",
	Usage: "retrieve registered instances from Eureka",

	Flags: []cli.Flag{
		endpointsFlag,
		appNameFlag,
		instanceIDFlag,
	},

	Action: func(c *cli.Context) error {
		endpoints := getEndpoints(c, "heartbeat")
		client := eureka.NewClient(endpoints)

		var instances []*eureka.Instance

		appName := c.String("app")
		instanceID := c.String("instance")

		switch {
		case instanceID != "" && appName != "":
			log.Printf("Retrieving instances for application '%s' and instance id '%s'...", appName, instanceID)

			instance, err := client.AppInstance(appName, instanceID)
			if err != nil {
				log.Printf("Error retrieving instance: %s\n", err)
				return err
			}

			instances = append(instances, instance)
		case instanceID != "" && appName == "":
			log.Printf("Retrieving instance with id '%s'...", instanceID)

			instance, err := client.Instance(instanceID)
			if err != nil {
				log.Printf("Error retrieving instance: %s\n", err)
				return err
			}

			instances = append(instances, instance)
		case instanceID == "" && appName != "":
			log.Printf("Retrieving instances for application '%s'...", appName)

			app, err := client.App(appName)
			if err != nil {
				log.Printf("Error retrieving application: %s\n", err)
				return err
			}

			instances = append(instances, app.Instances...)
		default:
			log.Println("Retrieving instances for all registered applications ...")

			apps, err := client.Apps()
			if err != nil {
				log.Printf("Error retrieving applications: %s\n", err)
				return err
			}

			for _, app := range apps {
				instances = append(instances, app.Instances...)
			}
		}

		output := struct {
			XMLName   xml.Name           `xml:"instances"`
			Instances []*eureka.Instance `xml:"instance"`
		}{
			Instances: instances,
		}

		data, err := xml.MarshalIndent(output, "", "  ")
		if err != nil {
			log.Printf("Error rendering output: %s\n", err)
			return err
		}

		log.Println(string(data))
		return nil
	},
}

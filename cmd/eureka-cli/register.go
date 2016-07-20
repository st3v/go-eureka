package main

import (
	"log"

	"github.com/codegangsta/cli"

	"github.com/st3v/go-eureka"
)

var registerCmd = cli.Command{
	Name:  "register",
	Usage: "register an instance with Eureka",

	Flags: []cli.Flag{
		instanceFlag,
		endpointsFlag,
	},

	Action: func(c *cli.Context) error {
		instance := getInstance(c, "register")
		endpoints := getEndpoints(c, "register")

		log.Printf("Registering instance '%s' for application '%s'... \n", instance.Id, instance.AppName)
		client := eureka.NewClient(endpoints)
		if err := client.Register(instance); err != nil {
			log.Printf("Error registering instance with Eureka: %s\n", err)
			return err
		}

		log.Println("Success")
		return nil
	},
}

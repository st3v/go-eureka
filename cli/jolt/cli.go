package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"

	"github.com/micro/cli"
	"github.com/st3v/jolt"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}

func main() {
	app := cli.NewApp()
	app.Name = "jolt"
	app.Usage = "Command-line client for Netflix Eureka"

	var xmlPath string

	app.Commands = []cli.Command{
		{
			Name:  "register",
			Usage: "register an instance with Eureka",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "instance, i",
					Value:       "",
					Usage:       "Path to an XML file that defines the Eureka instance",
					Destination: &xmlPath,
				},
				cli.StringSliceFlag{
					Name:  "endpoint, e",
					Value: &cli.StringSlice{},
					Usage: "Eureka service endpoint, e.g. http://127.0.0.1/eureka/v2",
				},
			},
			Action: func(c *cli.Context) {
				if c.String("instance") == "" {
					cli.ShowCommandHelp(c, "register")
					log.Fatalln("--instance flag is required")
				}

				endpoints := c.StringSlice("endpoint")
				if len(endpoints) == 0 {
					cli.ShowCommandHelp(c, "register")
					log.Fatalln("--endpoint flag is required")
				}

				data, err := ioutil.ReadFile(xmlPath)
				if err != nil {
					log.Fatalf("Error reading instance file: %s", err)
				}

				var instance jolt.Instance
				if err := xml.Unmarshal(data, &instance); err != nil {
					log.Fatalf("Error parsing instance file: %s", err)
				}

				log.Printf("Registering instance for application '%s'... \n", instance.App)
				client := jolt.NewClient(endpoints)
				if err := client.Register(instance); err != nil {
					log.Fatalf("Error registering instance with Eureka: %s", err)
				}

				log.Println("Success")
			},
		},
	}

	app.Run(os.Args)
}

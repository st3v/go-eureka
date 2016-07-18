// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/st3v/go-eureka"
)

var (
	xmlPath   string
	endpoints []string
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register an instance with Eureka",
	RunE: func(cmd *cobra.Command, args []string) error {
		instance, err := getInstance()
		if err != nil {
			return err
		}

		endpoints, err := getEndpoints()
		if err != nil {
			return err
		}

		cmd.Printf("Registering instance '%s' for application '%s'... \n", instance.Id, instance.AppName)
		client := eureka.NewClient(endpoints)
		if err := client.Register(*instance); err != nil {
			return fmt.Errorf("Error registering instance with Eureka: %s", err)
		}

		cmd.Println("Success")
		return nil
	},
}

func getEndpoints() ([]string, error) {
	if len(endpoints) == 0 {
		return nil, errors.New("Missing endpoints")
	}
	return endpoints, nil
}

func getInstance() (*eureka.Instance, error) {
	if xmlPath == "" {
		return nil, errors.New("Missing path to XML file")
	}

	data, err := ioutil.ReadFile(xmlPath)
	if err != nil {
		return nil, fmt.Errorf("Error reading instance file: %s", err)
	}

	instance := new(eureka.Instance)
	if err := xml.Unmarshal(data, instance); err != nil {
		return nil, fmt.Errorf("Error parsing instance file: %s", err)
	}

	return instance, nil
}

func init() {
	RootCmd.AddCommand(registerCmd)

	flags := registerCmd.Flags()
	flags.StringVarP(&xmlPath, "xml", "x", "", "Path to an XML file that defines the Eureka instance")
	flags.StringSliceVarP(&endpoints, "endpoints", "e", []string{}, "Eureka service endpoint, e.g. http://127.0.0.1/eureka/v2")
}

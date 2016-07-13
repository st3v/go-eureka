package jolt_test

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/st3v/jolt"
)

func TestJolt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jolt")
}

func removeIdendation(data []byte) []byte {
	r := regexp.MustCompile("\\n\\s*")
	return r.ReplaceAll(data, []byte{})
}

func instanceFixture() (*jolt.Instance, error) {
	fixture, err := os.Open(filepath.Join("fixtures", "instance.xml"))
	if err != nil {
		return nil, err
	}

	instance := new(jolt.Instance)
	return instance, xml.NewDecoder(fixture).Decode(&instance)
}

func appFixture() (*jolt.App, error) {
	instance, err := instanceFixture()
	if err != nil {
		return nil, err
	}

	return &jolt.App{
		XMLName:   xml.Name{Local: "application"},
		Name:      instance.AppName,
		Instances: []jolt.Instance{*instance},
	}, nil
}

package eureka_test

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/st3v/go-eureka"
)

func TestEureka(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "go-eureka")
}

func removeIdendation(data []byte) []byte {
	r := regexp.MustCompile("\\n\\s*")
	return r.ReplaceAll(data, []byte{})
}

func instanceFixture() (*eureka.Instance, error) {
	fixture, err := os.Open(filepath.Join("fixtures", "instance.xml"))
	if err != nil {
		return nil, err
	}
	defer fixture.Close()

	instance := new(eureka.Instance)
	return instance, xml.NewDecoder(fixture).Decode(&instance)
}

func appFixture() (*eureka.App, error) {
	instance, err := instanceFixture()
	if err != nil {
		return nil, err
	}

	return &eureka.App{
		XMLName:   xml.Name{Local: "application"},
		Name:      instance.AppName,
		Instances: []eureka.Instance{*instance},
	}, nil
}

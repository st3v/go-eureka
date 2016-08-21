package main_test

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/pborman/uuid"

	"github.com/st3v/go-eureka"
)

func TestCLI(t *testing.T) {
	config.DefaultReporterConfig.SlowSpecThreshold = 300.0
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration")
}

var binPath string

var _ = BeforeSuite(func() {
	var err error
	binPath, err = gexec.Build("github.com/st3v/go-eureka/cmd/eureka-cli")
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func execBin(args ...string) *gexec.Session {
	cmd := exec.Command(binPath, args...)
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())
	return session
}

type instancesResult struct {
	XMLName   xml.Name           `xml:"instances"`
	Instances []*eureka.Instance `xml:"instance"`
}

func (r *instancesResult) Contains(that *eureka.Instance) bool {
	for _, this := range r.Instances {
		if this.Equals(that) {
			return true
		}
	}
	return false
}

func instanceToFile(instance *eureka.Instance) string {
	f, err := ioutil.TempFile("", "instance_")
	Expect(err).ToNot(HaveOccurred())
	defer f.Close()

	err = xml.NewEncoder(f).Encode(instance)
	Expect(err).ToNot(HaveOccurred())

	return f.Name()
}

func testInstance() *eureka.Instance {
	return &eureka.Instance{
		ID:         uuid.New(),
		AppName:    uuid.New(),
		HostName:   "host-name",
		IpAddr:     "1.2.3.4",
		VipAddr:    "5.6.7.8",
		Port:       987,
		SecurePort: 789,
		Status:     eureka.StatusUp,
		Metadata:   map[string]string{"key": "value"},
	}
}

func endpoints() []string {
	endpoints := strings.Split(os.Getenv("EUREKA_URLS"), ",")
	for i, e := range endpoints {
		if !strings.HasPrefix(e, "http") {
			endpoints[i] = fmt.Sprintf("http://%s", e)
		}
	}

	return endpoints
}

func endpointFlags() []string {
	endpoints := endpoints()

	flags := make([]string, 0, len(endpoints)*2)
	for _, e := range endpoints {
		flags = append(flags, "-e", e)
	}

	return flags
}

func assertRegistration(instance *eureka.Instance) func() error {
	return func() error {
		endpoint := endpoints()[0]
		url := fmt.Sprintf("%s/apps", endpoint)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		req.Header.Set("Accepts", "application/xml")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Unexpected Status Code: %d", resp.StatusCode)
		}

		result := new(eureka.Registry)

		err = xml.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return err
		}

		for _, app := range result.Apps {
			if strings.ToLower(app.Name) != instance.AppName {
				continue
			}

			for _, i := range app.Instances {
				if i.Equals(instance) {
					return nil
				}
			}
		}

		return errors.New("Instance not registered")
	}
}

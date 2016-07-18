// +build integration

package main_test

import (
	"encoding/xml"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/st3v/go-eureka"
)

var _ = Describe("CLI", func() {
	var (
		instances         []*eureka.Instance
		instanceFilePaths []string
		session           *gexec.Session

		timeout = 2 * time.Minute
	)

	BeforeEach(func() {
		instances = []*eureka.Instance{testInstance(), testInstance()}

		instanceFilePaths = make([]string, 0, len(instances))
		for _, i := range instances {
			instanceFilePaths = append(instanceFilePaths, instanceToFile(i))
		}
	})

	AfterEach(func() {
		for _, p := range instanceFilePaths {
			os.Remove(p)
		}
	})

	It("handles the instance lifecycle", func() {
		// register
		for _, path := range instanceFilePaths {
			session = execBin(append([]string{"register", "-i", path}, endpointFlags()...)...)
			Eventually(session).Should(gexec.Exit(0))
		}

		// verify registration
		for _, i := range instances {
			Eventually(assertRegistration(i), timeout).ShouldNot(HaveOccurred())
		}

		// get all instances
		session = execBin(append([]string{"instances"}, endpointFlags()...)...)
		Eventually(session).Should(gexec.Exit(0))

		result := new(instancesResult)
		err := xml.Unmarshal(session.Out.Contents(), result)
		Expect(err).ToNot(HaveOccurred())

		Expect(result.Contains(instances[0])).To(BeTrue())
		Expect(result.Contains(instances[1])).To(BeTrue())

		// get instance by app name
		session = execBin(append([]string{"instances", "-a", instances[0].AppName}, endpointFlags()...)...)
		Eventually(session).Should(gexec.Exit(0))

		result = new(instancesResult)
		err = xml.Unmarshal(session.Out.Contents(), result)
		Expect(err).ToNot(HaveOccurred())

		Expect(result.Instances).To(HaveLen(1))
		Expect(result.Contains(instances[0])).To(BeTrue())

		// get instance by id
		session = execBin(append([]string{"instances", "-i", instances[1].Id}, endpointFlags()...)...)
		Eventually(session).Should(gexec.Exit(0))

		result = new(instancesResult)
		err = xml.Unmarshal(session.Out.Contents(), result)
		Expect(err).ToNot(HaveOccurred())

		Expect(result.Instances).To(HaveLen(1))
		Expect(result.Contains(instances[1])).To(BeTrue())

		// get instance by app name and id
		session = execBin(append([]string{"instances", "-a", instances[0].AppName, "-i", instances[0].Id}, endpointFlags()...)...)
		Eventually(session).Should(gexec.Exit(0))

		result = new(instancesResult)
		err = xml.Unmarshal(session.Out.Contents(), result)
		Expect(err).ToNot(HaveOccurred())

		Expect(result.Instances).To(HaveLen(1))
		Expect(result.Contains(instances[0])).To(BeTrue())

		// heartbeat
		session = execBin(append([]string{"heartbeat", "-i", instanceFilePaths[0]}, endpointFlags()...)...)
		Eventually(session).Should(gexec.Exit(0))

		// override status
		session = execBin(append([]string{"override-status", "OUT_OF_SERVICE", "-i", instanceFilePaths[0]}, endpointFlags()...)...)
		Eventually(session).Should(gexec.Exit(0))

		// verify override
		session = execBin(append([]string{"instances", "-i", instances[0].Id}, endpointFlags()...)...)
		Eventually(session).Should(gexec.Exit(0))

		result = new(instancesResult)
		err = xml.Unmarshal(session.Out.Contents(), result)
		Expect(err).ToNot(HaveOccurred())

		Expect(result.Instances).To(HaveLen(1))
		Expect(result.Instances[0].Status).To(Equal(eureka.StatusOutOfService))
		Expect(result.Instances[0].StatusOverride).To(Equal(eureka.StatusOutOfService))

		// remove status override
		session = execBin(append([]string{"remove-override", "UP", "-i", instanceFilePaths[0]}, endpointFlags()...)...)
		Eventually(session).Should(gexec.Exit(0))

		// verify override removal
		session = execBin(append([]string{"instances", "-i", instances[0].Id}, endpointFlags()...)...)
		Eventually(session).Should(gexec.Exit(0))

		result = new(instancesResult)
		err = xml.Unmarshal(session.Out.Contents(), result)
		Expect(err).ToNot(HaveOccurred())

		Expect(result.Instances).To(HaveLen(1))
		Expect(result.Instances[0].Status).To(Equal(eureka.StatusUp))
		Expect(result.Instances[0].StatusOverride).To(Equal(eureka.StatusUnknown))

		// deregister
		for _, path := range instanceFilePaths {
			session = execBin(append([]string{"deregister", "-i", path}, endpointFlags()...)...)
			Eventually(session).Should(gexec.Exit(0))
		}

		// verify deregistration
		for _, i := range instances {
			Eventually(assertRegistration(i), timeout).Should(HaveOccurred())
		}
	})
})

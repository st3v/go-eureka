package main_test

import (
	"os/exec"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var binPath string

func TestCLI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "eureka-cli")
}

func execBin(args ...string) *gexec.Session {
	cmd := exec.Command(binPath, args...)
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())
	return session
}

func removeIdendation(data []byte) []byte {
	r := regexp.MustCompile("\\n\\s*")
	return r.ReplaceAll(data, []byte{})
}

var _ = BeforeSuite(func() {
	var err error
	binPath, err = gexec.Build("github.com/st3v/go-eureka/cmd/eureka-cli")
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

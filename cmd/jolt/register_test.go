package main_test

import (
	"io/ioutil"
	"net/http"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("register command", func() {
	var (
		server      *ghttp.Server
		instanceXml []byte
		xmlPath     string
		statusCode  int
		args        []string
		session     *gexec.Session
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		xmlPath = filepath.Join("..", "..", "fixtures", "instance.xml")
		args = []string{"register", "-i", xmlPath, "-e", server.URL()}

		var err error
		instanceXml, err = ioutil.ReadFile(xmlPath)
		Expect(err).ToNot(HaveOccurred())

		statusCode = http.StatusNoContent
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/apps/myapp"),
				ghttp.VerifyContentType("application/xml"),
				ghttp.VerifyBody(removeIdendation(instanceXml)),
				ghttp.RespondWithPtr(&statusCode, nil),
			),
		)
	})

	AfterEach(func() {
		server.Close()
	})

	JustBeforeEach(func() {
		session = execBin(args...)
	})

	It("exits with exit code 0", func() {
		Eventually(session).Should(gexec.Exit(0))
	})

	It("sends the correct POST request", func() {
		Eventually(server.ReceivedRequests).Should(HaveLen(1))
	})

	It("provides basic logs on stdout", func() {
		Eventually(session).Should(gbytes.Say("Registering instance 'id' for application 'myapp'..."))
		Eventually(session).Should(gbytes.Say("Success"))
	})

	Context("when the --instance flag has not been specified", func() {
		BeforeEach(func() {
			args = []string{"register", "-e", server.URL()}
		})

		It("exits with exit code 1", func() {
			Eventually(session).Should(gexec.Exit(1))
		})

		It("provides a corresponding error message on stdout", func() {
			Eventually(session).Should(gbytes.Say("--instance flag is required"))
		})
	})

	Context("when no --endpoint flag has not been specified", func() {
		BeforeEach(func() {
			args = []string{"register", "-i", xmlPath}
		})

		It("exits with exit code 1", func() {
			Eventually(session).Should(gexec.Exit(1))
		})

		It("provides a corresponding error message on stdout", func() {
			Eventually(session).Should(gbytes.Say("--endpoint flag is required"))
		})
	})

	Context("when the instance file does not exist", func() {
		BeforeEach(func() {
			args = []string{"register", "-i", "/path/to/nowhere", "-e", server.URL()}
		})

		It("exits with exit code 1", func() {
			Eventually(session).Should(gexec.Exit(1))
		})

		It("does not send an HTTP request", func() {
			Eventually(session).Should(gexec.Exit())
			Expect(server.ReceivedRequests()).To(BeEmpty())
		})

		It("provides a corresponding error message on stdout", func() {
			Eventually(session).Should(gbytes.Say("Error reading instance file"))
		})

		It("does not report success on stdout", func() {
			Eventually(session).Should(gexec.Exit())
			Expect(session).ToNot(gbytes.Say("Success"))
		})
	})

	Context("when the instance file does not specifiy an instance", func() {
		BeforeEach(func() {
			path := filepath.Join("..", "..", "fixtures", "foo.xml")
			args = []string{"register", "-i", path, "-e", server.URL()}
		})

		It("exits with exit code 1", func() {
			Eventually(session).Should(gexec.Exit(1))
		})

		It("does not send an HTTP request", func() {
			Eventually(session).Should(gexec.Exit())
			Expect(server.ReceivedRequests()).To(BeEmpty())
		})

		It("provides a corresponding error message on stdout", func() {
			Eventually(session).Should(gbytes.Say("Error parsing instance file"))
		})

		It("does not report success on stdout", func() {
			Eventually(session).Should(gexec.Exit())
			Expect(session).ToNot(gbytes.Say("Success"))
		})
	})
})

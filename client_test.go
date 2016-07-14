package eureka_test

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/st3v/go-eureka"
)

var _ = Describe("eureka", func() {
	var (
		server     *ghttp.Server
		client     *eureka.Client
		instance   *eureka.Instance
		statusCode int
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = eureka.NewClient([]string{server.URL()})

		var err error
		instance, err = instanceFixture()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe(".Register", func() {
		BeforeEach(func() {
			instanceXml, err := ioutil.ReadFile(filepath.Join("fixtures", "instance.xml"))
			Expect(err).ToNot(HaveOccurred())

			route := fmt.Sprintf("/apps/%s", instance.AppName)
			statusCode = http.StatusNoContent
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", route),
					ghttp.VerifyContentType("application/xml"),
					ghttp.VerifyBody(removeIdendation(instanceXml)),
					ghttp.RespondWithPtr(&statusCode, nil),
				),
			)
		})

		It("sends the correct request", func() {
			client.Register(*instance)
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns no error", func() {
			err := client.Register(*instance)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				statusCode = http.StatusInternalServerError
			})

			It("returns an error", func() {
				err := client.Register(*instance)
				Expect(err).To(MatchError("Unexpected response code 500"))
			})
		})
	})

	Describe(".Deregister", func() {
		BeforeEach(func() {
			route := fmt.Sprintf("/apps/%s/%s", instance.AppName, instance.Id)
			statusCode = http.StatusOK
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", route),
					ghttp.RespondWithPtr(&statusCode, nil),
				),
			)
		})

		It("sends the correct request", func() {
			client.Deregister(*instance)
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns no error", func() {
			err := client.Deregister(*instance)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				statusCode = http.StatusInternalServerError
			})

			It("returns an error", func() {
				err := client.Deregister(*instance)
				Expect(err).To(MatchError("Unexpected response code 500"))
			})
		})
	})

	Describe(".Heartbeat", func() {
		BeforeEach(func() {
			route := fmt.Sprintf("/apps/%s/%s", instance.AppName, instance.Id)
			statusCode = http.StatusOK
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", route),
					ghttp.RespondWithPtr(&statusCode, nil),
				),
			)
		})

		It("sends the correct request", func() {
			client.Heartbeat(*instance)
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns no error", func() {
			err := client.Heartbeat(*instance)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				statusCode = http.StatusInternalServerError
			})

			It("returns an error", func() {
				err := client.Heartbeat(*instance)
				Expect(err).To(MatchError("Unexpected response code 500"))
			})
		})
	})

	Describe(".Apps", func() {
		var app *eureka.App

		BeforeEach(func() {
			var err error
			app, err = appFixture()
			Expect(err).ToNot(HaveOccurred())

			response := eureka.Registry{
				Apps: []eureka.App{*app, *app},
			}

			var body []byte
			body, err = xml.Marshal(response)
			Expect(err).ToNot(HaveOccurred())

			statusCode = http.StatusOK
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/apps"),
					ghttp.RespondWithPtr(&statusCode, &body),
				),
			)
		})

		It("sends the correct request", func() {
			client.Apps()
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns no error", func() {
			_, err := client.Apps()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns the correct apps", func() {
			apps, _ := client.Apps()
			Expect(apps).To(HaveLen(2))
			Expect(apps[0]).To(Equal(*app))
			Expect(apps[1]).To(Equal(*app))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				statusCode = http.StatusInternalServerError
			})

			It("returns an error", func() {
				_, err := client.Apps()
				Expect(err).To(MatchError("Unexpected response code 500"))
			})
		})
	})

	Describe(".App", func() {
		var app *eureka.App

		BeforeEach(func() {
			var err error
			app, err = appFixture()
			Expect(err).ToNot(HaveOccurred())

			var body []byte
			body, err = xml.Marshal(app)
			Expect(err).ToNot(HaveOccurred())

			route := fmt.Sprintf("/apps/%s", app.Name)
			statusCode = http.StatusOK
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", route),
					ghttp.RespondWithPtr(&statusCode, &body),
				),
			)
		})

		It("sends the correct request", func() {
			client.App(app.Name)
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns no error", func() {
			_, err := client.App(app.Name)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns the correct app", func() {
			actual, _ := client.App(app.Name)
			Expect(actual).To(Equal(*app))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				statusCode = http.StatusInternalServerError
			})

			It("returns an error", func() {
				_, err := client.App(app.Name)
				Expect(err).To(MatchError("Unexpected response code 500"))
			})
		})
	})

	Describe(".AppInstance", func() {
		BeforeEach(func() {
			var err error
			instance, err = instanceFixture()
			Expect(err).ToNot(HaveOccurred())

			var body []byte
			body, err = xml.Marshal(instance)
			Expect(err).ToNot(HaveOccurred())

			route := fmt.Sprintf("/apps/%s/%s", instance.AppName, instance.Id)
			statusCode = http.StatusOK
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", route),
					ghttp.RespondWithPtr(&statusCode, &body),
				),
			)
		})

		It("sends the correct request", func() {
			client.AppInstance(instance.AppName, instance.Id)
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns no error", func() {
			_, err := client.AppInstance(instance.AppName, instance.Id)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns the correct instance", func() {
			actual, _ := client.AppInstance(instance.AppName, instance.Id)
			Expect(actual).To(Equal(*instance))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				statusCode = http.StatusInternalServerError
			})

			It("returns an error", func() {
				_, err := client.AppInstance(instance.AppName, instance.Id)
				Expect(err).To(MatchError("Unexpected response code 500"))
			})
		})
	})

	Describe(".Instance", func() {
		BeforeEach(func() {
			var err error
			instance, err = instanceFixture()
			Expect(err).ToNot(HaveOccurred())

			var body []byte
			body, err = xml.Marshal(instance)
			Expect(err).ToNot(HaveOccurred())

			route := fmt.Sprintf("/instances/%s", instance.Id)
			statusCode = http.StatusOK
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", route),
					ghttp.RespondWithPtr(&statusCode, &body),
				),
			)
		})

		It("sends the correct request", func() {
			client.Instance(instance.Id)
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns no error", func() {
			_, err := client.Instance(instance.Id)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns the correct instance", func() {
			actual, _ := client.Instance(instance.Id)
			Expect(actual).To(Equal(*instance))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				statusCode = http.StatusInternalServerError
			})

			It("returns an error", func() {
				_, err := client.Instance(instance.Id)
				Expect(err).To(MatchError("Unexpected response code 500"))
			})
		})
	})
})

package jolt_test

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/st3v/jolt"
)

var _ = Describe("jolt", func() {
	var (
		server      *ghttp.Server
		client      *jolt.Client
		instanceXml []byte
		instance    jolt.Instance
		statusCode  int
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = jolt.NewClient([]string{server.URL()})

		var err error
		instanceXml, err = ioutil.ReadFile(filepath.Join("fixtures", "instance.xml"))
		Expect(err).ToNot(HaveOccurred())

		err = xml.Unmarshal(instanceXml, &instance)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe(".Register", func() {
		BeforeEach(func() {
			route := fmt.Sprintf("/apps/%s", instance.App)
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

		It("returns no error", func() {
			err := client.Register(instance)
			Expect(err).ToNot(HaveOccurred())
		})

		It("sends the correct POST request to the /apps route", func() {
			client.Register(instance)
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				statusCode = http.StatusInternalServerError
			})

			It("returns an error", func() {
				err := client.Register(instance)
				Expect(err).To(MatchError("Unexpected response code 500"))
			})
		})
	})

	Describe(".Deregister", func() {
		BeforeEach(func() {
			route := fmt.Sprintf("/apps/%s/%s", instance.App, instance.Id)
			statusCode = http.StatusOK
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", route),
					ghttp.RespondWithPtr(&statusCode, nil),
				),
			)
		})

		It("returns no error", func() {
			err := client.Deregister(instance)
			Expect(err).ToNot(HaveOccurred())
		})

		It("sends the correct POST request to the /apps route", func() {
			client.Deregister(instance)
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				statusCode = http.StatusInternalServerError
			})

			It("returns an error", func() {
				err := client.Deregister(instance)
				Expect(err).To(MatchError("Unexpected response code 500"))
			})
		})
	})
})

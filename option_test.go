package eureka

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("client options", func() {
	Describe("No option", func() {
		It("uses the default http client", func() {
			client := NewClient([]string{"endpoint"})
			Expect(client.httpClient).To(Equal(defaultHttpClient))
		})
	})

	Describe("HttpClient", func() {
		It("switches the specified http client", func() {
			hc := &http.Client{}
			client := NewClient([]string{"endpoint"}, HttpClient(hc))
			Expect(client.httpClient).To(Equal(hc))
		})
	})
})

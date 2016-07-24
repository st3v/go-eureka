package eureka

import (
	"net/http"
	"reflect"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/st3v/go-eureka/retry"
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

	Describe("Retry", func() {
		var (
			selector retry.Selector = func(_ []string) retry.Endpoint { return func(_ uint) string { return "" } }
			allow    retry.Allow    = func(_ uint) bool { return true }
			delay    retry.Delay    = func(_ uint) time.Duration { return 0 }

			client = NewClient([]string{"endpoint"}, Retry(selector, allow, delay))
		)

		It("sets retry selector", func() {
			Expect(reflect.ValueOf(client.retrySelector)).To(Equal(reflect.ValueOf(selector)))
			Expect(reflect.ValueOf(client.retryAllow)).To(Equal(reflect.ValueOf(allow)))
			Expect(reflect.ValueOf(client.retryDelay)).To(Equal(reflect.ValueOf(delay)))
		})
	})
})

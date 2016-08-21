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
			Expect(client.httpClient).To(Equal(DefaultHTTPClient))
		})

		It("uses the default retry selector", func() {
			client := NewClient([]string{"endpoint"})
			Expect(reflect.ValueOf(client.retrySelector)).To(Equal(reflect.ValueOf(DefaultRetrySelector)))
		})

		It("uses the default retry limit", func() {
			client := NewClient([]string{"endpoint"})
			Expect(reflect.ValueOf(client.retryLimit)).To(Equal(reflect.ValueOf(DefaultRetryLimit)))
		})

		It("uses the default retry delay", func() {
			client := NewClient([]string{"endpoint"})
			Expect(reflect.ValueOf(client.retryDelay)).To(Equal(reflect.ValueOf(DefaultRetryDelay)))
		})
	})

	Describe("HTTPClient", func() {
		It("switches to the specified http client", func() {
			hc := &http.Client{}
			client := NewClient([]string{"endpoint"}, HTTPClient(hc))
			Expect(client.httpClient).To(Equal(hc))
		})
	})

	Describe("RetrySelector", func() {
		var selector retry.Selector = func(_ []string) retry.Endpoint {
			return func(_ uint) string { return "" }
		}

		It("sets retry selector", func() {
			client := NewClient([]string{"endpoint"}, RetrySelector(selector))
			Expect(reflect.ValueOf(client.retrySelector)).To(Equal(reflect.ValueOf(selector)))
		})
	})

	Describe("RetryLimit", func() {
		var allow retry.Allow = func(_ uint) bool { return true }

		It("sets retry limit", func() {
			client := NewClient([]string{"endpoint"}, RetryLimit(allow))
			Expect(reflect.ValueOf(client.retryLimit)).To(Equal(reflect.ValueOf(allow)))
		})
	})

	Describe("RetryDelay", func() {
		var delay retry.Delay = func(_ uint) time.Duration { return 0 }

		It("sets retry delay", func() {
			client := NewClient([]string{"endpoint"}, RetryDelay(delay))
			Expect(reflect.ValueOf(client.retryDelay)).To(Equal(reflect.ValueOf(delay)))
		})
	})
})

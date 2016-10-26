package eureka

import (
	"crypto/tls"
	"net/http"
	"reflect"
	"time"

	"golang.org/x/oauth2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/st3v/go-eureka/retry"
)

var _ = Describe("client options", func() {
	Describe("No option", func() {
		It("uses the default http client", func() {
			actual := NewClient([]string{"endpoint"}).httpClient
			expected := &http.Client{
				Timeout:   DefaultTimeout,
				Transport: DefaultTransport,
			}

			Expect(actual).To(Equal(expected))
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

	Describe("HTTPTimeout", func() {
		It("sets the timeout for the internal HTTP client", func() {
			timeout := 123 * time.Second
			client := NewClient([]string{"endpoint"}, HTTPTimeout(timeout))
			Expect(client.httpClient.Timeout).To(Equal(timeout))
		})
	})

	Describe("HTTPTransport", func() {
		It("sets the transport for the internal HTTP client", func() {
			transport := new(http.Transport)
			client := NewClient([]string{"endpoint"}, HTTPTransport(transport))
			Expect(client.httpClient.Transport).To(BeIdenticalTo(transport))
		})
	})

	Describe("TLSConfig", func() {
		It("sets the tls config for the internal HTTP client", func() {
			tlsConfig := &tls.Config{InsecureSkipVerify: true}
			client := NewClient([]string{"endpoint"}, TLSConfig(tlsConfig))
			transport, ok := client.httpClient.Transport.(*http.Transport)
			Expect(ok).To(BeTrue())
			Expect(transport.TLSClientConfig).To(BeIdenticalTo(tlsConfig))
		})
	})

	Describe("Oauth2ClientCredentials", func() {
		It("wraps the internal http client transport in an oauth2 transport", func() {
			id, secret, uri, scope := "client-id", "client-secret", "token-uri", "scope"
			client := NewClient([]string{"endpoint"}, Oauth2ClientCredentials(id, secret, uri, scope))
			transport, ok := client.httpClient.Transport.(*oauth2.Transport)
			Expect(ok).To(BeTrue())
			Expect(transport.Base).To(BeIdenticalTo(DefaultTransport))
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

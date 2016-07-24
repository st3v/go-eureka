package retry_test

import (
	"errors"
	"math"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/st3v/go-eureka/retry"
)

func TestEureka(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "retry")
}

var _ = Describe("retry", func() {
	Describe(".Strategy", func() {
		Describe(".Apply", func() {
			var (
				retries int
				numErr  = 10
				someErr = errors.New("some error")

				action = func(endpoint string) error {
					defer func() { retries += 1 }()
					if retries < numErr {
						return someErr
					}
					return nil
				}
			)

			BeforeEach(func() {
				retries = 0
			})

			It("retries until it succeeds", func() {
				strategy := retry.NewStrategy(
					retry.RoundRobin([]string{"one"}),
					retry.NoLimit(),
					retry.NoDelay(),
				)

				err := strategy.Apply(action)

				Expect(err).ToNot(HaveOccurred())
				Expect(retries).To(Equal(numErr + 1))
			})

			It("retries only until the allowed limit", func() {
				var (
					limit    = numErr - 1
					strategy = retry.NewStrategy(
						retry.RoundRobin([]string{"one"}),
						retry.Limit(limit),
						retry.NoDelay(),
					)

					err = strategy.Apply(action)
				)

				Expect(err).To(MatchError(someErr))
				Expect(retries).To(Equal(limit))
			})

			It("follows the right strategy", func() {
				var (
					delayCalled    bool
					allowCalled    bool
					endpointCalled bool
					actionCalled   bool

					strategy = retry.NewStrategy(
						func(_ uint) string {
							endpointCalled = true
							return ""
						},
						func(_ uint) bool {
							allowCalled = true
							return true
						},
						func(_ uint) time.Duration {
							delayCalled = true
							return 0
						},
					)

					action = func(_ string) error {
						actionCalled = true
						return nil
					}
				)

				err := strategy.Apply(action)

				Expect(err).ToNot(HaveOccurred())
				Expect(endpointCalled).To(BeTrue())
				Expect(allowCalled).To(BeTrue())
				Expect(delayCalled).To(BeTrue())
				Expect(actionCalled).To(BeTrue())
			})
		})
	})

	Describe(".Endpoint", func() {
		var endpoints = []string{"one", "two", "three"}

		Describe(".RoundRobin", func() {
			It("returns the expected series of endpoint", func() {
				endpoint := retry.RoundRobin(endpoints)

				for i := 0; i < 100; i++ {
					Expect(endpoint(uint(i))).To(Equal(endpoints[i%len(endpoints)]))
				}
			})
		})

		Describe(".Random", func() {
			It("returns a random series of endpoint", func() {
				var (
					l = 1000
					a = make([]string, 0, l)
					b = make([]string, 0, l)

					endpoint = retry.Random(endpoints)
				)

				for i := 0; i < l; i++ {
					a = append(a, endpoint(uint(i)))
					b = append(b, endpoint(uint(i)))
				}

				Expect(strings.Join(a, "")).ToNot(Equal(strings.Join(b, "")))
			})
		})
	})

	Describe(".Allow", func() {
		Describe(".NoLimit", func() {
			It("always returns true", func() {
				allow := retry.NoLimit()
				for i := uint(0); i < 1000; i++ {
					Expect(allow(i)).To(BeTrue())
				}
			})
		})

		Describe(".Limit", func() {
			It("returns the expected bool", func() {
				var (
					limit = 100
					allow = retry.Limit(limit)
				)

				for i := 0; i < limit; i++ {
					Expect(allow(uint(i))).To(BeTrue())
				}

				for i := limit; i < limit*2; i++ {
					Expect(allow(uint(i))).To(BeFalse())
				}
			})
		})
	})

	Describe(".Delay", func() {
		Describe(".NoDelay", func() {
			It("always returns 0", func() {
				delay := retry.NoDelay()
				for i := uint(0); i < 100; i++ {
					Expect(delay(i)).To(Equal(time.Duration(0)))
				}
			})
		})

		Describe(".Constant", func() {
			It("always returns the specified delay", func() {
				var (
					term  = 123 * time.Second
					delay = retry.Constant(term)
				)

				for i := uint(0); i < 10; i++ {
					Expect(delay(i)).To(Equal(term))
				}
			})
		})

		Describe(".Linear", func() {
			It("returns the expected series of delays", func() {
				var (
					term  = 123 * time.Second
					delay = retry.Linear(term)
				)

				for i := uint(0); i < 10; i++ {
					Expect(delay(i)).To(Equal(term * time.Duration(i)))
				}
			})
		})

		Describe(".Exponential", func() {
			It("returns the expected series of delays", func() {
				var (
					term  = 123 * time.Second
					delay = retry.Exponential(term)
				)

				for i := uint(0); i < 10; i++ {
					want := time.Duration(math.Pow(2.0, float64(i))) * term
					Expect(delay(i)).To(Equal(want))
				}
			})
		})
	})
})

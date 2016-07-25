package retry

import (
	"math"
	"math/rand"
	"time"
)

type Strategy func(action Action) error

func (s Strategy) Apply(action Action) error {
	return s(action)
}

type Action func(endpoint string) error

type Endpoint func(attempt uint) string

type Selector func(endpoints []string) Endpoint

type Allow func(attempt uint) bool

type Delay func(attempt uint) time.Duration

func NewStrategy(endpoint Endpoint, allow Allow, delay Delay) Strategy {
	return func(action Action) error {
		var err error

		for i := uint(0); allow(i) && (i == 0 || err != nil); i++ {
			time.Sleep(delay(i))
			err = action(endpoint(i))
		}

		return err
	}
}

func RoundRobin(endpoints []string) Endpoint {
	return func(attempt uint) string {
		return endpoints[attempt%uint(len(endpoints))]
	}
}

func Random(endpoints []string) Endpoint {
	rand.Seed(time.Now().Unix())
	return func(_ uint) string {
		return endpoints[rand.Intn(len(endpoints))]
	}
}

func NoRetries() Allow {
	return func(attempt uint) bool {
		return attempt == 0
	}
}

func MaxRetries(max int) Allow {
	return func(attempt uint) bool {
		return attempt < uint(max)
	}
}

func NoDelay() Delay {
	return ConstantDelay(0)
}

func ConstantDelay(delay time.Duration) Delay {
	return func(attempt uint) time.Duration {
		if attempt == 0 {
			return 0
		}
		return delay
	}
}

func LinearBackoff(delay time.Duration) Delay {
	return func(attempt uint) time.Duration {
		return time.Duration(attempt) * delay
	}
}

func ExponentialBackoff(delay time.Duration) Delay {
	return func(attempt uint) time.Duration {
		return time.Duration(math.Pow(2.0, float64(attempt))) * delay
	}
}

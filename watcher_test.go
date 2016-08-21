package eureka

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Watcher", func() {
	var (
		interval = 10 * time.Millisecond

		existingApp *App
		registry    *mockRegistry
		watcher     *Watcher
	)

	BeforeEach(func() {
		existingApp = &App{
			Name: "existing",
			Instances: []*Instance{
				&Instance{
					ID:       "one",
					HostName: "one.example.com",
				},
				&Instance{
					ID:       "two",
					HostName: "two.example.com",
				},
			},
		}

		registry = newMockRegistry()
		registry.Register(existingApp)

		watcher = newWatcher(registry, interval)

		// should receive event for the above register
		Eventually(watcher.Events()).Should(Receive())
	})

	AfterEach(func() {
		watcher.Stop()
	})

	It("reports instances for newly registered apps", func() {
		instance := &Instance{
			ID:       "one",
			HostName: "new.example.com",
		}

		app := &App{
			Name:      "new",
			Instances: []*Instance{instance},
		}

		expectedEvent := Event{
			EventInstanceRegistered,
			instance,
		}

		registry.Register(app)

		Eventually(watcher.Events()).Should(Receive(Equal(expectedEvent)))
	})

	It("reports new instances for previously registered apps", func() {
		instance := &Instance{
			ID:       "two",
			HostName: "two.example.com",
		}
		existingApp.Instances = append(existingApp.Instances, instance)

		expectedEvent := Event{
			EventInstanceRegistered,
			instance,
		}

		registry.Register(existingApp)

		Eventually(watcher.Events()).Should(Receive(Equal(expectedEvent)))
	})

	It("reports instances that have been changed", func() {
		existingApp.Instances[0] = &Instance{
			ID:       "one",
			HostName: "updated.example.com",
		}

		expectedEvent := Event{
			EventInstanceUpdated,
			existingApp.Instances[0],
		}

		registry.Register(existingApp)

		Eventually(watcher.Events()).Should(Receive(Equal(expectedEvent)))
	})

	It("reports instances that have been deregistered", func() {
		expectedEvent := Event{
			EventInstanceDeregistered,
			existingApp.Instances[0],
		}

		existingApp.Instances = existingApp.Instances[1:]

		registry.Register(existingApp)

		Eventually(watcher.Events()).Should(Receive(Equal(expectedEvent)))
	})

	It("reports all instances when the app has been deregistered as a whole", func() {
		expectedEvent := Event{
			EventInstanceDeregistered,
			existingApp.Instances[0],
		}

		registry.Deregister(existingApp.Name)

		Eventually(watcher.Events()).Should(Receive(Equal(expectedEvent)))
	})
})

type mockRegistry struct {
	mtx  sync.RWMutex
	apps map[string]*App
}

func newMockRegistry() *mockRegistry {
	return &mockRegistry{
		apps: map[string]*App{},
	}
}

func (m *mockRegistry) Register(app *App) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.apps[app.Name] = app
}

func (m *mockRegistry) Deregister(name string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	delete(m.apps, name)
}

func (m *mockRegistry) Apps() ([]*App, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	apps := make([]*App, 0, len(m.apps))
	for _, a := range m.apps {
		apps = append(apps, a)
	}

	return apps, nil
}

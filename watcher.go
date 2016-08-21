package eureka

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
)

// DefaultPollInterval defines the default interval at which the watcher queries
// the registry.
const DefaultPollInterval = 30 * time.Second

// EventType defines the type of an observed event.
type EventType uint8

const (
	// EventInstanceRegistered indicates that a newly registered instance has
	// been observed.
	EventInstanceRegistered EventType = iota

	// EventInstanceDeregistered indicates that a previously registered instance
	// is no longer registered.
	EventInstanceDeregistered

	// EventInstanceUpdated indicates that a previously registered instance has
	// changed in the registry, e.g. status or metadata changes have been observed.
	EventInstanceUpdated
)

// Event holds information about the type and subject of an observation.
type Event struct {
	Type     EventType
	Instance *Instance
}

// Watcher can be used to observe the registry for changes with respect
// to the instances of particular app.
type Watcher struct {
	events    chan Event
	instances map[string]*Instance
	cancel    context.CancelFunc
}

// Registry is being used to poll for registered Apps.
type Registry interface {
	Apps() ([]*App, error)
}

func newWatcher(registry Registry, pollInterval time.Duration) *Watcher {
	ctx, cancel := context.WithCancel(context.Background())

	watcher := &Watcher{
		events: make(chan Event),
		cancel: cancel,
	}

	go watcher.poll(ctx, registry, pollInterval)

	return watcher
}

// Stop the watcher, i.e. the registry is no longer being polled.
func (w *Watcher) Stop() {
	w.cancel()
}

// Events returns a channel that can be used to listen for changes to the app
// observed by this watcher.
func (w *Watcher) Events() <-chan Event {
	return w.events
}

func (w *Watcher) poll(ctx context.Context, registry Registry, interval time.Duration) {
	tick := time.NewTicker(interval)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			if apps, err := registry.Apps(); err == nil {
				w.update(apps)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (w *Watcher) update(apps []*App) {
	current := map[string]*Instance{}

	// check if instances are new or have changed
	for _, a := range apps {
		for _, i := range a.Instances {
			key := key(a, i)
			current[key] = i

			prev, found := w.instances[key]
			if !found {
				w.notify(EventInstanceRegistered, i)
				continue
			}

			delete(w.instances, key)

			if !i.Equals(prev) {
				w.notify(EventInstanceUpdated, i)
			}
		}
	}

	// instances we haven't deleted above are not registered anymore
	for _, i := range w.instances {
		w.notify(EventInstanceDeregistered, i)
	}

	// reset instances
	w.instances = current
}

func (w *Watcher) notify(t EventType, i *Instance) {
	// blocking
	w.events <- Event{t, i}
}

func key(a *App, i *Instance) string {
	// instance ids might not be unique across apps
	return fmt.Sprintf("%s-%s", a.Name, i.ID)
}

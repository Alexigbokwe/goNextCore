package events

import (
	"context"
	"github.com/Alexigbokwe/gonext-framework/core/logger"
	"github.com/Alexigbokwe/gonext-framework/core/utils"
	"sync"

	"go.uber.org/zap"
)

// Event is the interface that all events must implement
type Event interface {
	Name() string
}

// Listener is a function that handles an event
type Listener func(ctx context.Context, event Event) error

// Dispatcher manages event listeners and dispatching
type Dispatcher struct {
	listeners map[string][]Listener
	mu        sync.RWMutex
}

var (
	dispatcherInstance *Dispatcher
	once               sync.Once
)

// GetDispatcher returns the singleton instance
func GetDispatcher() *Dispatcher {
	once.Do(func() {
		dispatcherInstance = &Dispatcher{
			listeners: make(map[string][]Listener),
		}
	})
	return dispatcherInstance
}

// Register adds a listener for a specific event
func (d *Dispatcher) Register(eventName string, listener Listener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.listeners[eventName] = append(d.listeners[eventName], listener)
}

// Dispatch fires an event to all registered listeners (Synchronous)
func (d *Dispatcher) Dispatch(ctx context.Context, event Event) error {
	d.mu.RLock()
	listeners, ok := d.listeners[event.Name()]
	d.mu.RUnlock()

	if !ok {
		return nil // No listeners
	}

	for _, listener := range listeners {
		if err := listener(ctx, event); err != nil {
			logger.Log.Error("Error in event listener",
				zap.String("event", event.Name()),
				zap.Error(err),
			)
			return err
		}
	}
	return nil
}

// DispatchAsync fires an event in the background
func (d *Dispatcher) DispatchAsync(event Event) {
	utils.RunBackground(func() error {
		return d.Dispatch(context.Background(), event)
	})
}

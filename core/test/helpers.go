package test

import (
	"time"

	"github.com/Alexigbokwe/gonext-framework/core"

	"github.com/hibiken/asynq"
)

// MockTaskQueue for testing without Redis
type MockTaskQueue struct {
	Tasks []struct {
		TypeName string
		Payload  interface{}
	}
}

func (m *MockTaskQueue) Enqueue(typeName string, payload interface{}, opts ...asynq.Option) error {
	m.Tasks = append(m.Tasks, struct {
		TypeName string
		Payload  interface{}
	}{TypeName: typeName, Payload: payload})
	return nil
}

func (m *MockTaskQueue) Start()    {}
func (m *MockTaskQueue) Shutdown() {}

// CreateTestApp creates an app instance with mocked dependencies for testing
func CreateTestApp() *core.App {
	// Initialize minimal config
	// Initialize minimal config
	// config := &config.Config{
	// 	Server: config.ServerConfig{Port: "0"},
	// }

	a := core.NewApp()
	// Register mocks here if needed
	return a
}

func WaitForServer(a *core.App, timeout time.Duration) bool {
	// Simple helper to wait for server (if needed in E2E tests)
	return true
}

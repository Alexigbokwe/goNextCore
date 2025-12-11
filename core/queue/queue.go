package queue

// TaskQueue defines the contract for background job processing.
// Implement this interface using your preferred provider (Redis, RabbitMQ, Kafka, etc.).
type TaskQueue interface {
	Enqueue(typeName string, payload interface{}, opts ...interface{}) error
	RegisterHandler(typeName string, handler func(payload []byte) error)
	Start()
	Shutdown()
}

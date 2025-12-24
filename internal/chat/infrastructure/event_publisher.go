package infrastructure

type EventPublisher interface {
    Publish(event any)
}

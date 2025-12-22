package usecase

type EventPublisher interface {
	Publish(event any)
}

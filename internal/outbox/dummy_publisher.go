package outbox

import "context"

type DummyPublisher struct{}

func (dp *DummyPublisher) Publish(
	ctx context.Context,
	topic string,
	key string,
	payload []byte,
) error {
	return nil
}

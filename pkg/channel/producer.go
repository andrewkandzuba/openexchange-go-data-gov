package channel

type Producer interface {
	Send(topic string, data string) error
}

func Send(p Producer, topic string, data string) error {
	return p.Send(topic, data)
}

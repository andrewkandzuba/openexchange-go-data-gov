package channel

type Consumer interface {
	Subscribe(topic string, listener *Listener)
}

func Subscribe(c Consumer, topic string, listener *Listener) {
	c.Subscribe(topic, listener)
}

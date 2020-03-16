package channel

type Listener interface {
	OnData(topic string, data string)
	OnError(err error)
}

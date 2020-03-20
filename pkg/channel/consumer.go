package channel

import "io"

type Consumer interface {
	io.Closer
	Subscribe(topic string, listener *Listener)
}
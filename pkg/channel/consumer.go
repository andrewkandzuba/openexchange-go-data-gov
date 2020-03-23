package channel

import "io"

type Consumer interface {
	io.Closer
	Subscribe(topic string, fn func (value interface{}))
}
package channel

import "io"

type Producer interface {
	io.Closer
	Send(topic string, data string) error
}

package closer

import "context"

type Closer struct {
	cancelFunc context.CancelFunc
}

func NewCloser(cancelFunc context.CancelFunc) *Closer {
	return &Closer{cancelFunc: cancelFunc}
}

func (c *Closer) Close() error {
	c.cancelFunc()
	return nil
}

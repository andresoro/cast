package cast

import "sync"

// Broadcast implements an interface for a fan out pattern
// the goal is to have one input source pass data to multiple output sources
// with each source getting the same data
type Broadcast interface {
	Start()
	Close()
	Add(chan []byte)
}

type cast struct {
	input  <-chan []byte
	output []chan<- []byte
	end    chan struct{}
	add    chan chan []byte
	mu     sync.Mutex
}

func New() Broadcast {
	return &cast{}
}

func (c *cast) Start() {

	go func() {
		for {
			select {
			// handle incoming data
			case data := <-c.input:
				for _, ch := range c.output {
					ch <- data
				}
			case ch := <-c.add:
				c.output = append(c.output, ch)
			// handle end signal
			case <-c.end:
				for _, ch := range c.output {
					close(ch)
				}
			}

		}
	}()

}

func (c *cast) Add(ch chan []byte) {
	c.add <- ch
}

func (c *cast) Close() {
	c.end <- struct{}{}
}

package cast

import "sync"

// Relay implements an interface for a fan out pattern
// the goal is to have one input source pass data to multiple output recievers
// with each reciever getting the same data
type Relay interface {
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

// New returns a broacast interface from an input channel
func New(input <-chan []byte) Relay {
	return &cast{
		input:  input,
		output: make([]chan<- []byte, 0),
		end:    make(chan struct{}),
		add:    make(chan chan []byte),
	}
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
				return
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

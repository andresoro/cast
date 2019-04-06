package cast

import "sync"

// Relay implements an interface for a fan out pattern
// the goal is to have one input source pass data to multiple output recievers
// with each reciever getting the same data
type Relay interface {
	Start()
	Close()
	New() chan []byte
}

type cast struct {
	input   <-chan []byte
	outputs []chan<- []byte
	end     chan struct{}
	add     chan chan []byte
	running bool
	mu      sync.Mutex
}

// New returns a broacast interface from an input channel
func New(input <-chan []byte) Relay {
	return &cast{
		input:   input,
		outputs: make([]chan<- []byte, 0),
		end:     make(chan struct{}),
		add:     make(chan chan []byte),
	}
}

func (c *cast) Start() {
	c.running = true
	go func() {
		for {
			select {
			// handle incoming data
			case data := <-c.input:
				for _, ch := range c.outputs {
					ch <- data
				}
			case ch := <-c.add:
				c.outputs = append(c.outputs, ch)
			// handle end signal
			case <-c.end:
				for _, ch := range c.outputs {
					close(ch)
				}
				c.running = false
				return
			}

		}
	}()

}

// New returns a recieve channel for a consumer
func (c *cast) New() chan []byte {

	ch := make(chan []byte, 1)

	if c.running {
		c.add <- ch
	} else {
		c.outputs = append(c.outputs, ch)
	}

	return ch

}

func (c *cast) Close() {
	c.end <- struct{}{}
}

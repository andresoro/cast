# Cast
An interface to satisfy a one to many message passing pattern (fan out). This problem arose when I wanted mulitple one-to-many messaging patterns in a backend server. Cast allows for an input channel to fan-out to multiple output channels with each output channel recieving the same value. For example, this can be used as a dispatcher to relay messages from a process in a server to multiple clients.


## Example usage

```golang
package main

import "github.com/andresoro/cast"

func main() {

    // where our values come from
	producer := make(chan []byte)

	// output channels
	outputs := make([]chan []byte, 0)
	for i := 0; i < 5; i++ {
		outputs = append(outputs, make(chan []byte, 1))
	}

	// init relay, add channels
	relay := New(producer)
	for _, ch := range outputs {
		relay.Add(ch)
	}
    // start relay, now listening to producer channel
	relay.Start()

	// value added to producer will be sent to every output channel
	producer <- []byte("Hello World")

}
```


### Disclaimer

Wrote this up quickly in an effort to avoid writing something like this multiple times for my webserver. Better tests and more features will probably be added later. 

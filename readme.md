# Cast
An interface to satisfy a one to many message passing pattern (fan out). This problem arose when I wanted mulitple one-to-many messaging patterns in a backend server. Cast allows for an input channel to fan-out to multiple output channels with each output channel recieving the same value. For example, this can be used as a dispatcher to relay messages from a process in a server to multiple clients.


## Example usage

```
package main

import "github.com/andresoro/cast"

func main() {

    // our input channel
    producer := make(chan []byte)

    // broadcast struct, run Start() before adding chans
    cast := cast.New(producer)
    cast.Start()


    // output channels
    chans := make([]chan []byte, 0)
    for i := 0; i < 5; i++ {
        // init chan
        output := make(chan []byte, 1)
        chans = chans.append(output)

        // add chan to broadcast interface
        cast.Add(output)
    }

    // this value is recieved by all channels
    producer <- []byte("Hello World)

}
```


### Disclaimer

Wrote this up quickly in an effort to avoid writing something like this multiple times for my webserver. Better tests and more features will probably be added later. 
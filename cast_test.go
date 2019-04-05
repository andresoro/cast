package cast

import (
	"testing"
)

func TestOutputs(t *testing.T) {

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

	relay.Start()

	// value added to producer
	producer <- []byte("Hello World")

	for i, ch := range outputs {
		val := <-ch
		if string(val) != "Hello World" {
			t.Errorf("Channel %d not recieving value", i)
		}

	}
}

func TestAdding(t *testing.T) {
	// adding before or after calling Start should not change behavior

	producer := make(chan []byte)

	output1 := make(chan []byte, 1)
	output2 := make(chan []byte, 1)

	relay := New(producer)

	relay.Add(output1)
	relay.Start()
	relay.Add(output2)

	producer <- []byte("Testing")

	val1 := <-output1
	if string(val1) != "Testing" {
		t.Error("Output channel added before start not recieving value")
	}

	val2 := <-output2
	if string(val2) != "Testing" {
		t.Error("Output channel after Start not recieving value")
	}

}

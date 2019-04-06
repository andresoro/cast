package cast

import (
	"testing"
)

func TestOutputs(t *testing.T) {

	// where our values come from
	producer := make(chan []byte)

	// output channels
	outputs := make([]chan []byte, 0)

	// init relay, add channels
	relay := New(producer)
	for i := 0; i < 10; i++ {
		outputs = append(outputs, relay.New())
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

	relay := New(producer)

	output1 := relay.New()

	relay.Start()

	output2 := relay.New()

	producer <- []byte("Value")

	val1 := <-output1
	if string(val1) != "Value" {
		t.Error("Error Adding channel 1")
	}

	val2 := <-output2
	if string(val2) != "Value" {
		t.Error("Error Adding channel 2")
	}
}

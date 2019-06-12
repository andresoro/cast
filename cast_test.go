package cast

import (
	"fmt"
	"sync"
	"testing"
	"time"
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

func TestBlockingWrite(t *testing.T) {

	producer := make(chan []byte)

	relay := New(producer)
	relay.Start()

	output1 := relay.New()

	output1 <- []byte("Value")

	producer <- []byte("Value 2")

	out := <-output1

	if string(out) != "Value" {
		t.Error("Wrong value")
	}

}

// Test should send a newly added channel the last data
func TestLastSent(t *testing.T) {
	producer := make(chan []byte)

	relay := New(producer)

	relay.Start()

	producer <- []byte("Test data")

	output := relay.New()

	val := <-output
	if string(val) != "Test data" {
		t.Error("Output channel not receiving last sent data")
	}

	output2 := relay.New()

	val = <-output2
	if string(val) != "Test data" {
		t.Error("Second output channel not receiving last sent data")
	}
}

func TestClose(t *testing.T) {
	producer := make(chan []byte)

	relay := New(producer)

	output1 := relay.New()
	output2 := relay.New()
	output3 := relay.New()
	output4 := relay.New()

	relay.Start()
	relay.Close()

	<-output1
	<-output2
	<-output3
	<-output4
}

// Tests slow and fast receivers
// Producer will write data at 2hz
// Slow receiver reads at 1hz
// Fast receiver reads at 4hz
// Test passes when 10 messages are received by either channel
func TestSlowReceiverCase(t *testing.T) {

	producer := make(chan []byte)

	relay := New(producer)

	// slow reciver, reads at .75 seconds
	output1 := relay.New()

	// faster reciever, can keep up, reads at .25 seconds
	output2 := relay.New()

	relay.Start()
	t.Log("Starting test producer")

	messagesSent := 0

	// func to output messages through producer
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			t.Logf("Message %d sent", messagesSent)

			producer <- []byte(fmt.Sprintf("Testing %d", messagesSent))
			messagesSent++
		}
	}()

	// message counters for both output channels
	messagesCh1 := 0
	messagesCh2 := 0

	var wg sync.WaitGroup
	wg.Add(2)

	// func to handle messages to output1
	go func() {
	loop1:
		for {
			select {
			case <-output1:
				time.Sleep(1000 * time.Millisecond)
				messagesCh1++
				if messagesCh1 >= 10 || messagesCh2 >= 10 {
					break loop1
				}
			case <-time.After(4 * time.Second):
				t.Fatalf("Failed slow receiver")
				t.Logf("Messages received on #1: %d", messagesCh1)
				t.Logf("Messages received on #2: %d", messagesCh2)
				break loop1
			}
		}
		wg.Done()
	}()

	// func to handle messages to output2
	go func() {
	loop2:
		for {
			select {
			case <-output2:
				time.Sleep(250 * time.Millisecond)
				messagesCh2++
				if messagesCh1 >= 10 || messagesCh2 >= 10 {
					break loop2
				}
			case <-time.After(4 * time.Second):
				t.Fatal("Failed slow receiver")
				t.Logf("Messages received on #1: %d", messagesCh1)
				t.Logf("Messages received on #2: %d", messagesCh2)
				break loop2
			}
		}
		wg.Done()
	}()

	wg.Wait()
	t.Logf("Messages received on #1: %d", messagesCh1)
	t.Logf("Messages received on #2: %d", messagesCh2)

}

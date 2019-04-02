package cast

import "testing"

func TestCast(t *testing.T) {

	// where our values come from
	producer := make(chan []byte)

	outputCh1 := make(chan []byte, 1)
	outputCh2 := make(chan []byte, 1)
	outputCh3 := make(chan []byte, 1)

	bcast := New(producer)
	bcast.Start()
	bcast.Add(outputCh1)
	bcast.Add(outputCh2)
	bcast.Add(outputCh3)

	val := []byte("Hello World")

	producer <- val

	msg1 := <-outputCh1
	if string(msg1) != string(val) {
		t.Error("Not receiving value on ch1")
	}

	msg2 := <-outputCh2
	if string(msg2) != string(val) {
		t.Error("Not recieving value on ch2")
	}

	msg3 := <-outputCh3
	if string(msg3) != string(val) {
		t.Error("Not recieving value on ch3")
	}

}

package utils

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	//see: https://github.com/stretchr/testify/assert
)

// TestPrintFields test PrintFields
func TestPrintFields(t *testing.T) {

	t.Skipf("Skipped due to reflection issues")

	// create an anonymous struct for testing
	data := struct {
		Field1 string
		Field2 string
	}{
		Field1: "aaaa",
		Field2: "bbbb",
	}

	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintFields("Anonymous", data)

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	// testing the output
	assert.Contains(t, out, "aaaa")
}

// TestSprintFields test SprintFields
func TestSprintFields(t *testing.T) {

	t.Skipf("Skipped due to reflection issues")

	// create an anonymous struct for testing
	data := struct {
		Field1 string
		Field2 string
	}{
		Field1: "aaaa",
		Field2: "bbbb",
	}

	out := SprintFields("Anonymous", data)

	// testing the output
	assert.Contains(t, out, "aaaa")
}

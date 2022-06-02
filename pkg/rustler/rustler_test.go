package rustler

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func TestRustler(t *testing.T) {

	// Given
	var target *Target
	if tgt, err := NewTarget("localhost:8500"); err == nil {
		target = tgt
	} else {
		t.Fatalf("Unable to connect to target. Is port-forwarding set up on localhost:8500?: %v", err)
	}

	// When

	var output []string
	if proc, ok := target.Rustle("echo", []string{"hello, world"}); ok {
		if pout, e := proc.Output(context.Background()); e == nil {
			for o := range pout {
				output = append(output, o)
			}
		} else {
			t.Errorf("unable to read output: %v", e)
		}
	} else {
		t.Error("error running command")
	}

	// Then
	if assert.Len(t, output, 1) {
		assert.Equal(t, "hello, world", output[0])
	}
}

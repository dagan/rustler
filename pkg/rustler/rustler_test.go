/*
 * Copyright © 2022 Dagan Henderson <dagan@techdagan.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package rustler

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func TestRustler(t *testing.T) {

	if testing.Short() {
		t.Skipf("Skipping tests for short mode")
	}

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

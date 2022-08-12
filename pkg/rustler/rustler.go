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
	"crypto/rand"
	"github.com/longhorn/longhorn-instance-manager/pkg/api"
	"github.com/longhorn/longhorn-instance-manager/pkg/client"
	"golang.org/x/net/context"
)

type Target struct {
	client *client.ProcessManagerClient
}

// NewTarget creates a new Rustler Target with the specified address.
// A check connection will be attempted and is successful if err is nil.
func NewTarget(addr string) (t *Target, err error) {
	c := client.NewProcessManagerClient(addr)
	_, err = c.VersionGet() // quick client/connection test
	if err == nil {
		t = &Target{
			client: c,
		}
	}
	return
}

// Rustle creates (executes) a remote process on the Target using the provided
// command and arguments.
// The returned Process can be used to retrieve command output.
func (t *Target) Rustle(command string, args []string) (process *Process, ok bool) {

	// Generate a random process ID
	var id string
	if id, ok = generateId(8); !ok {
		return
	}

	// Create the remote process
	if p, e := t.client.ProcessCreate(id, command, 0, args, []string{}); e == nil {
		process = &Process{
			client:  t.client,
			process: p,
		}
		ok = true
	}
	return
}

type Process struct {
	client  *client.ProcessManagerClient
	process *api.Process
}

// Output
func (p *Process) Output(ctx context.Context) (<-chan string, error) {

	out := make(chan string)
	var err error

	// Get the Longhorn log stream
	var log *api.LogStream
	if log, err = p.client.ProcessLog(p.process.Name); err != nil {
		return nil, err
	}

	// Buffer logstream messages until error (which occurs after LogStream.Close())
	buf := make(chan string, 1)
	go func(logs *api.LogStream, out chan<- string) {
		defer close(out)
		for {
			if msg, e := log.Recv(); e == nil {
				out <- msg
			} else {
				return
			}
		}
	}(log, buf)

	// Transfer log messages until the logstream closes or the context is done.
	go func(ctx context.Context, log *api.LogStream, in <-chan string, out chan<- string) {
		defer close(out)
		done := ctx.Done()
		for {
			select {
			case _ = <-done:
				_ = log.Close()
				done = nil
			case msg, ok := <-in:
				if !ok {
					return
				}
				out <- msg
			}
		}
	}(ctx, log, buf, out)

	return out, err
}

var idAlphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// generateId generates unnecessarily secure random IDs of the specified length.
// Given an alphabet of 52 characters (a-z and A-Z), this function generates:
//  - 45 bits of entropy for n = 8
//  - 64 bits of entropy for n = 12
//  - 128 bits of entropy for n = 23
//  - 256 bits of entropy for n = 45
func generateId(n int) (string, bool) {

	id := make([]rune, n)
	src := make([]byte, n)

	// Generate a cryptographically secure random string by only using byte values 0-207
	for pos := 0; pos < n; {

		// After the first pass, restrict src len to required size
		src = src[:n-pos]

		// Get entropy
		if _, e := rand.Read(src); e != nil {
			return "", false // Never return a partial ID
		}

		// Generate ID from entropy
		for i := 0; i < len(src); i++ {
			// Ensure even distribution across idAlphabet by discarding byte values > 207
			if src[i] < 208 {
				j := src[i] % 52
				id[pos] = idAlphabet[j]
				pos++
			}
		}
	}

	return string(id), true
}

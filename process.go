/*
Copyright 2018 Nathan Frederick Smith.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plumb

import (
	"errors"
	"reflect"
)

// Process is a function
type Process struct {
	reflect.Value // A reflect Function
	in            []Channel
	out           []Channel
	done          chan struct{} // Global kill channel
	workers       int
}

// NewProcess creates a process from a function
func NewProcess(f interface{}) (Process, error) {
	t := reflect.TypeOf(f)
	v := reflect.ValueOf(f)

	if t.Kind() != reflect.Func {
		return Process{}, errors.New("f must be a function")
	}

	// Make input channel
	in := make([]Channel, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		c := NewChannel(t.In(i), 1)
		in = append(in, c)
	}

	// Make output channels
	out := make([]Channel, 0)
	for j := 0; j < t.NumOut(); j++ {
		c := NewChannel(t.Out(j), 1)
		out = append(out, c)
	}

	done := make(chan struct{})
	return Process{v, in, out, done, 0}, nil
}

// Run spawn a worker for this process. Can by called
// as many times as you which to scale up the Process
func (p *Process) Run() {
	p.workers++
	go func() {
		inputs := make([]reflect.Value, p.Type().NumIn())
		outputs := make([]reflect.Value, p.Type().NumOut())
		var ok bool

		for {
			select {
			case <-p.done:
				p.workers--
				return
			default:
				// Read input from input channels
				for i := range inputs {
					inputs[i], ok = p.in[i].Recv()
					if ok == false {
						// TODO: Propagate the closure to output channels?
						p.workers--
						return
					}
				}

				// Compute
				outputs = p.Call(inputs)

				// Distribute output to output channels
				for j, output := range outputs {
					p.out[j].Send(output)
				}
			}
		}
	}()
}

// NumWorkers returns the number of worker in this
// Process
func (p *Process) NumWorkers() int {
	return p.workers
}

// SetWorkers scales or shrinks the number of workers
// running this Process
func (p *Process) SetWorkers(target int) {
	for p.workers < target {
		p.Run()
	}
	for p.workers > target {
		p.done <- struct{}{}
	}
}

// Connect links the i'th output of p to the j'th input
// of pr
func (p *Process) Connect(pr *Process, i, j int) {
	pr.in[j] = p.out[i]
	return
}

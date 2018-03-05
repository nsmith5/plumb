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
	F       reflect.Value // A reflect Function
	In      []Channel
	Out     []Channel
	Done    chan struct{} // Global kill channel
	Workers int
}

// NewProcess creates a process from a function
func NewProcess(f interface{}) (p Process, err error) {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		err = errors.New("f must by a function. Supplied " + t.String())
		return
	}

	p.F = reflect.ValueOf(f)
	// Init input channels
	for i := 0; i < t.NumIn(); i++ {
		c := NewChannel(t.In(i), 1)
		p.In = append(p.In, c)
	}
	// Init output channels
	for j := 0; j < t.NumOut(); j++ {
		c := NewChannel(t.Out(j), 1)
		p.Out = append(p.Out, c)
	}

	p.Done = make(chan struct{})
	p.Workers = 0
	return
}

// recv aggragates receiving from all input channels
func (p *Process) recv() ([]reflect.Value, bool) {
	var allOk, ok bool
	inputs := make([]reflect.Value, p.F.Type().NumIn())

	for i := range inputs {
		inputs[i], ok = p.In[i].Recv()
		allOk = ok || allOk
	}

	return inputs, allOk
}

// send aggregates outputs and sends them on their
// respective output channels
func (p *Process) send(outputs []reflect.Value) {
	for i, output := range outputs {
		p.Out[i].Send(output)
	}
}

// Run spawns a worker for this process. Can by called
// as many times as you which to scale up the Process
func (p *Process) Run() {
	p.Workers++
	go func() {
		var inputs, outputs []reflect.Value
		var ok bool

		for {
			select {
			// If there is a kill signal in the
			// Done channel, decrement Workers and
			// get out.
			case <-p.Done:
				p.Workers--
				return

			// This is the normal course of action:
			// (1) Receive all arguments
			// (2) Run computation
			// (3) Distribute outputs
			default:
				inputs, ok = p.recv()
				if !ok {
					_ = ok
					// TODO: What should happen here?
				}
				outputs = p.F.Call(inputs)
				p.send(outputs)
			}
		}
	}()
}

// NumWorkers returns the number of worker in this
// Process
func (p *Process) NumWorkers() int {
	return p.Workers
}

// SetWorkers scales or shrinks the number of workers
// running this Process
func (p *Process) SetWorkers(target int) {
	for p.Workers < target {
		p.Run()
	}
	for p.Workers > target {
		p.Done <- struct{}{}
	}
}

// Connect links the i'th output of p to the j'th input
// of pr
func (p *Process) Connect(pr *Process, i, j int) {
	pr.In[j] = p.Out[i]
	return
}

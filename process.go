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
}

// NewProcess creates a process from a function
func NewProcess(f interface{}) (Process, error) {
	t := reflect.TypeOf(f)
	v := reflect.ValueOf(f)

	if t.Kind() != reflect.Func {
		return Process{}, errors.New("f must be a function")
	}

	// Make input channel
	in := make([]Channel, 0)
	for i := 0; i < t.NumIn(); i++ {
		ctype := reflect.ChanOf(reflect.BothDir, t.In(i))
		c := Channel{reflect.MakeChan(ctype, 1)}
		in = append(in, c)
	}

	// Make output channels
	out := make([]Channel, 0)
	for j := 0; j < t.NumOut(); j++ {
		ctype := reflect.ChanOf(reflect.BothDir, t.Out(j))
		c := Channel{reflect.MakeChan(ctype, 1)}
		out = append(out, c)
	}

	return Process{v, in, out}, nil
}

// Run spawn a worker for this process. Can by called
// as many times as you which to scale up the Process
func (p *Process) Run() {
	go func() {
		inputs := make([]reflect.Value, p.Type().NumIn())
		outputs := make([]reflect.Value, p.Type().NumOut())
		var ok bool

		for {
			// Read input from input channels
			for i := range inputs {
				inputs[i], ok = p.in[i].Recv()
				if ok == false {
					for _, out := range p.out {
						out.Close()
					}
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
	}()
}

// Connect links the i'th output of p to the j'th input
// of pr
func (p *Process) Connect(pr *Process, i, j int) {
	pr.in[j] = p.out[i]
	return
}

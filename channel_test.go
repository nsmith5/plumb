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
	"log"
	"testing"
	"time"
)

func source() string {
	return "Hello World!"
}

func logAndPass(in string) string {
	log.Println(in)
	time.Sleep(1 * time.Second)
	return in
}

func sink(in string) {
	return
}

func TestRun(t *testing.T) {
	src, _ := NewProcess(source)
	middle, _ := NewProcess(logAndPass)
	snk, _ := NewProcess(sink)

	src.Connect(&middle, 0, 0)
	middle.Connect(&snk, 0, 0)

	src.Run()
	middle.Run()
	snk.Run()

	time.Sleep(2 * time.Second)
	middle.SetWorkers(2)
	time.Sleep(2 * time.Second)
	middle.SetWorkers(1)
	time.Sleep(2 * time.Second)
}

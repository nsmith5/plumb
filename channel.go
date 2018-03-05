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
	"reflect"
)

// Channel is the basic connection between
type Channel struct {
	reflect.Value // Require Channel.Kind() == Chan
}

// NewChannel constructs a new Channel with buffer length
// `buffer` and content type `t`.
func NewChannel(t reflect.Type, buffer int) Channel {
	ctype := reflect.ChanOf(reflect.BothDir, t)
	return Channel{reflect.MakeChan(ctype, buffer)}
}

package plumb

import (
	"io"
	"net/http"
)

func (p *Process) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World!")
}

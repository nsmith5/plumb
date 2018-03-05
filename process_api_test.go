package plumb

import (
	"net/http"
	"testing"
)

func TestProcessAPI(t *testing.T) {
	src, _ := NewProcess(source)
	middle, _ := NewProcess(logAndPass)
	snk, _ := NewProcess(sink)

	src.Connect(&middle, 0, 0)
	middle.Connect(&snk, 0, 0)

	src.Run()
	middle.Run()
	snk.Run()

	mux := http.NewServeMux()
	mux.Handle("/middle", &middle)
	go http.ListenAndServe(":8080", mux)

	http.Get("http://localhost:8080/middle")
	return
}

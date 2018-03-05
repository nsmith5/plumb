package plumb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (p *Process) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"workers": p.Workers,
		})

	case "PATCH", "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(err)
			return
		}
		_map := make(map[string]int)
		if err = json.Unmarshal(body, &_map); err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(err)
			return
		}

		p.SetWorkers(_map["workers"])
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode("Method not allowed")
	}
}

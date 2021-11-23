package main

import (
	"net/http"
	"testing"
)

func TestDdos(t *testing.T) {
	// create http server for testing
	go http.ListenAndServe(":8000", nil)
	ddoser, error := NewDDoser("GET", "http://localhost:8000")

	if error != nil {
		t.Error(error)
	}

	ddoser.Do(getUserAgents(10), 5)
}

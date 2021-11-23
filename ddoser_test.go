package main

import (
	"net/http"
	"testing"
)

func TestDdos(t *testing.T) {
	// create http server for testing
	go http.ListenAndServe(":8000", nil)
	ddoser, error := NewDDoser("abc", "GET")

	if error != nil {
		t.Error(error)
	}

	go ddoser.Do(getUserAgents(10), 5)
}

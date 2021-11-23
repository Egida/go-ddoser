package main

import (
	"net/http"
	"testing"
)

func TestDdos(t *testing.T) {
	req, error := http.NewRequest("GET", "https://github.com", nil)

	if error != nil {
		t.Error(error)
	}

	ddoser, error := NewDDoser(req, "443", "72.210.208.101:4145")

	if error != nil {
		t.Error(error)
	}

	ddoser.Do(getUserAgents(10), 1)
}

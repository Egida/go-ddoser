package main

import (
	"testing"
)

func TestDdos(t *testing.T) {
	ddoser, error := NewDDoser("GET", "https://github.com:443", "72.210.208.101:4145")

	if error != nil {
		t.Error(error)
	}

	ddoser.Do(getUserAgents(10), 5)
}

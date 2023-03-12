package main

import "testing"

func TestItAll(t *testing.T) {
	if err := run(); err != nil {
		t.Fatal(err)
	}
}

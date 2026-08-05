package main

import "testing"

func TestNewServer(t *testing.T) {
	server, err := newServer(":0")
	if err != nil {
		t.Fatalf("newServer() error = %v", err)
	}
	if server.Addr != ":0" {
		t.Fatalf("server.Addr = %q, want %q", server.Addr, ":0")
	}
}

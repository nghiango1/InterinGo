package server_test

import (
	"interingo/pkg/server"
	"testing"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		{"Server can be created"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := server.NewServer()
			// TODO: update the condition below to compare got with tt.want.
			if got == nil {
				t.Errorf("NewServer() failed to create server implementation")
			}
		})
	}
}

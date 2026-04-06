package object_test

import (
	"interingo/pkg/object"
	"testing"
)

func TestArray_Inspect(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		arr  object.Array
		want string
	}{
		{
			"validate",
			object.Array{
				Value: []object.Object{
					&object.Integer{Value: 1},
				},
			},
			"[1]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			got := tt.arr.Inspect()
			// TODO: update the condition below to compare got with tt.want.
			if tt.want != got {
				t.Errorf("Inspect() = %v, want %v", got, tt.want)
			}
		})
	}
}

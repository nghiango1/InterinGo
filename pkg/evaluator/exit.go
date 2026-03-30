package evaluator

import "interingo/pkg/object"

func SystemExit(n int) (object.Object, *SignalExitInfo) {
	return nil, &SignalExitInfo{
		ReturnSignal: n,
	}
}

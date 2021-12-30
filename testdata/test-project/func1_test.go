package testproject

import "testing"

func TestFunc1(t *testing.T) {
	tests := map[string]struct {
		bool1 bool
		bool2 bool
	}{
		"bool1": {bool1: true},
	}
	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			Func1(tt.bool1, tt.bool2)
		})
	}
}

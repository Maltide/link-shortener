package hash

import (
	"testing"

	"go.uber.org/zap"
)

func TestHashAndRedirectHash(t *testing.T) {
	log := zap.NewNop().Sugar()
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{"zero", 0, true},
		{"one", 1, false},
		{"sixty two", 62, false},
		{"large", 123456, false},
		{"negative", -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			short, err := hash(tt.id, log)
			if (err != nil) != tt.wantErr {
				t.Errorf("%q: unexpected err: %v, wantErr=%v", tt.name, err, tt.wantErr)
				return
			}
			if err == nil {
				gotId, err := redirectHash(short, log)
				if err != nil {
					t.Errorf("%q: redirectHash error: %v", tt.name, err)
				}
				if gotId != tt.id {
					t.Errorf("%q: got id %d, want %d", tt.name, gotId, tt.id)
				}
			}
		})
	}
}

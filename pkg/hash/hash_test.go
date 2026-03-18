package hash

import (
	"database/sql"
	"testing"

	"github.com/Maltide/link-shortener/pkg/config"
	"github.com/Maltide/link-shortener/pkg/helpers"
	"go.uber.org/zap"
)

func TestHashAndRedirectHash(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	log := zap.NewNop().Sugar()

	connStr := helpers.ConnStr(cfg)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

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

	pair_counter = 1000000

	short, err := Hash(1, log, db)
	if err != nil || short == "" {
		t.Errorf("hash failed after cleanup: err=%v, short=%q", err, short)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			short, err := Hash(tt.id, log, db)
			if (err != nil) != tt.wantErr {
				t.Errorf("%q: unexpected err: %v, wantErr=%v", tt.name, err, tt.wantErr)
				return
			}
			if err == nil {
				gotId, err := RedirectHash(short, log)
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

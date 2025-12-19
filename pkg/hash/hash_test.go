package hash

import (
	"testing"

	"go.uber.org/zap"
)

func TestHash(t *testing.T) {
	log := zap.NewNop().Sugar()
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid URL",
			input:   "https://example.com",
			want:    "a6",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "URL with path",
			input:   "https://example.com/page1",
			want:    "rd",
			wantErr: false,
		},
		{
			name:    "URL with params",
			input:   "https://example.com/?a=1",
			want:    "tf",
			wantErr: false,
		},
		{
			name:    "URL with query params",
			input:   "https://example.com/?a=1&b=2",
			want:    "4e",
			wantErr: false,
		},
		{
			name:    "long URL",
			input:   "https://example.com/this/is/a/very/long/url/with/many/segments/and/parameters?foo=bar&baz=qux",
			want:    "S4",
			wantErr: false,
		},
		{
			name:    "URL without scheme",
			input:   "example.com/page2",
			want:    "da",
			wantErr: false,
		},
		{
			name:    "URL with special characters",
			input:   "https://example.com/!@#$%^&*()_+",
			want:    "38",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := hash(tt.input, log)
			if (err != nil) != tt.wantErr {
				t.Errorf("%q:unexpected err:%v, wantErr =%v", tt.name, err, tt.wantErr)
				return
			}
			if result != tt.want {
				t.Errorf("%q: got %q, want %q", tt.name, result, tt.want)
			}
		})
	}
}

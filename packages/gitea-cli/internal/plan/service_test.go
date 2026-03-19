package plan

import (
	"testing"
)

// go test -v -run TestGenNextId ./internal/plan/
func TestGenNextId(t *testing.T) {
	tests := []struct {
		name     string
		existIds []string
		want     string
		wantErr  bool
	}{
		{
			name:     "start",
			existIds: []string{},
			want:     "100",
			wantErr:  false,
		},
		{
			name:     "end",
			existIds: []string{"900"},
			want:     "",
			wantErr:  true,
		},
		{
			name:     "a00",
			existIds: []string{"100"},
			want:     "200",
			wantErr:  false,
		},
		{
			name:     "a00 pattern 1",
			existIds: []string{"100", "200"},
			want:     "300",
			wantErr:  false,
		},
		{
			name:     "a00 pattern 2",
			existIds: []string{"100", "200", "300"},
			want:     "400",
			wantErr:  false,
		},
		{
			name:     "a00 pattern 3",
			existIds: []string{"100", "300"},
			want:     "400",
			wantErr:  false,
		},
		{
			name:     "ab0 pattern 1",
			existIds: []string{"100", "110"},
			want:     "200",
			wantErr:  false,
		},
		{
			name:     "ab0 pattern 2",
			existIds: []string{"100", "210"},
			want:     "300",
			wantErr:  false,
		},
		{
			name:     "ab0 pattern 3",
			existIds: []string{"100", "890"},
			want:     "900",
			wantErr:  false,
		},
		{
			name:     "abc pattern 1",
			existIds: []string{"123", "125"},
			want:     "200",
			wantErr:  false,
		},
		{
			name:     "abc pattern 2",
			existIds: []string{"123", "225"},
			want:     "300",
			wantErr:  false,
		},
		{
			name:     "abc pattern 3",
			existIds: []string{"123", "899"},
			want:     "900",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := genNextId(tt.existIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("genNextId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("genNextId() = %v, want %v", got, tt.want)
			}
		})
	}
}

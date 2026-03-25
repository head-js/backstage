package plan

import (
	"testing"
)

// go test -v -run TestGenNextTierId ./internal/plan/
func TestGenNextTierId(t *testing.T) {
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
			got, err := genNextTierId(tt.existIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("genNextTierId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("genNextTierId() = %v, want %v", got, tt.want)
			}
		})
	}
}

// go test -v -run TestGenNextBlockId ./internal/plan/
func TestGenNextBlockId(t *testing.T) {
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
			existIds: []string{"990"},
			want:     "",
			wantErr:  true,
		},
		{
			name:     "ab0 to ab10",
			existIds: []string{"100"},
			want:     "110",
			wantErr:  false,
		},
		{
			name:     "ab0 pattern 1",
			existIds: []string{"100", "110"},
			want:     "120",
			wantErr:  false,
		},
		{
			name:     "ab0 pattern 2",
			existIds: []string{"100", "120"},
			want:     "130",
			wantErr:  false,
		},
		{
			name:     "ab0 pattern 3",
			existIds: []string{"100", "180"},
			want:     "190",
			wantErr:  false,
		},
		{
			name:     "ab0 to next hundred",
			existIds: []string{"190"},
			want:     "200",
			wantErr:  false,
		},
		{
			name:     "abc pattern 1",
			existIds: []string{"123", "125"},
			want:     "130",
			wantErr:  false,
		},
		{
			name:     "abc pattern 2",
			existIds: []string{"123", "225"},
			want:     "230",
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
			got, err := genNextBlockId(tt.existIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("genNextBlockId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("genNextBlockId() = %v, want %v", got, tt.want)
			}
		})
	}
}

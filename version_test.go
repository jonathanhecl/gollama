package gollama

import (
	"context"
	"testing"
)

func TestGollama_Version(t *testing.T) {
	tests := []struct {
		name    string
		c       *Gollama
		want    string
		wantErr bool
	}{
		{
			name:    "Version",
			c:       New("llama3.2"),
			want:    "0.4.2",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Version(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Gollama.Version() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Gollama.Version() = %v, want %v", got, tt.want)
			}
		})
	}
}

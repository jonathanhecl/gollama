package gollama

import (
	"context"
	"testing"
)

func TestGollama_Embedding(t *testing.T) {
	type args struct {
		Prompt string
	}
	tests := []struct {
		name    string
		c       *Gollama
		args    args
		wantLen int
		wantErr bool
	}{
		{
			name:    "Embedding",
			c:       New("llama3.2"),
			args:    args{Prompt: "hello"},
			wantLen: 3072,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Embedding(context.Background(), tt.args.Prompt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Gollama.Embedding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("Gollama.Embedding() = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestCosenoSimilarity(t *testing.T) {
	type args struct {
		vector1 []float64
		vector2 []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "CosenoSimilarity",
			args: args{vector1: []float64{1, 2, 3}, vector2: []float64{4, 5, 6}},
			want: 0.9746318461970762,
		},
		{
			name: "CosenoSimilarity",
			args: args{vector1: []float64{1, 1, 1}, vector2: []float64{1, 1, 1}},
			want: 1.0000000000000002,
		},
		{
			name: "CosenoSimilarity",
			args: args{vector1: []float64{1, 0, 0}, vector2: []float64{0, 1, 0}},
			want: 0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CosenoSimilarity(tt.args.vector1, tt.args.vector2); got != tt.want {
				t.Errorf("CosenoSimilarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

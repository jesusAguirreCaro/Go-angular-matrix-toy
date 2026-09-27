package main

import (
	"math"
	"testing"
)

func TestDeterminant(t *testing.T) {
	cases := []struct {
		name string
		in   [][]float64
		want float64
	}{
		{"1x1", [][]float64{{5}}, 5},
		{"2x2", [][]float64{{1, 2}, {3, 4}}, -2},
		{"identity3", [][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}, 1},
		{"3x3", [][]float64{{6, 1, 1}, {4, -2, 5}, {2, 8, 7}}, -306},
		{"singular", [][]float64{{1, 2}, {2, 4}}, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Determinant(tc.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tc.want) > 1e-9 {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDeterminantRejectsNonSquare(t *testing.T) {
	if _, err := Determinant([][]float64{{1, 2, 3}, {4, 5, 6}}); err == nil {
		t.Fatal("expected error for non-square matrix")
	}
}

package main

import (
	"math/big"
	"testing"
)

func mustRat(s string) *big.Rat {
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		panic("bad rat: " + s)
	}
	return r
}

func TestDeterminantRat(t *testing.T) {
	cases := []struct {
		name string
		in   [][]string
		want string
	}{
		{"1x1", [][]string{{"5"}}, "5"},
		{"2x2", [][]string{{"1", "2"}, {"3", "4"}}, "-2"},
		{"identity3", [][]string{{"1", "0", "0"}, {"0", "1", "0"}, {"0", "0", "1"}}, "1"},
		{"3x3", [][]string{{"6", "1", "1"}, {"4", "-2", "5"}, {"2", "8", "7"}}, "-306"},
		{"singular", [][]string{{"1", "2"}, {"2", "4"}}, "0"},
		{"rational", [][]string{{"1/2", "0"}, {"0", "2"}}, "1"},
		{"float-input", [][]string{{"0.5", "0"}, {"0", "2"}}, "1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := make([][]*big.Rat, len(tc.in))
			for i, row := range tc.in {
				m[i] = make([]*big.Rat, len(row))
				for j, s := range row {
					m[i][j] = mustRat(s)
				}
			}
			got, err := DeterminantRat(m)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.RatString() != tc.want {
				t.Fatalf("got %s, want %s", got.RatString(), tc.want)
			}
		})
	}
}

func TestDeterminantRatRejectsNonSquare(t *testing.T) {
	m := [][]*big.Rat{
		{mustRat("1"), mustRat("2"), mustRat("3")},
		{mustRat("4"), mustRat("5"), mustRat("6")},
	}
	if _, err := DeterminantRat(m); err == nil {
		t.Fatal("expected error for non-square matrix")
	}
}

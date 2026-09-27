package main

import (
	"errors"
	"fmt"
	"math/big"
)

// ValidateSquareRat returns an error unless m is a non-empty square matrix.
func ValidateSquareRat(m [][]*big.Rat) error {
	n := len(m)
	if n == 0 {
		return errors.New("matrix must not be empty")
	}
	for i, row := range m {
		if len(row) != n {
			return fmt.Errorf("row %d has length %d, expected %d", i, len(row), n)
		}
		for j, v := range row {
			if v == nil {
				return fmt.Errorf("nil entry at [%d][%d]", i, j)
			}
		}
	}
	return nil
}

// DeterminantRat computes det(a) exactly using rational arithmetic. if too slow could change to Bareiss algorithm (fraction-free Gaussian elimination)
// Gaussian elimination with partial pivoting. No floating-point rounding.
func DeterminantRat(a [][]*big.Rat) (*big.Rat, error) {
	if err := ValidateSquareRat(a); err != nil {
		return nil, err
	}

	n := len(a)
	m := make([][]*big.Rat, n)
	for i := range a {
		m[i] = make([]*big.Rat, n)
		for j := range a[i] {
			m[i][j] = new(big.Rat).Set(a[i][j])
		}
	}

	det := big.NewRat(1, 1)
	zero := new(big.Rat)

	for col := 0; col < n; col++ {
		// Find a non-zero pivot in this column.
		pivot := -1
		for r := col; r < n; r++ {
			if m[r][col].Cmp(zero) != 0 {
				pivot = r
				break
			}
		}
		if pivot == -1 {
			return new(big.Rat), nil // singular
		}
		if pivot != col {
			m[pivot], m[col] = m[col], m[pivot]
			det.Neg(det)
		}

		det.Mul(det, m[col][col])

		for r := col + 1; r < n; r++ {
			factor := new(big.Rat).Quo(m[r][col], m[col][col])
			for c := col; c < n; c++ {
				m[r][c].Sub(m[r][c], new(big.Rat).Mul(factor, m[col][c]))
			}
		}
	}

	return det, nil
}

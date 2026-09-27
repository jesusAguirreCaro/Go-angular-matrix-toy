package main

import (
	"errors"
	"fmt"
	"math"
)

const epsilon = 1e-12

// ValidateSquare returns an error unless m is a non-empty square matrix
// containing only finite values.
func ValidateSquare(m [][]float64) error {
	n := len(m)
	if n == 0 {
		return errors.New("matrix must not be empty")
	}
	for i, row := range m {
		if len(row) != n {
			return fmt.Errorf("row %d has length %d, expected %d", i, len(row), n)
		}
		for _, v := range row {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return errors.New("matrix contains non-finite values")
			}
		}
	}
	return nil
}

// Determinant computes det(a) via Gaussian elimination with partial
// pivoting. O(n^3). Returns 0 for singular matrices.
func Determinant(a [][]float64) (float64, error) {
	if err := ValidateSquare(a); err != nil {
		return 0, err
	}

	n := len(a)
	m := make([][]float64, n)
	for i := range a {
		m[i] = make([]float64, n)
		copy(m[i], a[i])
	}

	det := 1.0
	for col := 0; col < n; col++ {
		// Partial pivoting: pick the largest magnitude entry in this column.
		pivot := col
		for r := col + 1; r < n; r++ {
			if math.Abs(m[r][col]) > math.Abs(m[pivot][col]) {
				pivot = r
			}
		}
		if math.Abs(m[pivot][col]) < epsilon {
			return 0, nil // singular
		}
		if pivot != col {
			m[pivot], m[col] = m[col], m[pivot]
			det = -det
		}

		det *= m[col][col]
		for r := col + 1; r < n; r++ {
			factor := m[r][col] / m[col][col]
			if factor == 0 {
				continue
			}
			for c := col; c < n; c++ {
				m[r][c] -= factor * m[col][c]
			}
		}
	}
	return det, nil
}

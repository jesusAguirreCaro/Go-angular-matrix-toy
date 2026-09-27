package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"
)

type calcRequest struct {
	Matrix [][]json.Number `json:"matrix"`
}

type determinantResponse struct {
	Determinant string `json:"determinant"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handleHealth)
	mux.HandleFunc("POST /api/determinant", handleDeterminant)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           withCORS(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("listening on http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}

func trimTrailingZeros(s string) string {
	if !strings.Contains(s, ".") {
		return s
	}
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func handleDeterminant(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.UseNumber() // preserve the exact text of each number

	var req calcRequest
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	matrix, err := parseMatrix(req.Matrix)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	det, err := DeterminantRat(matrix)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// RatString() returns "306", "-2", or "3/2" — always exact. FloatString(5) returns trimmed decimals
	writeJSON(w, http.StatusOK, determinantResponse{Determinant: trimTrailingZeros(det.FloatString(5))})
}

// parseMatrix converts the raw JSON numbers into exact rationals.
// big.Rat.SetString accepts "5", "-3", "2.5", and even "1/2".
func parseMatrix(rows [][]json.Number) ([][]*big.Rat, error) {
	if len(rows) == 0 {
		return nil, errors.New("matrix must not be empty")
	}
	out := make([][]*big.Rat, len(rows))
	for i, row := range rows {
		out[i] = make([]*big.Rat, len(row))
		for j, num := range row {
			rat, ok := new(big.Rat).SetString(num.String())
			if !ok {
				return nil, fmt.Errorf("invalid number at row %d, col %d: %q", i, j, num.String())
			}
			out[i][j] = rat
		}
	}
	return out, nil
}

// --- helpers ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

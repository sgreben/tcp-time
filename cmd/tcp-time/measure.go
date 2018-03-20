package main

import (
	"sort"
	"time"
)

type sample struct {
	Success  bool
	Duration time.Duration
}

type measurements []sample

func (m *measurements) append(s sample) {
	*m = append(*m, s)
}

func (m *measurements) insuccessCount() (out int) {
	for _, s := range *m {
		if !s.Success {
			out++
		}
	}
	return
}
func (m *measurements) allSeconds() []float64 {
	out := make([]float64, len(*m))
	for i, s := range *m {
		out[i] = float64(s.Duration)
	}
	sort.Float64s(out)
	return out
}

func (m *measurements) successSeconds() []float64 {
	out := make([]float64, 0, len(*m))
	for _, s := range *m {
		if s.Success {
			out = append(out, float64(s.Duration))
		}
	}
	sort.Float64s(out)
	return out
}

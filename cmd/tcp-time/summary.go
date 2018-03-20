package main

import (
	"log"
	"time"

	"github.com/gonum/floats"
	"github.com/gonum/stat"
)

type bucket struct {
	Label string
	Value time.Duration
	Count int
}

type distribution struct {
	Mean      float64
	StdDev    float64
	Quantiles []float64
	Histogram []bucket
}

type summary struct {
	All   *distribution
	Valid *distribution `json:",omitempty"`
}

func makeDistribution(x []float64) (out distribution) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	if len(x) == 0 {
		return
	}
	out.Mean = stat.Mean(x, nil)
	out.StdDev = stat.StdDev(x, nil)
	quantiles := []float64{0, 0.25, 0.5, 0.75, 1.0}
	for _, q := range quantiles {
		out.Quantiles = append(out.Quantiles, stat.Quantile(q, 1, x, nil))
	}
	dividers := make([]float64, config.HistogramBins+1)
	min := floats.Min(x)
	max := floats.Max(x)
	floats.Span(dividers, min, max+1)
	hist := stat.Histogram(nil, dividers, x, nil)
	for i := range hist {
		out.Histogram = append(out.Histogram, bucket{
			Label: time.Duration(dividers[i]).String(),
			Value: time.Duration(dividers[i]),
			Count: int(hist[i]),
		})
	}
	return
}

func (m *measurements) summary() (out summary) {
	all := makeDistribution(m.allSeconds())
	out.All = &all

	if m.invalidCount() > 0 {
		valid := makeDistribution(m.validSeconds())
		out.Valid = &valid
	}
	return
}

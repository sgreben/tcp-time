package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/schollz/progressbar"
	"golang.org/x/sync/semaphore"
)

type configuration struct {
	Target        string
	N             int
	Parallel      int
	HistogramBins int
	Debug         bool
	Progress      bool
}

type output struct {
	Measurements measurements
	Summary      summary
}

var config configuration

func init() {
	flag.StringVar(&config.Target, "target", "duckduckgo.com:443", "host:port to ping.")
	flag.IntVar(&config.N, "n", 10, "Number of pings to make.")
	flag.IntVar(&config.Parallel, "p", 3, "Number of pings to make in parallel.")
	flag.IntVar(&config.HistogramBins, "b", 5, "Number of histogram bins.")
	flag.BoolVar(&config.Debug, "debug", false, "Print debug logs to stderr.")
	flag.BoolVar(&config.Progress, "progress", false, "Print a progress bar to stderr.")
	flag.Parse()
	log.SetOutput(os.Stderr)
	if !config.Debug {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {
	var output output
	var mu sync.Mutex
	var wg sync.WaitGroup
	ctx := context.Background()
	limit := semaphore.NewWeighted(int64(config.Parallel))
	var bar *progressbar.ProgressBar
	if config.Progress {
		bar = progressbar.New(config.N)
		bar.SetWriter(os.Stderr)
	}

	for i := 0; i < config.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if config.Progress {
				defer bar.Add(1)
			}
			limit.Acquire(ctx, 1)
			defer limit.Release(1)
			d, err := connectDuration(config.Target)
			mu.Lock()
			defer mu.Unlock()
			output.Measurements.append(sample{
				Valid:    err == nil,
				Duration: d,
			})
		}()
	}
	wg.Wait()
	output.Summary = output.Measurements.summary()
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(output)
}

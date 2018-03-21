package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/schollz/progressbar"
	"golang.org/x/time/rate"
)

type configuration struct {
	Target        string
	N             int
	Parallel      int
	HistogramBins int
	RateLimit     float64
	Debug         bool
	Progress      bool
	CSV           bool
	CSVFile       string
}

type output struct {
	Measurements  measurements
	Summary       summary
	Configuration configuration
}

var config configuration
var out output
var csvFile *os.File

func init() {
	flag.StringVar(&config.Target, "target", "duckduckgo.com:443", "host:port to connect to.")
	flag.IntVar(&config.N, "n", 10, "Number of connections to make.")
	flag.Float64Var(&config.RateLimit, "rate-limit", 0.0, "Rate limit (connections per second) to apply.")
	flag.IntVar(&config.Parallel, "p", 3, "Number of connections to make in parallel.")
	flag.IntVar(&config.HistogramBins, "b", 5, "Number of histogram bins.")
	flag.BoolVar(&config.Debug, "debug", false, "Print debug logs to stderr.")
	flag.BoolVar(&config.Progress, "progress", false, "Print a progress bar to stderr.")
	flag.BoolVar(&config.CSV, "csv", false, "Print CSV (index,success,duration) instead of JSON")
	flag.StringVar(&config.CSVFile, "csv-file", "", "Write CSV (index,success,duration) to a file.")
	flag.Parse()
	log.SetOutput(os.Stderr)
	if !config.Debug {
		log.SetOutput(ioutil.Discard)
	}
}

func worker(wg *sync.WaitGroup, mu *sync.Mutex, work chan int) {
	wg.Add(1)
	defer wg.Done()
	for i := range work {
		d, err := connectDuration(config.Target)
		sample := sample{
			Success:  err == nil,
			Duration: d,
		}
		success := 0
		if sample.Success {
			success = 1
		}
		if !config.CSV {
			mu.Lock()
			out.Measurements.append(sample)
			mu.Unlock()
		} else {
			fmt.Printf("%d,%d,%d", i, success, sample.Duration)
			fmt.Println()
		}
		if csvFile != nil {
			fmt.Fprintf(csvFile, "%d,%d,%d", i, success, sample.Duration)
			fmt.Fprintln(csvFile)
		}
	}
}

func main() {
	var bar *progressbar.ProgressBar
	if config.Progress {
		bar = progressbar.New(config.N)
		bar.SetWriter(os.Stderr)
	}
	if config.CSVFile != "" {
		var err error
		csvFile, err = os.OpenFile(config.CSVFile, os.O_CREATE|os.O_RDWR, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	ctx := context.Background()
	var limiter *rate.Limiter

	if config.RateLimit > 0 {
		limiter = rate.NewLimiter(rate.Limit(config.RateLimit), 1)
	}

	work := make(chan int, config.Parallel)

	for i := 0; i < config.Parallel; i++ {
		go worker(&wg, &mu, work)
	}

	for i := 0; i < config.N; i++ {
		if limiter != nil {
			err := limiter.Wait(ctx)
			if err != nil {
				log.Println(err)
			}
		}
		work <- i
		if config.Progress {
			bar.Add(1)
		}
	}
	close(work)
	wg.Wait()
	if !config.CSV {
		out.Configuration = config
		out.Summary = out.Measurements.summary()
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(out)
	}
}

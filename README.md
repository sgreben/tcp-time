# tcp-time

Repeatedly measures TCP connection durations to a given target. Prints the measurements and a summary (mean, standard deviation, quantiles, histogram) as JSON to stdout.

<!-- TOC -->

- [Get it](#get-it)
- [Use it](#use-it)
- [Example](#example)

<!-- /TOC -->

## Get it

```bash
go get -u github.com/sgreben/tcp-time/cmd/tcp-time
```

Or [download the binary](https://github.com/sgreben/tcp-time/releases/latest) from the releases page. 

```bash
# Linux
curl -LO https://github.com/sgreben/tcp-time/releases/download/1.0.5/tcp-time_1.0.5_linux_x86_64.zip
unzip tcp-time_1.0.5_linux_x86_64.zip

# OS X
curl -LO https://github.com/sgreben/tcp-time/releases/download/1.0.5/tcp-time_1.0.5_osx_x86_64.zip
unzip tcp-time_1.0.5_osx_x86_64.zip

# Windows
curl -LO https://github.com/sgreben/tcp-time/releases/download/1.0.5/tcp-time_1.0.5_windows_x86_64.zip
unzip tcp-time_1.0.5_windows_x86_64.zip
```

## Use it

```text
Usage of tcp-time:
  -target string
    	host:port to ping. (default "duckduckgo.com:443")
  -n int
    	Number of connections to make. (default 10)
  -p int
    	Number of connections to make in parallel. (default 3)
  -b int
    	Number of histogram bins. (default 5)
  -rate-limit float
    	Rate limit (connections per second) to apply.
  -csv
    	Print CSV (index,success,duration) instead of JSON
  -csv-file string
        Write CSV (index,success,duration) to a file.
  -progress
    	Print a progress bar to stderr.
  -debug
    	Print debug logs to stderr.
```

## Example

```bash
$ tcp-time -target github.com:443 -n 100 | jp -type bar -xy '..[Label,Count]'
                                    47.0
                               ██████████████
                               ██████████████
                               ██████████████
                               ██████████████
                               ██████████████
                               ██████████████
                               ██████████████
                               ██████████████
                               ██████████████
                     23.0      ██████████████      21.0
                ██████████████ ██████████████ ██████████████
                ██████████████ ██████████████ ██████████████
                ██████████████ ██████████████ ██████████████
                ██████████████ ██████████████ ██████████████
                ██████████████ ██████████████ ██████████████
                ██████████████ ██████████████ ██████████████
       4.0      ██████████████ ██████████████ ██████████████       5.0
 ██████████████ ██████████████ ██████████████ ██████████████ ██████████████

  118.413316ms   122.026494ms   125.639672ms   129.252851ms   132.866029ms
```

```bash
$ tcp-time -target github.com:443 -n 1000
{
  "Measurements": [
    {
      "Success": true,
      "Duration": 137026172
    },
    {
      "Success": true,
      "Duration": 137579594
    },
    {
      "Success": true,
      "Duration": 138099878
    },
    # ...
  ],
  "Summary": {
    "All": {
      "Mean": 127847627.927,
      "StdDev": 3758632.3347467924,
      "Quantiles": [
        117717262,
        125723549,
        127955884,
        129484843,
        159901839
      ],
      "Histogram": [
        {
          "Label": "117.717262ms",
          "Value": 117717262,
          "Count": 293
        },
        {
          "Label": "126.154177ms",
          "Value": 126154177,
          "Count": 684
        },
        {
          "Label": "134.591093ms",
          "Value": 134591093,
          "Count": 20
        },
        {
          "Label": "143.028008ms",
          "Value": 143028008,
          "Count": 1
        },
        {
          "Label": "151.464924ms",
          "Value": 151464924,
          "Count": 2
        }
      ]
    }
  }
}
```

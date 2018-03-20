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
curl -LO https://github.com/sgreben/tcp-time/releases/download/1.0.0/tcp-time_1.0.0_linux_x86_64.zip
unzip tcp-time_1.0.0_linux_x86_64.zip

# OS X
curl -LO https://github.com/sgreben/tcp-time/releases/download/1.0.0/tcp-time_1.0.0_osx_x86_64.zip
unzip tcp-time_1.0.0_osx_x86_64.zip

# Windows
curl -LO https://github.com/sgreben/tcp-time/releases/download/1.0.0/tcp-time_1.0.0_windows_x86_64.zip
unzip tcp-time_1.0.0_windows_x86_64.zip
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
  -progress
    	Print a progress bar to stderr.
  -debug
    	Print debug logs to stderr.
```

## Example

```bash
$ tcp-time -n 1000 | jq .
{
  "Measurements": [
      {
      "Valid": true,
      "Duration": 40605393
    },
    {
      "Valid": true,
      "Duration": 35469927
    },
    # ...
  ],
  "Summary": {
    "All": {
      "Mean": 36694356.154,
      "StdDev": 31847099.84842516,
      "Quantiles": [ 32873052, 34962729, 36098973, 37387095, 89838530 ],
      "Histogram": [
        {
          "Label": "32.873052ms",
          "Value": 32873052,
          "Count": 990
        },
        {
          "Label": "47.114421ms",
          "Value": 47114421,
          "Count": 2
        },
        {
          "Label": "61.355791ms",
          "Value": 61355791,
          "Count": 2
        },
        {
          "Label": "75.597161ms",
          "Value": 75597161,
          "Count": 6
        }
      ]
    }
  }
}
```

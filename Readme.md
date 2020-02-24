# Loxone Statistics

go-loxone-statistics is a tiny library to fetch the statistic values from a Loxone Miniserver.

## Installation

```bash
$ go get github.com/denkfabrik-neueMedien/go-loxone-stats
```

## Usage

```go
package main

import (
    "log"
    lxstats "github.com/denkfabrik-neueMedien/go-loxone-stats"
)

func main() {
    //
    ms := NewMiniserver("host", "user", "pass")
    
    // fetch the list of available statistics
    err := ms.FetchStatistics()
    if err != nil {
        log.Fatal(err)
    }

    //
    uuid := "0d01a765-026e-085a-ffff6f4bfad385ea"
    month := 12
    year := 2015

    //
    s, err := ms.GetStatistic(uuid, month, year)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Found '%s' from %d/%d.", s.Statistics.Name, s.Month, s.Year)
}
```
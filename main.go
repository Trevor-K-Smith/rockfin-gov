//go:generate go build -o rockfin-gov

package main

import (
	"flag"
	"fmt"
	"os"
	"rockfin-gov/internal/aggregator/federal"
)

func main() {
	sourcesCmd := flag.NewFlagSet("sources", flag.ExitOnError)
	samgovCmd := flag.NewFlagSet("samgov", flag.ExitOnError)

	pullCmd := flag.NewFlagSet("pull", flag.ExitOnError)
	pullRawCmd := flag.NewFlagSet("pull-raw", flag.ExitOnError)

	var limit int
	var limitPullRaw int
	pullCmd.IntVar(&limit, "limit", 0, "Limit the number of requests")
	pullRawCmd.IntVar(&limitPullRaw, "limit", 5, "Limit the number of requests")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [command] [subcommand] [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  rockfin-gov sources samgov pull [--limit=1000]\n")
		fmt.Fprintf(os.Stderr, "  rockfin-gov sources samgov pull-raw [--limit=5]\n")
		flag.PrintDefaults()
	}

	if len(os.Args) < 3 {
		flag.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "sources":
		switch os.Args[2] {
		case "samgov":
			switch os.Args[3] {
			case "pull":
				pullCmd.Parse(os.Args[4:])
				samgovCmd.Parse(os.Args[3:4])
				sourcesCmd.Parse(os.Args[2:3])
				federal.CallSamGovAPI(limit)
			case "pull-raw":
				pullRawCmd.Parse(os.Args[4:])
				samgovCmd.Parse(os.Args[3:4])
				sourcesCmd.Parse(os.Args[2:3])
				rawJSON := federal.CallSamGovAPIWithLimit(limitPullRaw)
				fmt.Println(rawJSON)
			default:
				flag.Usage()
				os.Exit(1)
			}
		default:
			flag.Usage()
			os.Exit(1)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}

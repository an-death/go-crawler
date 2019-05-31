package main

import (
	"flag"
)

func parseAgrs() (string, int) {
	var rps int
	var startUrlStr string
	flag.IntVar(&rps, "rps", -1, "requests per second limit")
	flag.StringVar(&startUrlStr, "url", "", "define url for crawler")
	flag.Parse()
	if startUrlStr == "" {
		panic("No host available")
	}

	return startUrlStr, rps
}

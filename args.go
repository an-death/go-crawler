package main

import (
	"flag"
)

func parseAgrs() (string, uint64) {
	var startUrlStr string
	var rps uint64
	flag.Uint64Var(&rps, "rps", 10, "requests per second limit")
	flag.StringVar(&startUrlStr, "url", "", "define url for crawler")
	flag.Parse()
	if startUrlStr == "" {
		panic("No host available")
	}
	return startUrlStr, rps
}

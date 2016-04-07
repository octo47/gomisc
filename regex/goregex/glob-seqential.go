package main

import (
	"github.com/octo47/gomisc/regex"
	"regexp"
)

var pattern = regexp.MustCompile("^[^#]*#[15][15]1110+$")

func main() {
	regex.RunBenchmark(func(line string) bool {
		return pattern.MatchString(line)
	}, "goregex")
}
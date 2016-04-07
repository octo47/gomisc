package main

import (
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"github.com/octo47/gomisc/regex"
)

//var regex = regexp.MustCompile(".*#[15][15]1110+$")
var pattern = pcre.MustCompile("^[^#]*#[15][15]1110+$", 0)


func main() {
	regex.RunBenchmark(func(line string) bool {
		return pattern.MatcherString(line, 0).Matches()
	}, "pcre")
}
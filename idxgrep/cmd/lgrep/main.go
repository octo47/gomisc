package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	_ "github.com/blevesearch/bleve/index/firestorm"
	_ "github.com/blevesearch/bleve/index/store/goleveldb"
	"github.com/golang/glog"
	"github.com/octo47/gomisc/idxgrep"
)

var (
	indexPath  = flag.String("index", "", "index path")
	cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func main() {
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	// open a new index
	if *indexPath == "" {
		log.Fatal("must specify index path")
	}

	if flag.NArg() == 0 {
		log.Fatal("expected at least one file to index")
	}

	indexer := idxgrep.NewBleveIndexer(*indexPath)

	defer func() {
		cerr := indexer.Close()
		if cerr != nil {
			glog.Fatalf("error closing index: %v", cerr)
		}
	}()

	idxgrep.IndexPath(indexer, flag.Args())
}

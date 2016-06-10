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
	engine     = flag.String("engine", "bleve", "Indexer engine to use [bleve]")
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

	var indexer idxgrep.Indexer
	switch *engine {
	case "bleve":
		indexer = idxgrep.NewBleveIndexer(*indexPath)
	default:
		log.Fatalln("Unknown indexing enging", *engine)
	}

	defer func() {
		cerr := indexer.Close()
		if cerr != nil {
			glog.Fatalf("error closing index: %v", cerr)
		}
	}()

	idxgrep.IndexPath(indexer, flag.Args())
}

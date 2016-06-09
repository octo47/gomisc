package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/golang/glog"
)

type Data struct {
	Host    string
	File    string
	Line    string
	Content string
}

var (
	indexPath  = flag.String("index", "", "index path")
	cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")
	batchSize  = flag.Int("batch", 10000, "Lines per indexing batch")
)

func openIndex(idxPath string, overwrite bool) (bleve.Index, error) {
	fileInfo, err := os.Stat(idxPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("can't access index path %v", err)
	}
	if fileInfo != nil {
		if overwrite {
			glog.Info("Overwriting index at ", idxPath)
			os.RemoveAll(idxPath)
		} else {
			glog.Info("Reusing index at ", idxPath)
			index, err := bleve.Open(idxPath)
			if err != nil {
				return nil, fmt.Errorf("Unable to open index at %s: %v", idxPath, err)
			}
			return index, nil
		}
	}
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(*indexPath, mapping)
	if err != nil {
		return nil, fmt.Errorf("Unable to create index at %s: %v", idxPath, err)
	}
	return index, nil
}

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

	glog.Info("Opening index at ", *indexPath)
	index, err := openIndex(*indexPath, true)
	if err != nil {
		glog.Fatal(err)
	}
	defer func() {
		cerr := index.Close()
		if cerr != nil {
			glog.Fatalf("error closing index: %v", err)
		}
	}()
	batch := index.NewBatch()
	for f := range handleArgs(flag.Args()) {
		b := bytes.NewBufferString(f.filename)
		b.WriteByte(':')
		b.WriteString(strconv.FormatInt(f.linenum, 10))
		batch.Index(b.String(), f)
		if batch.Size() > *batchSize {
			index.Batch(batch)
			batch.Reset()
		}
	}
	if batch.Size() > 0 {
		index.Batch(batch)
	}
}

type file struct {
	filename string
	linenum  int64
	contents []byte
}

func handleArgs(files []string) chan file {
	ch := make(chan file)
	go collectFiles(files, ch)
	return ch
}

func collectFiles(files []string, ch chan file) {
	for _, fname := range files {
		fname = filepath.Clean(fname)
		err := filepath.Walk(fname, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				glog.Warningln(err)
				return err
			}
			if info.IsDir() {
				return nil
			}
			glog.Infoln("Indexing file ", path)
			f, err := os.Open(path)
			if err != nil {
				glog.Warningf("unable to read file %s %v", path, err)
				return nil
			}
			defer f.Close()

			var linenum int64 = 0
			scanner := bufio.NewScanner(bufio.NewReaderSize(f, 1<<24))

			for scanner.Scan() {
				linenum++
				bytes := scanner.Bytes()
				if len(bytes) == 0 {
					continue
				}
				ch <- file{
					filename: filepath.Base(path),
					linenum:  linenum,
					contents: bytes,
				}
			}
			return nil
		})
		if err != nil {
			glog.Fatal(err)
		}
	}
	close(ch)
}

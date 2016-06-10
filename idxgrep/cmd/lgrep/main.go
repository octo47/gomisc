package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strconv"

	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/index/firestorm"
	_ "github.com/blevesearch/bleve/index/store/goleveldb"
	"github.com/golang/glog"
)

const (
	chunkSize   = 1 << 16
	channelSize = 4
)

var (
	indexPath  = flag.String("index", "", "index path")
	cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")
)

// Structure used for indexing.
// Keep in sync with indexMapping() function
type Data struct {
	HostName string "json:hostname"
	FileName string "json:filename"
	Offset   int64  "json:offset"
	Contents []byte "json:contents"
}

func (d *Data) indexKey() string {
	var b bytes.Buffer
	b.WriteString(d.HostName)
	b.WriteByte(1)
	b.WriteString(d.FileName)
	b.WriteByte(1)
	b.WriteString(strconv.FormatInt(d.Offset, 10))
	return b.String()
}

func (d *Data) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("data[")
	buffer.WriteString(" hostname:")
	buffer.WriteString(d.HostName)
	buffer.WriteString(" filename:")
	buffer.WriteString(d.FileName)
	buffer.WriteString(" offset:")
	buffer.WriteString(strconv.FormatInt(d.Offset, 10))
	buffer.WriteString(" len(contentes):")
	buffer.WriteString(strconv.Itoa(len(d.Contents)))
	buffer.WriteString("]")
	return buffer.String()
}

func buildIndexMapping() *bleve.IndexMapping {
	mapping := bleve.NewDocumentStaticMapping()
	textField := bleve.NewTextFieldMapping()
	textField.IncludeTermVectors = false
	mapping.AddFieldMappingsAt("hostname", textField)
	mapping.AddFieldMappingsAt("filename", textField)
	mapping.AddFieldMappingsAt("offset", bleve.NewNumericFieldMapping())
	contentsField := bleve.NewTextFieldMapping()
	contentsField.Store = false
	contentsField.IncludeTermVectors = false
	contentsField.IncludeInAll = false
	mapping.AddFieldMappingsAt("contents", contentsField)
	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("file", mapping)
	return indexMapping
}

func openIndex(path string) bleve.Index {
	glog.Info("Opening index ", path)
	index, err := bleve.Open(path)
	if err == bleve.ErrorIndexPathDoesNotExist {
		glog.Infof("Creating new index...")
		// create a mapping
		indexMapping := buildIndexMapping()
		index, err = bleve.NewUsing(path, indexMapping, "upside_down", "goleveldb", nil)
		if err != nil {
			glog.Fatal(err)
		}
	} else if err == nil {
		glog.Info("Opening existing index...")
	} else {
		glog.Fatal(err)
	}
	return index
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

	index := openIndex(*indexPath)
	defer func() {
		cerr := index.Close()
		if cerr != nil {
			glog.Fatalf("error closing index: %v", cerr)
		}
	}()
	for f := range handleArgs(flag.Args()) {
		if glog.V(3) {
			glog.Info("indexing ", f.String())
		}
		index.Index(f.indexKey(), f)
		if glog.V(3) {
			glog.Info("done indexing ", f.indexKey())
		}
	}
}

func handleArgs(files []string) chan Data {
	ch := make(chan Data, channelSize)
	go collectFiles(files, ch)
	return ch
}

func collectFiles(files []string, ch chan Data) {
	var hostname, err = os.Hostname()
	if err != nil {
		glog.Fatal("unable to get hostname", err)
	}
	for _, fname := range files {
		fname = filepath.Clean(fname)
		fname, _ = filepath.Abs(fname)
		err := filepath.Walk(fname, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				glog.Warningln(err)
				return err
			}
			if info.IsDir() {
				return nil
			}
			glog.Infoln("Reading file ", path)
			f, err := os.Open(path)
			if err != nil {
				glog.Warningf("unable to read file %s %v", path, err)
				return nil
			}
			defer f.Close()

			reader := bufio.NewReaderSize(f, chunkSize)
			offset := int64(0)

			var buffer []byte = make([]byte, chunkSize, chunkSize)
			for {
				size, err := reader.Read(buffer)
				if err != nil {
					if err == io.EOF {
						break
					}
					glog.Fatal("Error during reading file", path, err)
				}
				// read to the end of the line
				// we want to ensure that we not truncating tokens
				line, lerr := reader.ReadBytes('\n')
				if lerr != nil && lerr != io.EOF {
					glog.Fatal("Error during reading file", path, lerr)
				}
				// merge and send data to indexer
				chunkData := make([]byte, size+len(line))
				copy(chunkData, buffer[:size])
				chunkData = append(chunkData, line...)
				data := Data{
					HostName: hostname,
					FileName: path,
					Offset:   offset,
					Contents: chunkData,
				}
				if glog.V(3) {
					glog.Info(data.String())
				}
				ch <- data
				// advance offset for next chunk
				offset += int64(size) + int64(len(line))
			}
			return nil
		})
		if err != nil {
			glog.Fatal(err)
		}
		glog.Infoln("Walk complete for ", fname)
	}
	close(ch)
}

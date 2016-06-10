package idxgrep

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/blevesearch/bleve/index/firestorm"
	_ "github.com/blevesearch/bleve/index/store/goleveldb"
	"github.com/golang/glog"
)

type Indexer interface {
	Index(document *Data)
}

const (
	chunkSize   = 1 << 16
	channelSize = 4
)

// Structure used for indexing.
type Data struct {
	HostName string "json:hostname"
	FileName string "json:filename"
	Offset   int64  "json:offset"
	Contents []byte "json:contents"
}

func (d *Data) IndexKey() string {
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

func IndexPath(i Indexer, pathsToIndex []string) {
	for f := range handleArgs(pathsToIndex) {
		if glog.V(3) {
			glog.Info("indexing ", f.String())
		}
		i.Index(&f)
		if glog.V(3) {
			glog.Info("done indexing ", f.IndexKey())
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

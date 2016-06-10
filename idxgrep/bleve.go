package idxgrep

import (
	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/index/firestorm"
	_ "github.com/blevesearch/blevex/cznicb"
	_ "github.com/blevesearch/blevex/leveldb"
	_ "github.com/blevesearch/blevex/rocksdb"
	"github.com/golang/glog"
)

type BleveIndexer struct {
	index bleve.Index
}

func NewBleveIndexer(path string, store string) *BleveIndexer {
	return &BleveIndexer{
		index: openIndex(path, store),
	}
}

func (i *BleveIndexer) Index(d *Data) {
	key := d.IndexKey()
	i.index.Index(key, d)
}

func (i *BleveIndexer) Close() error {
	return i.index.Close()
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

func openIndex(path string, engine string) bleve.Index {
	glog.Info("Opening index ", path)
	index, err := bleve.Open(path)
	if err == bleve.ErrorIndexPathDoesNotExist {
		glog.Infof("Creating new index...")
		// create a mapping
		indexMapping := buildIndexMapping()
		index, err = bleve.NewUsing(path, indexMapping, "upside_down", engine, nil)
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

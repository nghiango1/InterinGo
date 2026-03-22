package store

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

type EditedFile struct {
	data *protocol.TextDocumentItem
}

func (ef *EditedFile) Unwrap() *protocol.TextDocumentItem {
	return ef.data
}

func Wrap(tdi *protocol.TextDocumentItem) *EditedFile {
	return &EditedFile {
		data: tdi,
	}
}

func (ef *EditedFile) Update(changeData protocol.TextDocumentContentChangeEvent ) {
	start, end := changeData.Range.IndexesIn(ef.data.Text)
	ef.data.Text = ef.data.Text[:start] + changeData.Text + ef.data.Text[end:]
}

func (ef *EditedFile) UpdateWhole(changeData protocol.TextDocumentContentChangeEventWhole ) {
	ef.data.Text = changeData.Text
}

type ServerFileSyncStore struct {
	clientOpenedFile map[protocol.URI]*EditedFile
}

var store ServerFileSyncStore;

func Init() {
	store.clientOpenedFile  = make(map[protocol.URI]*EditedFile)
}

func GetStore() *ServerFileSyncStore {
	if store.clientOpenedFile == nil {
		Init()
	}
	return &store
}

func (store *ServerFileSyncStore) Add(ef *EditedFile){
	store.clientOpenedFile[ef.Unwrap().URI] = ef
}

func (store *ServerFileSyncStore) Get(uri protocol.URI) (*EditedFile, error){
	ef, ok := store.clientOpenedFile[uri]
	if !ok {
		return nil, errors.New("File isn't registed as document open file by client")
	}
	return ef, nil
}

func (store *ServerFileSyncStore) populateURIFileData(uri string) error {
	textDocItem := &protocol.TextDocumentItem {
		URI: uri,
	}
	ef := Wrap(textDocItem)
	store.Add(ef)

	ef.Unwrap().Version = 0
	ef.Unwrap().LanguageID = ""

	filepath, err := url.Parse(uri)
	if err != nil {
		return errors.New(fmt.Sprint("Can't parse file uri provided", uri))
	}

	data, err := os.ReadFile(filepath.Path)
	if err != nil {
		return errors.New(fmt.Sprint("Can't read file at", filepath))
	}

	ef.Unwrap().Text = string(data)

	return nil
}


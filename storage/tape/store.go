package tape

import (
	"encoding/gob"
	"fmt"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/keys"
	"github.com/louisevanderlith/husk/persisted"
	"github.com/louisevanderlith/husk/records"
	"github.com/louisevanderlith/husk/storage"
	"github.com/louisevanderlith/husk/validation"
	"io"
	"os"
	"reflect"
)

func init() {
	gob.Register(&keys.TimeKey{})
}

//NewDefaultStore returns Storage using Gob for encoding/decoding
func NewDefaultStore(obj validation.Dataer) hsk.Storage {
	return NewStore(obj, storage.GobEncoder, storage.GobDecoder)
}

//NewStore returns Storage using the specified encoding/decoding
func NewStore(obj validation.Dataer, encfunc storage.NewEncoder, decfunc storage.NewDecoder) hsk.Storage {
	t := reflect.TypeOf(obj)

	if t.Kind() == reflect.Ptr {
		panic("obj must not be a pointer")
	}

	err := persisted.CreateDirectory("db")

	if err != nil {
		panic(err)
	}

	gob.Register(obj)

	trackName := fmt.Sprintf("db/%s.Data.husk", t.Name())

	track, err := os.OpenFile(trackName, os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		panic(err)
	}

	return &tapeStore{
		t:     t,
		enc:   encfunc,
		dec:   decfunc,
		track: track,
	}
}

//tapeStore 101[0]00100011
type tapeStore struct {
	t     reflect.Type
	enc   storage.NewEncoder
	dec   storage.NewDecoder
	track *os.File
}

func (ts *tapeStore) Name() string {
	return ts.t.Name()
}

func (ts *tapeStore) ZeroValue() validation.Dataer {
	return reflect.Zero(ts.t).Interface().(validation.Dataer)
}

func (ts *tapeStore) Read(p hsk.Point, res chan<- hsk.Record) {
	rec := records.NewRecord(ts.ZeroValue())
	r := io.NewSectionReader(ts.track, p.GetOffset(), p.GetLength())

	serial := ts.dec(r)
	err := serial.Decode(rec)

	if err != nil {
		panic(err)
	}

	res <- rec
}

//Write will append obj to the end of the file
func (ts *tapeStore) Write(obj hsk.Record, p chan<- hsk.Point) {
	nf, err := ts.track.Seek(0, io.SeekEnd)

	if err != nil {
		panic(err)
	}

	serial := ts.enc(ts.track)
	err = serial.Encode(obj)

	if err != nil {
		panic(err)
	}

	endWith, err := ts.GetSize()

	if err != nil {
		panic(err)
	}

	p <- hsk.NewPoint(nf, endWith-nf)
}

func (ts *tapeStore) GetSize() (int64, error) {
	inf, err := ts.track.Stat()

	if err != nil {
		return 0, err
	}

	return inf.Size(), nil
}

//Close closes the Data Track
func (ts *tapeStore) Close() error {
	//we can expand here, later
	return ts.track.Close()
}

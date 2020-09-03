package tape

import (
	"fmt"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/storers"
	"github.com/louisevanderlith/husk/validation"
	"io"
	"os"
	"reflect"
)

func NewStore(t reflect.Type, encfunc storers.NewEncoder, decfunc storers.NewDecoder) storers.Storage {
	trackName := fmt.Sprintf("db/%s.Data.husk", t.Name())

	track, err := os.OpenFile(trackName, os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		panic(err)
	}

	return &tapeStore{
		enc:   encfunc,
		dec:   decfunc,
		t:     t,
		track: track,
	}
}

//tapeStore 101[0]00100011
type tapeStore struct {
	enc   storers.NewEncoder
	dec   storers.NewDecoder
	t     reflect.Type
	track *os.File
}

func (ts *tapeStore) Read(p hsk.Point, res chan<- validation.Dataer) {
	dObj := reflect.New(ts.t)
	dInf := dObj.Interface()

	r := io.NewSectionReader(ts.track, p.GetOffset(), p.GetLength())

	serial := ts.dec(r)
	err := serial.Decode(dInf)

	if err != nil {
		panic(err)
	}

	res <- dObj.Elem().Interface().(validation.Dataer)
}

//Write will append obj to the end of the file
func (ts *tapeStore) Write(obj validation.Dataer, p chan<- hsk.Point) {
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

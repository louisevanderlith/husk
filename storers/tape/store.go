package tape

import (
	"encoding/gob"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/storers"
	"io"
	"os"
	"reflect"
)

//tapeStore 101[0]00100011
type tapeStore struct {
	t     reflect.Type
	track *os.File
}

func newStore(t reflect.Type) storers.Storer {
	trackName := getDataPath(t.Name())
	track, err := os.OpenFile(trackName, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}

	return &tapeStore{
		t:     t,
		track: track,
	}
}

func (ts *tapeStore) Read(p *hsk.Point, res chan<- hsk.Dataer) error {

	go func() {
		dObj := reflect.New(ts.t)
		dInf := dObj.Interface()

		r := io.NewSectionReader(ts.track, p.Offset, p.Len)

		serial := gob.NewDecoder(r)
		err := serial.Decode(dInf)

		if err != nil {
			panic(err)
		}

		res <- dObj.Elem().Interface().(hsk.Dataer)
	}()

	return nil
}

//Write will append obj to the end of the file
func (ts *tapeStore) Write(obj hsk.Dataer) (*hsk.Point, error) {
	nf, err := ts.track.Seek(0, io.SeekEnd)

	if err != nil {
		return nil, err
	}

	serial := gob.NewEncoder(ts.track)
	err = serial.Encode(obj)

	if err != nil {
		return nil, err
	}

	endWith, err := ts.GetSize()

	if err != nil {
		return nil, err
	}

	return hsk.NewPoint(nf, endWith-nf), nil
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

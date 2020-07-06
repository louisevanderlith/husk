package husk

import (
	"encoding/gob"
	"io"
	"os"
)

//Taper 101[0]00100011
type Taper interface {
	Read(point *Point, obj interface{}) error
	Write(obj Dataer) (*Point, error)
	Close()
}

type tape struct {
	track *os.File
}

func newTape(trackname string) Taper {
	track, err := os.OpenFile(trackname, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}

	return &tape{
		track: track,
	}
}

//Reads the data @point into obj
func (t *tape) Read(point *Point, obj interface{}) error {
	r := io.NewSectionReader(t.track, point.Offset, point.Len)

	serial := gob.NewDecoder(r)
	return serial.Decode(obj)
}

func (t *tape) GetSize() (int64, error) {
	inf, err := t.track.Stat()

	if err != nil {
		return 0, err
	}

	return inf.Size(), nil
}

//Write will append obj to the end of the file
func (t *tape) Write(obj Dataer) (*Point, error) {
	nf, err := t.track.Seek(0, io.SeekEnd)

	if err != nil {
		return nil, err
	}

	serial := gob.NewEncoder(t.track)
	err = serial.Encode(obj)

	if err != nil {
		return nil, err
	}

	endWith, err := t.GetSize()

	if err != nil {
		return nil, err
	}

	return newPoint(nf, endWith-nf), nil
}

//Close closes the Data Track
func (t *tape) Close() {
	//we can expand here, later
	t.track.Close()
}

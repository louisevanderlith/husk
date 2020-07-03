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
	track  *os.File
	offset int64
}

func newTape(trackname string) Taper {
	track, err := os.OpenFile(trackname, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}

	return &tape{track, int64(0)}
}

//Reads the data @point into obj
func (t *tape) Read(point *Point, obj interface{}) error {
	_, err := t.track.Seek(point.Offset, io.SeekStart)

	if err != nil {
		return err
	}

	defer t.track.Seek(0, io.SeekStart)

	serial := gob.NewDecoder(t.track)
	return serial.Decode(obj)
}

func (t *tape) GetSize() (int64, error) {
	inf, err := t.track.Stat()

	if err != nil {
		return 0, err
	}

	return inf.Size(), nil
}

func (t *tape) Write(obj Dataer) (*Point, error) {
	startOff, err := t.GetSize()

	if err != nil {
		return nil, err
	}

	result := newPoint(startOff, 0)

	nf, err := t.track.Seek(result.Offset, 1)

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

	t.offset += nf
	result.Len = endWith - startOff

	t.track.Seek(0, io.SeekStart)

	return result, nil
}

//Close closes the Data Track
func (t *tape) Close() {
	//we can expand here, later
	t.track.Close()
}

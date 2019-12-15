package husk

import (
	"fmt"
	"io"
	"log"
	"os"
)

//Taper 101[0]00100011
type Taper interface {
	Read(point *Point, obj interface{}) error
	Write(obj interface{}) (*Point, error)
	Close()
}

type tape struct {
	track  *os.File
	offset int64
	serial Serializer
}

func newTape(trackname string, serial Serializer) Taper {
	track, err := os.OpenFile(trackname, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}

	//enc := gob.NewEncoder(track)
	//dec := gob.NewDecoder(track)

	return &tape{track, int64(0), serial}
}

//Reads the data @point into obj
func (t *tape) Read(point *Point, obj interface{}) error {
	//len := point.Len
	//byts := make([]byte, len) //, len, len)

	nxtOff, err := t.track.Seek(point.Offset, io.SeekStart)

	if err != nil {
		return err
	}

	defer t.track.Seek(0, io.SeekStart)
	log.Println(nxtOff)
	/*read, err := t.track.ReadAt(byts, point.Offset)

	if err != nil && err != io.EOF {
		return err
	}

	//The database is still empty.
	if int64(read) != len {
		return fmt.Errorf("read %v, need %v", read, len)
	}*/

	//buffer := bytes.NewBuffer(byts)

	return t.serial.Decode(t.track, obj)
}

func (t *tape) Write(obj interface{}) (*Point, error) {
	result := newPoint(t.offset, 0)

	byts, err := t.serial.Encode(obj) //toBytes(obj)

	if err != nil {
		return nil, err
	}

	//wrote, err := t.track.WriteAt(byts, t.offset)
	_, err = t.track.Seek(t.offset, 1)
	if err != nil {
		return nil, err
	}

	wrote, err := t.track.Write(byts)

	if err != nil {
		return nil, err
	}

	if wrote != len(byts) {
		return nil, fmt.Errorf("Incomplete write. %v - %v", wrote, len(byts))
	}

	written := int64(wrote)
	t.offset += written
	result.Len = int64(len(byts))

	t.track.Seek(0, io.SeekStart)
	return result, nil
}

//Close closes the Data Track
func (t *tape) Close() {
	//we can expand here, later
	t.track.Close()
}

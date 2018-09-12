package husk

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
)

type Taper interface {
	Read(point *Point, obj interface{}) error
	Write(obj interface{}) (*Point, error)
	Close()
}

type tape struct {
	track  *os.File
	offset int64
}

func NewTape(trackname string) Taper {
	track, err := os.OpenFile(trackname, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}

	return &tape{track, int64(0)}
}

/*
var err error
	indx := make(map[int]*Point)
	written := int64(0) //total file len (when reading)
	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("%d \n- %d", i, i)
		bitMsg := []byte(msg)
		indx[i] = NewPoint(written, int64(len(bitMsg)))


		if err != nil {
			log.Print(err)
		}

		written += int64(len(bitMsg))
	}

	fmt.Println("DONE")

	fmt.Printf("Item 3: %+v", indx[2])

*/

func (t *tape) Read(point *Point, result interface{}) error {
	len := point.Len
	byts := make([]byte, len, len)

	read, err := t.track.ReadAt(byts, point.Offset)

	if err != nil {
		return err
	}

	if int64(read) != len {
		msg := fmt.Sprintf("Incorrect Read: read %d, len %d", read, len)
		return errors.New(msg)
	}

	buffer := bytes.NewBuffer(byts)
	err = gob.NewDecoder(buffer).Decode(result)

	return nil
}

func (t *tape) Write(obj interface{}) (*Point, error) {
	result := NewPoint(t.offset, 0)

	byts, err := toBytes(obj)

	if err != nil {
		return nil, err
	}

	wrote, err := t.track.WriteAt(byts, t.offset)

	if err != nil {
		return nil, err
	}

	written := int64(wrote)
	t.offset += int64(written)
	result.Len = written

	return result, nil
}

func (t *tape) Close() {
	//we can expand here, later
	t.track.Close()
}

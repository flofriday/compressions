package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func TestPointer(t *testing.T) {
	o1 := 4095
	l1 := 15

	tmp := createPointer(o1, l1)
	o2, l2 := readPointer(tmp)

	if l1 != l2 {
		t.Errorf("Length are not equal. l1=%d l2=%d", l1, l2)
	}

	if o1 != o2 {
		t.Errorf("Offset are not equal. o1=%d o2=%d", o1, o2)
	}
}

func TestBitMask(t *testing.T) {
	d1 := []bool{true, false, false, true, true, true, false, true}

	tmp := createBitMask(d1)
	d2 := readBitMask(tmp)

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("BitMasks are not equal. d1=%v d2=%v", d1, d2)
	}
}

func TestEncodingDecoding(t *testing.T) {
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		panic("No testdata found")
	}

	for _, file := range files {
		path := filepath.Join("testdata", file.Name())
		data, _ := ioutil.ReadFile(path)
		tmp := encode(data)
		data2 := decode(tmp)

		if !bytes.Equal(data, data2) {
			t.Errorf("File %s was not correctly encoded/decoded", file.Name())
		}
	}
}

package main

import "testing"

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

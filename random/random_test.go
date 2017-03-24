package random

import "testing"

func TestNewString(t *testing.T) {
	str, err := NewStr(10)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(str)
}

func TestUpperStr(t *testing.T) {
	str, err := UpperStr(10)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(str)
}

func TestNoUpper(t *testing.T) {
	str, err := NoUpperStr(10)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(str)
}

func TestLower(t *testing.T) {
	str, err := LowerStr(10)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(str)
}

func TestNoLower(t *testing.T) {
	str, err := NoLowerStr(10)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(str)
}

func TestNoNumber(t *testing.T) {
	str, err := NoNumberStr(10)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(str)
}

func TestUUID(t *testing.T) {
	str, err := UUID()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(str)
}

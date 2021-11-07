package boltutil

import (
	"bytes"
	"testing"
	"time"
)

func TestToKeyBytes(t *testing.T) {
	tests := []struct {
		Input    interface{}
		Expected []byte
	}{
		{Input: "aaa", Expected: []byte("aaa")},
		{Input: int64(1), Expected: []byte{0, 0, 0, 0, 0, 0, 0, 1}},
	}

	for _, test := range tests {
		b, err := ToKeyBytes(test.Input)
		if err != nil {
			t.Error(err)
		}

		if bytes.Compare(test.Expected, b) == -1 {
			t.Errorf("expected %v but result: %v", test.Expected, b)
		}
	}
}

func TestToKeyBytesFromTime(t *testing.T) {
	tm1 := time.Date(1980, 1, 27, 10, 59, 59, 0, time.UTC)
	tm2 := time.Date(1980, 1, 27, 10, 59, 59, 1, time.UTC)

	b1, err := ToKeyBytes(tm1)
	if err != nil {
		t.Error(err)
	}

	b2, err := ToKeyBytes(tm2)
	if err != nil {
		t.Error(err)
	}

	t.Log(b1)
	t.Log(b2)
}

type example struct {
	Str string
	Num int
	B   bool
	M   map[string]string
}

func TestSerializeStruct(t *testing.T) {
	e := example{
		Str: "aaa",
		Num: 100,
		B:   true,
		M: map[string]string{
			"hoge": "hoge111",
		},
	}

	b, err := Serialize(e)
	if err != nil {
		t.Error(err)
	}

	//t.Log(b)

	var ret example
	err = Deserialize(b, &ret)
	if err != nil {
		t.Error(err)
	}

	//t.Log(ret)
	if e.Str != ret.Str {
		t.Errorf("expected %v but result: %v", e.Str, ret.Str)
	}
}

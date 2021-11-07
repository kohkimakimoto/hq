package boltutil

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

func ToKeyBytes(in interface{}) ([]byte, error) {
	switch v := in.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case int:
		return numberToBytes(int64(v))
	case uint:
		return numberToBytes(uint64(v))
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		return numberToBytes(v)
	default:
		buf := &bytes.Buffer{}
		enc := gob.NewEncoder(buf)
		err := enc.Encode(v)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

func numberToBytes(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Serialize(in interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(in)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Deserialize(in []byte, out interface{}) error {
	buf := bytes.NewBuffer(in)
	return gob.NewDecoder(buf).Decode(out)
}

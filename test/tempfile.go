package test

import (
	"io/ioutil"
	"os"
)

// CreateTempfile creates temporary file for testing.
// You need to delete the file manually after tests.
func CreateTempfile(b []byte) (*os.File, error) {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}
	if err = ioutil.WriteFile(tempFile.Name(), b, 0644); err != nil {
		return nil, err
	}
	return tempFile, nil
}

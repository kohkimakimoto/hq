package boltutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	bolt "go.etcd.io/bbolt"
)

func testDB(t *testing.T) *bolt.DB {
	db, err := bolt.Open(filepath.Join(testTempDir(t), "test.bolt"), 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func testTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	return dir
}

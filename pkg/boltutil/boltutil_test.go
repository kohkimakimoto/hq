package boltutil

import (
	"bytes"
	"testing"

	bolt "go.etcd.io/bbolt"
)

func TestSet(t *testing.T) {
	db := testDB(t)

	err := db.Update(func(tx *bolt.Tx) error {
		return Set(tx, []interface{}{"foo", "bar"}, "key", "value")
	})
	if err != nil {
		t.Error(err)
	}

	var ret string
	err = db.View(func(tx *bolt.Tx) error {
		return Get(tx, []interface{}{"foo", "bar"}, "key", &ret)
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(ret)

	if ret != "value" {
		t.Errorf("expected 'value' but got: %s", ret)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		return Delete(tx, []interface{}{"foo", "bar"}, "key")
	})
	if err != nil {
		t.Error(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		return Get(tx, []interface{}{"foo", "bar"}, "key", &ret)
	})
	if err == nil || err != ErrNotFound {
		t.Error("should be 'Not found'")
	}
}

func TestCursor(t *testing.T) {
	db := testDB(t)

	err := db.Update(func(tx *bolt.Tx) error {
		if err := Set(tx, []interface{}{"foo", "bar", 10}, 0, "value0"); err != nil {
			return err
		}

		if err := Set(tx, []interface{}{"foo", "bar", 11}, 1, "value1"); err != nil {
			return err
		}

		if err := Set(tx, []interface{}{"foo", "bar", 12}, 2, "value2"); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Error(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		c, err := Cursor(tx, []interface{}{"foo", "bar"})
		if err != nil {
			return err
		}

		i := 0
		expected := []int64{10, 11, 12}
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			t.Logf("key = %d", k)
			b, err := numberToBytes(expected[i])
			if err != nil {
				return err
			}

			if bytes.Compare(b, k) != 0 {
				t.Errorf("expected %v but %v", b, k)
			}

			i++
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteBucket(t *testing.T) {
	db := testDB(t)

	err := db.Update(func(tx *bolt.Tx) error {
		if err := Set(tx, []interface{}{"foo", "bar"}, "key", "value"); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}

	var ret string
	err = db.View(func(tx *bolt.Tx) error {
		return Get(tx, []interface{}{"foo", "bar"}, "key", &ret)
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(ret)

	if ret != "value" {
		t.Errorf("unexpected value %v", ret)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		return DeleteBucket(tx, []interface{}{"foo", "bar"})
	})

	if err != nil {
		t.Error(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		return Get(tx, []interface{}{"foo", "bar"}, "key", &ret)
	})
	if err == nil || err != ErrNotFound {
		t.Error("should be 'Not found'")
	}
}

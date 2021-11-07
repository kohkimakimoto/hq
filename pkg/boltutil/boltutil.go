package boltutil

import (
	"errors"

	bolt "go.etcd.io/bbolt"
)

var (
	ErrNotFound = errors.New("Not Found")
)

// Get gets a value from a bucket. The nested bucket is supported
func Get(tx *bolt.Tx, bucketNames []interface{}, key interface{}, to interface{}) error {
	bucket, err := Bucket(tx, bucketNames)
	if err != nil {
		return err
	}

	if bucket == nil {
		return ErrNotFound
	}

	keyB, err := ToKeyBytes(key)
	if err != nil {
		return err
	}

	v := bucket.Get(keyB)
	if v == nil {
		return ErrNotFound
	}

	err = Deserialize(v, to)
	if err != nil {
		return err
	}

	return nil
}

func Set(tx *bolt.Tx, bucketNames []interface{}, key interface{}, value interface{}) error {
	bucket, err := CreateBucketIfNotExists(tx, bucketNames)
	if err != nil {
		return err
	}

	keyB, err := ToKeyBytes(key)
	if err != nil {
		return err
	}

	valueB, err := Serialize(value)
	if err != nil {
		return err
	}

	err = bucket.Put(keyB, valueB)
	if err != nil {
		return err
	}

	return nil
}

func Delete(tx *bolt.Tx, bucketNames []interface{}, key interface{}) error {
	bucket, err := Bucket(tx, bucketNames)
	if err != nil {
		return err
	}
	if bucket == nil {
		return nil
	}

	keyB, err := ToKeyBytes(key)
	if err != nil {
		return err
	}

	return bucket.Delete(keyB)
}

func DeleteBucket(tx *bolt.Tx, bucketNames []interface{}) error {
	// get and remove last bucket name.
	key := bucketNames[len(bucketNames)-1]
	bucketNames = bucketNames[:len(bucketNames)-1]

	keyB, err := ToKeyBytes(key)
	if err != nil {
		return err
	}

	if len(bucketNames) == 0 {
		return tx.DeleteBucket(keyB)
	} else {
		bucket, err := Bucket(tx, bucketNames)
		if err != nil {
			return err
		}
		if bucket == nil {
			return nil
		}

		return bucket.DeleteBucket(keyB)
	}
}

func Bucket(tx *bolt.Tx, bucketNames []interface{}) (*bolt.Bucket, error) {
	var bucket *bolt.Bucket
	for i, bucketName := range bucketNames {
		b, err := ToKeyBytes(bucketName)
		if err != nil {
			return nil, err
		}

		if i == 0 {
			bucket = tx.Bucket(b)
			if bucket == nil {
				return nil, nil
			}
		} else {
			bucket = bucket.Bucket(b)
			if bucket == nil {
				return nil, nil
			}
		}
	}

	return bucket, nil
}

func CreateBucketIfNotExists(tx *bolt.Tx, bucketNames []interface{}) (*bolt.Bucket, error) {
	var bucket *bolt.Bucket
	for i, bucketName := range bucketNames {
		b, err := ToKeyBytes(bucketName)
		if err != nil {
			return nil, err
		}

		if i == 0 {
			bc, err := tx.CreateBucketIfNotExists(b)
			if err != nil {
				return nil, err
			}
			bucket = bc
		} else {
			bc, err := bucket.CreateBucketIfNotExists(b)
			if err != nil {
				return nil, err
			}
			bucket = bc
		}
	}

	return bucket, nil
}

func Cursor(tx *bolt.Tx, bucketNames []interface{}) (*bolt.Cursor, error) {
	bucket, err := Bucket(tx, bucketNames)
	if err != nil {
		return nil, err
	}
	if bucket == nil {
		return nil, ErrNotFound
	}

	return bucket.Cursor(), nil
}

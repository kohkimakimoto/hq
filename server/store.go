package server

import (
	"encoding/binary"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/kohkimakimoto/boltutil"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/labstack/echo"
	"regexp"
)

type Store struct {
	app    *App
	db     *bolt.DB
	logger echo.Logger
}

const (
	BucketNameForJobs = "j"
)

// Job Error

type ErrJobNotFound struct {
	ID uint64
}

func (e *ErrJobNotFound) Error() string {
	return fmt.Sprintf("Job '%d' is not found", e.ID)
}

type ErrJobAlreadyExisted struct {
	ID   uint64
	Name string
}

func (e *ErrJobAlreadyExisted) Error() string {
	return fmt.Sprintf("'%d' (%s) is already exsited", e.ID, e.Name)
}

func (s *Store) Init() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if _, err := boltutil.CreateBucketIfNotExists(tx, []interface{}{BucketNameForJobs}); err != nil {
			return err
		}
		return nil
	})

	return nil
}

func (s *Store) CreateJob(job *hq.Job) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if err := boltutil.Get(tx, []interface{}{BucketNameForJobs}, job.ID, &hq.J{}); err == nil {
			return &ErrJobAlreadyExisted{ID: job.ID, Name: job.Name}
		}

		in := &hq.J{
			ID:         job.ID,
			Name:       job.Name,
			Comment:    job.Comment,
			URL:        job.URL,
			Payload:    job.Payload,
			Timeout:    job.Timeout,
			CreatedAt:  job.CreatedAt,
			FinishedAt: job.FinishedAt,
			Failure:    job.Failure,
			Success:    job.Success,
			StatusCode: job.StatusCode,
			Err:        job.Err,
			Output:     job.Output,
		}

		if err := boltutil.Set(tx, []interface{}{BucketNameForJobs}, job.ID, in); err != nil {
			return err
		}

		return nil
	})
}

func (s *Store) UpdateJob(job *hq.Job) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if err := boltutil.Get(tx, []interface{}{BucketNameForJobs}, job.ID, &hq.J{}); err != nil {
			if err == boltutil.ErrNotFound {
				return &ErrJobNotFound{ID: job.ID}
			} else {
				return err
			}
		}

		in := &hq.J{
			ID:         job.ID,
			Name:       job.Name,
			Comment:    job.Comment,
			URL:        job.URL,
			Payload:    job.Payload,
			Timeout:    job.Timeout,
			CreatedAt:  job.CreatedAt,
			FinishedAt: job.FinishedAt,
			Failure:    job.Failure,
			Success:    job.Success,
			StatusCode: job.StatusCode,
			Err:        job.Err,
			Output:     job.Output,
		}

		if err := boltutil.Set(tx, []interface{}{BucketNameForJobs}, job.ID, in); err != nil {
			return err
		}

		return nil
	})
}

func (s *Store) DeleteJob(id uint64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if err := boltutil.Get(tx, []interface{}{BucketNameForJobs}, id, &hq.J{}); err != nil {
			if err == boltutil.ErrNotFound {
				return &ErrJobNotFound{ID: id}
			} else {
				return err
			}
		}

		if err := boltutil.Delete(tx, []interface{}{BucketNameForJobs}, id); err != nil {
			return err
		}

		return nil
	})
}

func (s *Store) FetchJob(id uint64, job *hq.Job) error {
	qm := s.app.QueueManager
	return s.db.View(func(tx *bolt.Tx) error {
		out := &hq.J{}
		if err := boltutil.Get(tx, []interface{}{BucketNameForJobs}, id, out); err != nil {
			if err == boltutil.ErrNotFound {
				return &ErrJobNotFound{ID: id}
			} else {
				return err
			}
		}

		job.ID = out.ID
		job.Name = out.Name
		job.Comment = out.Comment
		job.URL = out.URL
		job.Payload = out.Payload
		job.Timeout = out.Timeout
		job.CreatedAt = out.CreatedAt
		job.Failure = out.Failure
		job.FinishedAt = out.FinishedAt
		job.Success = out.Success
		job.StatusCode = out.StatusCode
		job.Err = out.Err
		job.Output = out.Output

		job = qm.UpdateRunningStatus(job)

		return nil
	})
}

type ListJobsQuery struct {
	Name    string
	Begin   *uint64
	Reverse bool
	Limit   int
}

func (s *Store) ListJobs(query *ListJobsQuery, ret *hq.JobList) error {
	return s.db.View(func(tx *bolt.Tx) error {
		c, err := boltutil.Cursor(tx, []interface{}{BucketNameForJobs})
		if err != nil {
			if err == boltutil.ErrNotFound {
				return nil
			} else {
				return err
			}
		}

		if query.Reverse {
			if query.Begin != nil {
				begin := *query.Begin
				beginB, err := boltutil.ToKeyBytes(begin)
				if err != nil {
					return err
				}

				var k, v []byte
				if c.Bucket().Get(beginB) != nil {
					k, v = c.Seek(beginB)
				} else {
					// If the seeking key does not exist then the next key is used.
					k, v = c.Seek(beginB)
					if k == nil {
						k, v = c.Last()
					} else {
						k, v = c.Prev()
					}
				}

				for ; k != nil; k, v = c.Prev() {
					if err := s.appendJob(v, query, ret); err != nil {
						return err
					}

					if query.Limit > 0 && len(ret.Jobs) >= query.Limit {
						break
					}
				}
			} else {
				for k, v := c.Last(); k != nil; k, v = c.Prev() {
					if err := s.appendJob(v, query, ret); err != nil {
						return err
					}

					if query.Limit > 0 && len(ret.Jobs) >= query.Limit {
						break
					}
				}
			}

			if k, _ := c.Prev(); k != nil {
				ret.HasNext = true
				n := binary.BigEndian.Uint64(k)
				ret.Next = &n
			} else {
				ret.HasNext = false
			}
		} else {
			if query.Begin != nil {
				begin := *query.Begin
				beginB, err := boltutil.ToKeyBytes(begin)
				if err != nil {
					return err
				}

				for k, v := c.Seek(beginB); k != nil; k, v = c.Next() {
					if err := s.appendJob(v, query, ret); err != nil {
						return err
					}

					if query.Limit > 0 && len(ret.Jobs) >= query.Limit {
						break
					}
				}
			} else {
				for k, v := c.First(); k != nil; k, v = c.Next() {
					if err := s.appendJob(v, query, ret); err != nil {
						return err
					}

					if query.Limit > 0 && len(ret.Jobs) >= query.Limit {
						break
					}
				}
			}

			if k, _ := c.Next(); k != nil {
				ret.HasNext = true
				n := binary.BigEndian.Uint64(k)
				ret.Next = &n
			} else {
				ret.HasNext = false
			}
		}

		ret.Count = len(ret.Jobs)

		return nil
	})
}

func (s *Store) appendJob(v []byte, query *ListJobsQuery, ret *hq.JobList) error {
	qm := s.app.QueueManager

	in := &hq.J{}
	if err := boltutil.Deserialize(v, in); err != nil {
		return err
	}

	job := &hq.Job{
		ID:         in.ID,
		Name:       in.Name,
		Comment:    in.Comment,
		URL:        in.URL,
		Payload:    in.Payload,
		Timeout:    in.Timeout,
		CreatedAt:  in.CreatedAt,
		FinishedAt: in.FinishedAt,
		Failure:    in.Failure,
		Success:    in.Success,
		StatusCode: in.StatusCode,
		Err:        in.Err,
		Output:     in.Output,
	}

	job = qm.UpdateRunningStatus(job)

	if query.Name == "" {
		ret.Jobs = append(ret.Jobs, job)
		return nil
	}

	// filter job name
	r, err := regexp.Compile(query.Name)
	if err != nil {
		return err
	}

	if r.MatchString(job.Name) {
		ret.Jobs = append(ret.Jobs, job)
		return nil
	}

	return nil
}

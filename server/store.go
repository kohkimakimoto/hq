package server

import (
	"encoding/binary"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/kayac/go-katsubushi"
	"github.com/kohkimakimoto/boltutil"
	"github.com/kohkimakimoto/hq/structs"
	"github.com/labstack/echo"
	"regexp"
)

type Store struct {
	app    *App
	db     *bolt.DB
	logger echo.Logger
	gen    katsubushi.Generator
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

func (s *Store) CreateJob(job *structs.Job) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if err := boltutil.Get(tx, []interface{}{BucketNameForJobs}, job.ID, &structs.J{}); err == nil {
			return &ErrJobAlreadyExisted{ID: job.ID, Name: job.Name}
		}

		in := &structs.J{
			ID:         job.ID,
			Name:       job.Name,
			Comment:    job.Comment,
			Code:       job.Code,
			CreatedAt:  job.CreatedAt,
			FinishedAt: job.FinishedAt,
			Failure:    job.Failure,
			Success:    job.Success,
			Err:        job.Err,
			Output:     job.Output,
		}

		if err := boltutil.Set(tx, []interface{}{BucketNameForJobs}, job.ID, in); err != nil {
			return err
		}

		return nil
	})
}

func (s *Store) UpdateJob(job *structs.Job) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if err := boltutil.Get(tx, []interface{}{BucketNameForJobs}, job.ID, &structs.J{}); err != nil {
			if err == boltutil.ErrNotFound {
				return &ErrJobNotFound{ID: job.ID}
			} else {
				return err
			}
		}

		in := &structs.J{
			ID:         job.ID,
			Name:       job.Name,
			Comment:    job.Comment,
			Code:       job.Code,
			CreatedAt:  job.CreatedAt,
			FinishedAt: job.FinishedAt,
			Failure:    job.Failure,
			Success:    job.Success,
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
		if err := boltutil.Get(tx, []interface{}{BucketNameForJobs}, id, &structs.J{}); err != nil {
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

func (s *Store) FetchJob(id uint64, job *structs.Job) error {
	//jm := s.srv.RuntimeJobManager
	return s.db.View(func(tx *bolt.Tx) error {
		out := &structs.J{}
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
		job.Code = out.Code
		job.CreatedAt = out.CreatedAt
		job.Failure = out.Failure
		job.FinishedAt = out.FinishedAt
		job.Success = out.Success
		job.Err = out.Err
		job.Output = out.Output

		//		job = jm.SetRuntimeInfo(job)

		return nil
	})
}

func (s *Store) ListJobs(query *structs.ListJobsQuery, ret *structs.JobList) error {
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
			if query.HasBegin {
				beginB, err := boltutil.ToKeyBytes(query.Begin)
				if err != nil {
					return err
				}

				for k, v := c.Seek(beginB); k != nil; k, v = c.Prev() {
					if err := s.appendJob(v, query, ret); err != nil {
						return err
					}

					if len(ret.Jobs) >= query.Limit {
						break
					}
				}
			} else {
				for k, v := c.Last(); k != nil; k, v = c.Prev() {
					if err := s.appendJob(v, query, ret); err != nil {
						return err
					}

					if len(ret.Jobs) >= query.Limit {
						break
					}
				}
			}

			if k, _ := c.Prev(); k != nil {
				ret.HasNext = true
				n := binary.BigEndian.Uint64(k)
				ret.NextJob = &n
			} else {
				ret.HasNext = false
			}
		} else {
			if query.HasBegin {
				beginB, err := boltutil.ToKeyBytes(query.Begin)
				if err != nil {
					return err
				}

				for k, v := c.Seek(beginB); k != nil; k, v = c.Next() {
					if err := s.appendJob(v, query, ret); err != nil {
						return err
					}

					if len(ret.Jobs) >= query.Limit {
						break
					}
				}
			} else {
				for k, v := c.First(); k != nil; k, v = c.Next() {
					if err := s.appendJob(v, query, ret); err != nil {
						return err
					}

					if len(ret.Jobs) >= query.Limit {
						break
					}
				}
			}

			if k, _ := c.Next(); k != nil {
				ret.HasNext = true
				n := binary.BigEndian.Uint64(k)
				ret.NextJob = &n
			} else {
				ret.HasNext = false
			}
		}

		ret.Count = len(ret.Jobs)

		return nil
	})
}

func (s *Store) appendJob(v []byte, query *structs.ListJobsQuery, ret *structs.JobList) error {
	// jm := s.srv.RuntimeJobManager

	in := &structs.J{}
	if err := boltutil.Deserialize(v, in); err != nil {
		return err
	}

	job := &structs.Job{
		ID:         in.ID,
		Name:       in.Name,
		Comment:    in.Comment,
		Code:       in.Code,
		CreatedAt:  in.CreatedAt,
		FinishedAt: in.FinishedAt,
		Failure:    in.Failure,
		Success:    in.Success,
		Err:        in.Err,
		Output:     in.Output,
	}

	//job = jm.SetRuntimeInfo(job)

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

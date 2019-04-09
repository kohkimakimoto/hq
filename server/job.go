package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cjoudrey/gluahttp"
	"github.com/kohkimakimoto/gluaenv"
	"github.com/kohkimakimoto/hq/structs"
	"github.com/labstack/echo"
	"github.com/otm/gluash"
	glualibs "github.com/vadv/gopher-lua-libs"
	"github.com/yuin/gopher-lua"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type JobWorkerManager struct {
	app             *App
	mutex           *sync.Mutex
	queue           chan *structs.Job
	NumExecutorsAll int64
}

func NewJobWorkerManager(app *App) *JobWorkerManager {
	return &JobWorkerManager{
		app:             app,
		mutex:           new(sync.Mutex),
		queue:           make(chan *structs.Job, app.Config.Queues),
		NumExecutorsAll: 0,
	}
}

func (m *JobWorkerManager) Start() {
	config := m.app.Config

	for i := int64(0); i < config.Workers; i++ {
		worker := &JobWorker{
			m:            m,
			NumExecutors: 0,
		}

		go worker.Start()
	}
}

func (m *JobWorkerManager) Enqueue(job *structs.Job) {
	m.queue <- job
}

type JobWorker struct {
	m            *JobWorkerManager
	NumExecutors int64
}

func (w *JobWorker) Start() {
	m := w.m
	logger := m.app.Logger
	config := m.app.Config

	for {
		job := <-w.m.queue
		logger.Debugf("dequeue job: %d", job.ID)

		if atomic.LoadInt64(&config.MaxExecutors) <= 0 {
			// sync
			w.executeJob(job)
		} else if atomic.LoadInt64(&w.NumExecutors) < atomic.LoadInt64(&config.MaxExecutors) {
			// async
			atomic.AddInt64(&w.NumExecutors, 1)
			atomic.AddInt64(&m.NumExecutorsAll, 1)

			go func(job *structs.Job) {
				w.executeJob(job)
				atomic.AddInt64(&w.NumExecutors, -1)
				atomic.AddInt64(&m.NumExecutorsAll, -1)
			}(job)

		} else {
			// sync
			w.executeJob(job)
		}
	}
}

func (w *JobWorker) executeJob(job *structs.Job) {
	logger := w.m.app.Logger

	L := lua.NewState()
	defer L.Close()

	// modules
	glualibs.Preload(L)
	L.PreloadModule("env", gluaenv.Loader)
	L.PreloadModule("sh", gluash.Loader)
	L.PreloadModule("httpclient", gluahttp.NewHttpModule(&http.Client{}).Loader)

	e := &Executor{
		App:       w.m.app,
		L:         L,
		Job:       job,
		OutBuffer: &bytes.Buffer{},
		Logger:    logger,
	}

	L.SetGlobal("print", L.NewFunction(luaPrint(e)))
	if err := L.DoString(`
-- disabled os.exit
os.exit = nil
`); err != nil {
		logger.Error(err)
	}

	// set context
	var ctx context.Context
	var cancel context.CancelFunc
	if job.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(job.Timeout)*time.Second)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	L.SetContext(ctx)
	e.CancelFunc = cancel

	defer e.Close()
	if err := e.Run(); err != nil {
		e.Err = err
		logger.Error(err)
	}
}

type Executor struct {
	App        *App
	L          *lua.LState
	Job        *structs.Job
	OutBuffer  *bytes.Buffer
	Err        error
	Logger     echo.Logger
	CancelFunc func()
}

func (e *Executor) Run() error {
	e.Logger.Infof("job: %d executing", e.Job.ID)
	defer e.Logger.Infof("job: %d executed", e.Job.ID)

	if err := e.L.DoString(e.Job.Code); err != nil {
		return err
	}
	return nil
}

func (e *Executor) Close() {
	logger := e.Logger
	logger.Debugf("job: %d closing", e.Job.ID)

	err := e.Err
	job := e.Job
	outBuffer := e.OutBuffer
	store := e.App.Store

	// Update result status (success or failure).
	// If the evaluator has an error, write it to the output buf.
	if err != nil {
		job.Success = false
		job.Failure = true

		if outBuffer != nil {
			job.Err = strings.Replace(err.Error(), "\n", "", -1)
			fmt.Fprintf(outBuffer, err.Error())
		}
	} else {
		job.Success = true
		job.Failure = false
	}

	// Truncate millisecond. It is compatible time for katsubushi ID generator time stamp.
	now := time.Now().UTC().Truncate(time.Millisecond)

	// update finishedAt
	job.FinishedAt = &now
	if outBuffer != nil {
		job.Output = outBuffer.String()
	}

	if e := store.UpdateJob(job); e != nil {
		err = e
	}

	e.CancelFunc()

	logger.Debugf("job: %d closed", e.Job.ID)
}

func luaPrint(e *Executor) func(*lua.LState) int {
	return func(L *lua.LState) int {
		top := L.GetTop()
		for i := 1; i <= top; i++ {
			fmt.Fprint(e.OutBuffer, L.ToStringMeta(L.Get(i)).String())
			if i != top {
				fmt.Fprint(e.OutBuffer, "\t")
			}
		}
		fmt.Fprintln(e.OutBuffer, "")
		return 0
	}
}

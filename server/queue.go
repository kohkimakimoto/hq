package server

import (
	"bytes"
	"fmt"
	"github.com/kohkimakimoto/hq/structs"
	"net/http"
	"sync/atomic"
	"time"
)

type QueueManager struct {
	app           *App
	queue         chan *structs.Job
	dispatchers   []*Dispatcher
	numWorkersAll int64
}

func NewQueueManager(app *App) *QueueManager {
	return &QueueManager{
		app:           app,
		queue:         make(chan *structs.Job, app.Config.Queues),
		dispatchers:   []*Dispatcher{},
		numWorkersAll: 0,
	}
}

func (m *QueueManager) Start() {
	config := m.app.Config

	for i := int64(0); i < config.Dispatchers; i++ {
		d := &Dispatcher{
			manager:    m,
			numWorkers: 0,
		}

		go d.loop()
	}
}

func (m *QueueManager) Enqueue(job *structs.Job) {
	m.queue <- job
}

type Dispatcher struct {
	manager    *QueueManager
	numWorkers int64
}

func (d *Dispatcher) loop() {
	m := d.manager
	logger := m.app.Logger
	config := m.app.Config

	for {
		job := <-m.queue
		logger.Debugf("dequeue job: %d", job.ID)

		if atomic.LoadInt64(&config.MaxWorkers) <= 0 {
			// sync
			d.work(job)
		} else if atomic.LoadInt64(&d.numWorkers) < atomic.LoadInt64(&config.MaxWorkers) {
			// async
			atomic.AddInt64(&d.numWorkers, 1)
			atomic.AddInt64(&m.numWorkersAll, 1)

			go func(job *structs.Job) {
				d.work(job)
				atomic.AddInt64(&d.numWorkers, -1)
				atomic.AddInt64(&m.numWorkersAll, -1)
			}(job)
		} else {
			// sync
			d.work(job)
		}
	}
}

func (d *Dispatcher) work(job *structs.Job) {
	app := d.manager.app
	logger := app.Logger
	store := app.Store

	// worker variables
	var err error
	var output string
	logger.Infof("job: %d working", job.ID)

	// the terminating logic
	defer func() {
		logger.Infof("job: %d worked", job.ID)
		logger.Debugf("job: %d closing", job.ID)

		// Update result status (success or failure).
		// If the evaluator has an error, write it to the output buf.
		if err != nil {
			job.Success = false
			job.Failure = true
			job.Err = err.Error()
		} else {
			job.Success = true
			job.Failure = false
		}

		// Truncate millisecond. It is compatible time for katsubushi ID generator time stamp.
		now := time.Now().UTC().Truncate(time.Millisecond)

		// update finishedAt
		job.FinishedAt = &now
		job.Output = output

		if e := store.UpdateJob(job); e != nil {
			logger.Error(e)
		}

		logger.Debugf("job: %d closed", job.ID)
	}()

	// worker
	client := &http.Client{
		Timeout: time.Duration(job.Timeout) * time.Second,
	}
	req, err := http.NewRequest(
		"POST",
		job.URL,
		bytes.NewReader(job.Payload),
	)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return
	//}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(http.StatusText(resp.StatusCode))
		return
	}
}

//
//
//type Worker struct {
//	app *App
//	job *structs.Job
//	Err error
//}
//
//
//func (w *Dispatcher) aa(job *structs.Job) {
//	logger := w.m.app.Logger
//
//	L := lua.NewState()
//	defer L.Close()
//
//	// modules
//	glualibs.Preload(L)
//	L.PreloadModule("env", gluaenv.Loader)
//	L.PreloadModule("sh", gluash.Loader)
//	L.PreloadModule("httpclient", gluahttp.NewHttpModule(&http.Client{}).Loader)
//
//	e := &Executor{
//		App:       w.m.app,
//		L:         L,
//		Job:       job,
//		OutBuffer: &bytes.Buffer{},
//		Logger:    logger,
//	}
//
//	L.SetGlobal("print", L.NewFunction(luaPrint(e)))
//	if err := L.DoString(`
//-- disabled os.exit
//os.exit = nil
//`); err != nil {
//		logger.Error(err)
//	}
//
//	// set context
//	var ctx context.Context
//	var cancel context.CancelFunc
//	if job.Timeout > 0 {
//		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(job.Timeout)*time.Second)
//	} else {
//		ctx, cancel = context.WithCancel(context.Background())
//	}
//	L.SetContext(ctx)
//	e.CancelFunc = cancel
//
//	defer e.Close()
//	if err := e.Run(); err != nil {
//		e.Err = err
//		logger.Error(err)
//	}
//}
//
//type Executor struct {
//	App        *App
//	L          *lua.LState
//	Job        *structs.Job
//	OutBuffer  *bytes.Buffer
//	Err        error
//	Logger     echo.Logger
//	CancelFunc func()
//}
//
//func (e *Executor) Run() error {
//	e.Logger.Infof("job: %d executing", e.Job.ID)
//	defer e.Logger.Infof("job: %d executed", e.Job.ID)
//
//	if err := e.L.DoString(e.Job.Code); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (e *Executor) Close() {
//	logger := e.Logger
//	logger.Debugf("job: %d closing", e.Job.ID)
//
//	err := e.Err
//	job := e.Job
//	outBuffer := e.OutBuffer
//	store := e.App.Store
//
//	// Update result status (success or failure).
//	// If the evaluator has an error, write it to the output buf.
//	if err != nil {
//		job.Success = false
//		job.Failure = true
//
//		if outBuffer != nil {
//			job.Err = strings.Replace(err.Error(), "\n", "", -1)
//			fmt.Fprintf(outBuffer, err.Error())
//		}
//	} else {
//		job.Success = true
//		job.Failure = false
//	}
//
//	// Truncate millisecond. It is compatible time for katsubushi ID generator time stamp.
//	now := time.Now().UTC().Truncate(time.Millisecond)
//
//	// update finishedAt
//	job.FinishedAt = &now
//	if outBuffer != nil {
//		job.Output = outBuffer.String()
//	}
//
//	if e := store.UpdateJob(job); e != nil {
//		err = e
//	}
//
//	e.CancelFunc()
//
//	logger.Debugf("job: %d closed", e.Job.ID)
//}
//
//func luaPrint(e *Executor) func(*lua.LState) int {
//	return func(L *lua.LState) int {
//		top := L.GetTop()
//		for i := 1; i <= top; i++ {
//			fmt.Fprint(e.OutBuffer, L.ToStringMeta(L.Get(i)).String())
//			if i != top {
//				fmt.Fprint(e.OutBuffer, "\t")
//			}
//		}
//		fmt.Fprintln(e.OutBuffer, "")
//		return 0
//	}
//}

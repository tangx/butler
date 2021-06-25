package butler

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type Butler struct {
	// active workers
	workers int
	// max workers
	workersCap  int
	workerQueue chan *worker

	// max jobs queue
	jobs     int
	jobQueue chan func()

	ctx context.Context
	wg  sync.WaitGroup
}

// Default return a butler , which workers' number is equal GOMAXPROC
// and jobs' number is double of workers
// https://golang.org/doc/effective_go#parallel
func Default() *Butler {
	b := &Butler{}
	b.initial()
	return b
}

// WithOptions return a new Butler
func (b *Butler) WithOptions(funcs ...OptionFunc) *Butler {
	for _, fn := range funcs {
		fn(b)
	}
	return b
}

func (b *Butler) Init() {
	b.initial()
}

func (b *Butler) trace() {
	for {
		time.Sleep(500 * time.Millisecond)
		fmt.Println(len(b.workerQueue), cap(b.workerQueue), "<===>", len(b.jobQueue), cap(b.jobQueue))
	}
}

// Work start
func (b *Butler) Work() {
	// https://colobu.com/2015/10/09/Linux-Signals/
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// go b.trace()

Loop:
	for {
		select {
		// catch a signal and break out loop
		case sig := <-sigs:
			log.Printf(">>>>>>>> catch signal %v \n", sig)
			break Loop
		// ctx timeout
		case <-b.ctx.Done():
			log.Printf(">>>>>>>> context cancel %v \n", b.ctx.Err())
			break Loop

		// work when worker and jobs are both ready
		case job := <-b.jobQueue:

		JobLoop:
			// pervent block by worker when worker queue empty
			for {
				select {
				// catch a signal and break out loop
				case sig := <-sigs:
					log.Printf(">>>>>>>> catch signal %v \n", sig)
					break Loop
				// ctx timeout
				case <-b.ctx.Done():
					log.Printf(">>>>>>>> context cancel %v \n", b.ctx.Err())
					break Loop
				case worker := <-b.workerQueue:
					job := job
					go b.assign(worker, job)
					break JobLoop
				default:
					// try to create a new worker
					b.hire()
				}
			}
		}

	}

	// wait all jobs done
	b.wg.Wait()

}

func (b *Butler) AddJobs(funcs ...func()) {
	for _, fn := range funcs {
		b.jobQueue <- fn
	}
}

// SetDefaults set default value for butler
func (b *Butler) SetDefaults() {
	if b.workersCap < 1 {
		b.workersCap = runtime.GOMAXPROCS(0)
	}

	b.jobs = b.workersCap * 2

	if b.ctx == nil {
		b.ctx = context.Background()
	}
}

// initial butler
func (b *Butler) initial() {
	b.SetDefaults()

	b.jobQueue = make(chan func(), b.jobs)

	b.workerQueue = make(chan *worker, b.workersCap)

}

// assign a job to a worker
func (b *Butler) assign(w *worker, job func()) {
	b.wg.Add(1)
	defer b.wg.Done()

	defer func() {
		if err := recover(); err != nil {
			log.Printf("catch panic: %v", err)
		}

		b.workerQueue <- w
	}()

	w.do(job)
}

func (b *Butler) hire() {
	if b.workers < b.workersCap {
		log.Printf("<<<<<--- hire a new worker\n")
		b.workerQueue <- newWorker()
		b.workers++
	}
}

// OptionFunc
type OptionFunc = func(b *Butler)

// WithWorkers set concurrency worker numbers
func WithWorkers(n int) OptionFunc {
	return func(b *Butler) {
		b.workersCap = n
	}
}

// WithJobs set max jobs queue
func WithJobs(n int) OptionFunc {
	return func(b *Butler) {
		b.jobs = n
	}
}

func WithContext(ctx context.Context) OptionFunc {
	if ctx == nil {
		ctx = context.Background()
	}
	return func(b *Butler) {
		b.ctx = ctx
	}
}

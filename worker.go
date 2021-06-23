package bulter

type worker struct{}

// newWorker return a new worker
func newWorker() *worker {
	return &worker{}
}

// do stuffs
func (w *worker) do(job func()) {
	job()
}

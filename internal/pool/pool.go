package pool

import (
	"fmt"
	"log/slog"
)

type Pool[Data any] struct {
	poolSize int
	pool     chan *Worker
	workers  []*Worker
	handler  func(int, Data) error
}

func NewPool[Data any](size int, handler func(int, Data) error) *Pool[Data] {
	return &Pool[Data]{
		poolSize: size,
		pool:     make(chan *Worker, size),
		handler:  handler,
	}
}

func (pl *Pool[Data]) Create() {
	for i := range pl.poolSize {
		worker := &Worker{id: i + 1}
		pl.workers = append(pl.workers, worker)
		pl.pool <- worker
	}
}

func (pl *Pool[Data]) Handle(data Data) {
	op := "pool.Handle"
	w := <-pl.pool

	go func() {
		if err := pl.handler(w.id, data); err != nil {
			w.failed++
			slog.Error("Failed to processing", "error", err, "operation", op)
		} else {
			w.completed++
		}
		w.total++
		pl.pool <- w
	}()
}

func (pl *Pool[Data]) Wait() {
	for range len(pl.workers) {
		<-pl.pool
	}
}

func (pl *Pool[Data]) Stats() {
	fmt.Println("____________RESULT____________")
	for _, worker := range pl.workers {
		fmt.Printf(
			"worker: %d; total: %d; completed %d; failed %d;\n",
			worker.id,
			worker.total,
			worker.completed,
			worker.failed,
		)
	}
	fmt.Println("_____________________________")
}

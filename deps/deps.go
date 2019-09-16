package deps

import (
	"github.com/brad-jones/goasync/await"
	"github.com/brad-jones/goasync/task"
	"github.com/go-errors/errors"
)

func Serial(tasks ...func(t *task.Internal)) func(t *task.Internal) {
	return func(t *task.Internal) {
		for _, v := range tasks {
			if t.ShouldStop() {
				t.Reject(errors.New("deps: told to stop before all tasks could be executed"))
				break
			}
			t2 := task.New(v)
			t2.Stopper = t.Stopper
			t2.MustResult()
		}
	}
}

func Parallel(tasks ...func(t *task.Internal)) func(t *task.Internal) {
	return func(t *task.Internal) {
		awaitables := []await.Awaitable{}
		for _, v := range tasks {
			if t.ShouldStop() {
				t.Reject(errors.New("deps: told to stop before all tasks could be executed"))
				break
			}
			t2 := task.New(v)
			t2.Stopper = t.Stopper
			awaitables = append(awaitables, t2)
		}
		await.MustFastAllOrError(awaitables...)
	}
}

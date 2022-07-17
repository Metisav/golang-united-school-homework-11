package batch

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

type user struct {
	ID int64
}

func getOne(id int64, ch chan user) error {
	time.Sleep(time.Millisecond * 100)
	ch <- user{ID: id}
	return nil
}

func getBatch(n int64, pool int64) (res []user) {
	var ch = make(chan user, n)
	errG, _ := errgroup.WithContext(context.Background())
	errG.SetLimit(int(pool))
	for i := 0; i < int(n); i++ {
		copyI := i
		errG.Go(func() error {
			return getOne(int64(copyI), ch)
		})
	}
	errG.Wait()
	close(ch)

	for us := range ch {
		res = append(res, us)
	}
	return res
}

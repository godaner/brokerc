package sync

import "sync"

type OnceError struct {
	sync.Once
}

func (o *OnceError) Do(f func() error) (err error) {
	o.Once.Do(func() {
		err = f()
	})
	return err
}

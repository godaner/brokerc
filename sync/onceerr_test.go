package sync

import (
	"errors"
	"fmt"
	"testing"
)

func TestOnceError_Do(t *testing.T) {
	once := OnceError{}
	fmt.Println(once.Do(func() error {
		return errors.New("mock err")
	}))
}

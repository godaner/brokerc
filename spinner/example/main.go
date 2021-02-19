package main

import (
	"github.com/godaner/brokerc/spinner"
	"time"
)

func main() {
	s := spinner.DownloadSpinner{}
	s.Start() // Start the spinner
	defer s.Stop()
	go func() {
		v := uint64(0)
		download := uint64(0)
		for ; ; {
			<-time.After(time.Second)
			v += 1000
			total := uint64(59959)
			download += 1000
			s.UpdateStatus(&spinner.Status{
				Download: &download,
				Total:    &total,
				V:        &v,
			})
		}
	}()
	<-time.After(1000 * time.Second)
}

package spinner

import (
	"fmt"
	"github.com/briandowns/spinner"
	"sync"
	"time"
)

type status struct {
	download, total, v uint64
	fileName           string
}
type Status struct {
	Download, Total, V *uint64
	FileName           *string
}
type Spinner struct {
	sync.Once
	status
	s *spinner.Spinner
}

func (s *Spinner) init() {
	s.Do(func() {
		s.s = spinner.New(s.computeStatus(), time.Second)
	})
}
func (s *Spinner) Start() {
	s.init()
	s.s.Start()
}
func (s *Spinner) Stop() {
	s.init()
	s.s.Stop()
}
func (s *Spinner) UpdateStatus(sts *Status) {
	s.init()
	if sts == nil {
		return
	}
	if sts.Total != nil {
		s.status.total = *sts.Total
	}
	if sts.Total != nil {
		s.status.download = *sts.Download
	}
	if sts.Total != nil {
		s.status.v = *sts.V
	}
	if sts.FileName != nil {
		s.status.fileName = *sts.FileName
	}
	s.s.UpdateCharSet(s.computeStatus())
}

func (s *Spinner) computeStatus() []string {
	p := int64(0)
	if s.total != 0 {
		p = int64((float64(s.download) / float64(s.total)) * 100)
	}
	return []string{fmt.Sprint("Downloading file: ", s.fileName, ", ", p, "%, ", formatSize(int64(s.download)), "/", formatSize(int64(s.total))+", ", formatSize(int64(s.v)), "/s")}
}

func formatSize(size int64) (s string) {
	if size < 1024 {
		// return strconv.FormatInt(size, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(size)/float64(1))
	} else if size < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(size)/float64(1024))
	} else if size < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(size)/float64(1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(size)/float64(1024*1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(size)/float64(1024*1024*1024*1024))
	} else { // if size < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(size)/float64(1024*1024*1024*1024*1024))
	}
}

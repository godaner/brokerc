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
	spentTime          time.Duration
}
type Status struct {
	Download, Total, V *uint64
	FileName           *string
	SpentTime          *time.Duration
}
type DownloadSpinner struct {
	sync.Once
	status
	s *spinner.Spinner
}

func (s *DownloadSpinner) init() {
	s.Do(func() {
		s.s = spinner.New(s.computeStatus(), time.Second)
	})
}
func (s *DownloadSpinner) Start() {
	s.init()
	s.s.Start()
}
func (s *DownloadSpinner) Stop() {
	s.init()
	s.s.Stop()
}
func (s *DownloadSpinner) UpdateStatus(sts *Status) {
	s.init()
	if sts == nil {
		return
	}
	if sts.Total != nil {
		s.status.total = *sts.Total
	}
	if sts.Download != nil {
		s.status.download = *sts.Download
	}
	if sts.V != nil {
		s.status.v = *sts.V
	}
	if sts.FileName != nil {
		s.status.fileName = *sts.FileName
	}
	if sts.SpentTime != nil {
		s.status.spentTime = *sts.SpentTime
	}
	s.s.UpdateCharSet(s.computeStatus())
}

func (s *DownloadSpinner) computeStatus() []string {
	p := int64(0)
	if s.total != 0 {
		p = int64((float64(s.download) / float64(s.total)) * 100)
	}
	leftTime := time.Hour * 999
	if s.v != 0 {
		leftTime = time.Duration((s.total-s.download)/s.v) * time.Second
	}
	return []string{fmt.Sprint("Downloading file: ", s.fileName, ", ", p, "%, ", formatSize(int64(s.download)), "/", formatSize(int64(s.total))+", ", formatSize(int64(s.v)), "/s, spent time: ", (s.spentTime/time.Second)*time.Second, ", left time: ", leftTime)}
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

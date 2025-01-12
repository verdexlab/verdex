//go:build !js
// +build !js

package ui

import (
	"io"

	"github.com/cheggaaa/pb/v3"
)

type ProgressBar struct {
	ProgressBar *pb.ProgressBar
	Reader      io.Reader
}

func ProgressBarStart(r io.Reader, total int64) ProgressBar {
	bar := pb.New64(total).SetMaxWidth(100)
	bar.Start()
	return ProgressBar{
		ProgressBar: bar,
		Reader:      bar.NewProxyReader(r).Reader,
	}
}

func ProgressBarFinish(bar ProgressBar) {
	bar.ProgressBar.Finish()
}

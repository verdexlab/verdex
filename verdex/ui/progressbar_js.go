//go:build js
// +build js

package ui

import (
	"io"
)

type ProgressBar struct {
	Reader io.Reader
}

func ProgressBarStart(r io.Reader, total int64) ProgressBar {
	return ProgressBar{
		Reader: r,
	}
}

func ProgressBarFinish(bar ProgressBar) {
	// nothing
}

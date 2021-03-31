package main

import (
	"fmt"
	"os"
)

type output struct {
	filePath     string
	writtenFiles []string
}

func newOutput(filePath string) *output {
	return &output{filePath: filePath}
}

func (out *output) CreateFile(outputParameters ...interface{}) (f *os.File, err error) {
	outputPath := fmt.Sprintf(out.filePath) // TODO - check parameter count vs. placeholders in output file path
	f, err = os.Create(outputPath)
	if err == nil {
		out.writtenFiles = append(out.writtenFiles, outputPath)
	}
	return
}

func (out *output) WrittenFiles() (info []os.FileInfo, err error) {
	info = make([]os.FileInfo, len(out.writtenFiles))
	for i, filePath := range out.writtenFiles {
		var s os.FileInfo
		s, err = os.Stat(filePath)
		if err != nil {
			return
		}
		info[i] = s
	}
	return
}

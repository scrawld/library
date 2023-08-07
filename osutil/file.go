package osutil

import (
	"fmt"
	"os"
	"sync"
)

// File is an io.WriteCloser that writes to the specified filename.
// opens or creates the logfile on first Write.
type File struct {
	Filename string `json:"filename"`

	file *os.File
	mu   sync.Mutex
}

func NewFile(filename string) *File {
	return &File{Filename: filename}
}

// Write to the file and create it if it does not exist
func (f *File) Write(p []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.file == nil {
		f.file, err = os.OpenFile(f.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return 0, fmt.Errorf("open file error: %s", err)
		}
	}
	return f.file.Write(p)
}

// Close implements io.Closer, and closes the current logfile.
func (f *File) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.file != nil {
		err := f.file.Close()
		f.file = nil
		return err
	}
	return nil
}

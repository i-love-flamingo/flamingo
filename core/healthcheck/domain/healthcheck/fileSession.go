package healthcheck

import (
	"os"
)

type (
	// FileSession is the healthcheck for session file handling
	FileSession struct {
		fileName string
	}
)

var (
	_ Status = &FileSession{}
)

// Inject dependencies
func (s *FileSession) Inject(cfg *struct {
	FileName string `inject:"config:session.file"`
}) {
	s.fileName = cfg.FileName
}

// Status checks if the session file is available
func (s *FileSession) Status() (bool, string) {
	_, err := os.Stat(s.fileName)
	if err == nil {
		return true, "success"
	}

	return false, err.Error()
}

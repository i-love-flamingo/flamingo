package healthcheck

import (
	"os"

	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
)

// FileSession session backend health check
type FileSession struct {
	fileName string
}

var _ healthcheck.Status = &FileSession{}

// Inject configuration for session backend
func (s *FileSession) Inject(cfg *struct {
	FileName string `inject:"config:flamingo.session.file"`
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

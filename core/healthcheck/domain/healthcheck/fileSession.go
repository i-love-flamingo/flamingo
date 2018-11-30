package healthcheck

import (
	"os"
)

type (
	FileSession struct {
		fileName string
	}
)

var (
	_ Status = &FileSession{}
)

func (s *FileSession) Inject(cfg *struct {
	FileName string `inject:"config:session.file"`
}) {
	s.fileName = cfg.FileName
}

func (s *FileSession) Status() (bool, string) {
	_, err := os.Stat(s.fileName)
	if err == nil {
		return true, "success"
	}

	return false, err.Error()
}

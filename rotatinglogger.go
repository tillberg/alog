package alog

import (
	"errors"
	"os"
	"path/filepath"
)

const ROTATE_SIZE = 10 * (1 << 20)

type RotatingLogger struct {
	*Logger
	loggerInt PrintLogger
	path      string
	file      *os.File
	size      int64
}

func NewRotatingLogger(path string, loggerInt PrintLogger) (*RotatingLogger, error) {
	var err error
	l := &RotatingLogger{}
	l.path = path
	l.loggerInt = loggerInt
	stat, err := os.Stat(l.path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		// Make the directory, in case it doesn't exist yet
		os.MkdirAll(filepath.Dir(path), 0755)
	} else {
		l.size = stat.Size()
	}
	l.Logger = New(l, "@(dim:{isodate}) ", 0)
	err = l.openfile()
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (l *RotatingLogger) rotate() {
	os.Rename(l.path, l.path+".old")
	err := l.openfile()
	if err != nil {
		l.loggerInt.Printf("(@error:Error opening new log file %s on rotation: %v)\n", l.path, err)
	}
	l.size = 0
}

func (l *RotatingLogger) openfile() error {
	var err error
	if l.file != nil {
		go l.file.Close()
	}
	l.file, err = os.OpenFile(l.path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (l *RotatingLogger) Write(buf []byte) (int, error) {
	// Println(l)
	if l.file == nil {
		return 0, errors.New("Logfile not open")
	}
	nn, err := l.file.Write(buf)
	l.size += int64(nn)
	if l.size > ROTATE_SIZE {
		l.rotate()
	}
	if err != nil {
		l.loggerInt.Printf("(@error:Error writing to log file: %v)\n", err)
	}
	return nn, err
}

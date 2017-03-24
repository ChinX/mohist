package files

import (
	"io"
	"os"
	"path/filepath"
	"bufio"
	"bytes"
	"syscall"
)

type LineHandler func([]byte) bool

// IsDirExist checks if a dir exists
func IsDirExist(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir() || os.IsExist(err)
}

// IsFileExist checks if a file exists
func IsFileExist(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir() || os.IsExist(err)
}

// CreateDir create a dir
func CreateDir(path string) error {
	path = filepath.Dir(path)
	if _, err := os.Stat(path); err != nil && !os.IsExist(err) {
		return os.Mkdir(path, 0750)
	}
	return nil
}

// CopyFile copy content to new file
func CopyFile(from string, to string) error {
	ff, err := os.Open(from)
	if err != nil {
		return err
	}
	tf, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err == nil {
		ff.Close()
		return err
	}
	_, err = io.CopyBuffer(tf, ff, make([]byte, 4096))
	if err1 := ff.Close(); err == nil {
		err = err1
	}
	if err1 := tf.Close(); err == nil {
		err = err1
	}
	return err
}

// AppendFile Append content to file
func AppendFile(path string, byteArr []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	n, err := w.Write(byteArr)
	if err == nil && n < len(byteArr) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// InsertLine insert content to file on condition
func InsertLine(path string, condition, byteArr []byte, offset int64) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	if err == nil {
		h := func(line []byte) bool {
			if bytes.Contains(line, condition) {
				f.Seek(offset, syscall.FILE_CURRENT)
				f.Write(byteArr)
				f.Seek(-offset, syscall.FILE_CURRENT)
				return true
			}
			return false
		}

		if err = readLine(bufio.NewReader(f), h); err == nil {
			f.Sync()
		}
	}

	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// ReadLine Read file content in line
func ReadLine(path string, handler LineHandler) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	err = readLine(bufio.NewReader(f), handler)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// FindLine Find bytes of file line
func FindLine(path string, byteArr []byte) (fl []byte, err error) {
	h := func(line []byte) bool {
		fond := bytes.Contains(line, byteArr)
		if fond {
			fl = line
		}
		return fond
	}
	err = ReadLine(path, h)
	return
}

func readLine(buf *bufio.Reader, handler LineHandler) (err error) {
	var line []byte
	isOver := false
	for !isOver {
		line, _, err = buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		isOver = handler(line)
	}
	return
}

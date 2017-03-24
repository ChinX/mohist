package compress

import (
	"bytes"
	"compress/bzip2"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
)

// Zlib compress by zlib
func Zlib(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

// UnZlib un compress by zlib
func UnZlib(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

// Gzip compress by gzip
func Gzip(src []byte) []byte {
	var in bytes.Buffer
	w := gzip.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

// UnGzip un compress by gzip
func UnGzip(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := gzip.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

// UnBzip2 un compress by bzip2
func UnBzip2(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r := bzip2.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

// Flate compress by flate
func Flate(src []byte, level int) ([]byte, error) {
	var in bytes.Buffer
	w, err := flate.NewWriter(&in, level)
	if err != nil {
		return nil, err
	}
	w.Write(src)
	w.Close()
	return in.Bytes(), nil
}

// UnFlate un compress by flate
func UnFlate(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r := flate.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

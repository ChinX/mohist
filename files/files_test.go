package files

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/chinx/mohist/compress"
)

// main functions shows how to TarGz a directory/file and
// UnTarGz a file
// Gzip and tar from source directory or file to destination file
// you need check file exist before you call this function

func TestTarGz(t *testing.T) {
	os.Mkdir("/home/ty4z2008/tar", 0777)
	err := CopyFile("/home/ty4z2008/tar/1.pdf", "/home/ty4z2008/src/1.pdf")
	//targetfile,sourcefile
	if err != nil {
		fmt.Println(err.Error())
	}

	TarGz("/home/ty4z2008/tar/1.pdf", "/home/ty4z2008/test.tar.gz") //压缩
	//UnTarGz("/home/ty4z2008/1.tar.gz", "/home/ty4z2008")     //解压
	os.RemoveAll("/home/ty4z2008/tar")

	fmt.Println("ok")
}

func TestDoZlibUnCompress(t *testing.T) {
	buff := []byte{120, 156, 202, 72, 205, 201, 201, 215, 81, 40, 207,
		47, 202, 73, 225, 2, 4, 0, 0, 255, 255, 33, 231, 4, 147}
	b := bytes.NewReader(buff)
	r, err := zlib.NewReader(b)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, r)
	r.Close()

	zip := compress.Zlib([]byte("hello, world\n"))
	fmt.Println(zip)
	fmt.Println(string(compress.UnZlib(zip)))
}

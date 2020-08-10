package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

// 一般不会直接将上传的图片名暴露出来，因此对图片名进行 MD5 来达到这个效果
func Md5(text string) string {
	m := md5.New()
	m.Write([]byte(text))
	return hex.EncodeToString(m.Sum(nil))
}

func Md5Buffer(buf []byte) string {
	m := md5.New()
	m.Write(buf)
	return hex.EncodeToString(m.Sum(nil))
}

func Md5File(reader io.Reader) string {
	var buf = make([]byte, 4096)
	m := md5.New()
	for {
		n, err := reader.Read(buf)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			break
		}
		m.Write(buf[:n])
	}
	return hex.EncodeToString(m.Sum(nil))
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

package main

import (
	"fmt"
	"os"

	"github.com/toxyl/glog"
)

type FileReader struct {
	Name   string
	File   *os.File
	isOpen bool
	offset int
}

func (fr *FileReader) String() string {
	name := fr.Name
	if fr.isOpen {
		name += "*"
	}

	return fmt.Sprintf("%s:%s", glog.Auto(name), glog.PadLeft(glog.Auto(fr.offset), 4, ' '))
}

func (fr *FileReader) Close() {
	if fr.File != nil {
		fr.File.Close()
		fr.isOpen = false
	}
}

func (fr *FileReader) Open(filename string) error {
	if fr.isOpen {
		fr.Close()
	}
	fr.Name = filename
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	fr.File = file
	fr.isOpen = true
	return nil
}

// LookAhead read as many bytes as bCompare has from the FileReader.
// If all bytes match the method returns `true`. Otherwise it will
// reset the offset to the previous position and return `false`.
func (fr *FileReader) LookAhead(bCompare []byte) bool {
	pos := fr.offset
	bytes, err := fr.ReadBytes(len(bCompare))
	if err != nil {
		log.Error("%s Could not read bytes for look ahead: %s", fr.String(), glog.Error(err))
		fr.File.Seek(int64(pos), 0)
		fr.offset = pos
		return false
	}

	for i, b := range bCompare {
		if bytes[i] != b {
			fr.File.Seek(int64(pos), 0)
			fr.offset = pos
			return false
		}
	}

	return true
}

func (fr *FileReader) MoveCursor(n int) error {
	o, err := fr.File.Seek(int64(n), 1)
	fr.offset = int(o)
	return err
}

func (fr *FileReader) ReadBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	read, err := fr.File.Read(bytes)
	fr.offset += read
	return bytes, err
}

func (fr *FileReader) ReadInt() (int, error) {
	bytes := make([]byte, 1)
	read, err := fr.File.Read(bytes)
	fr.offset += read
	return int(bytes[0]), err
}

func (fr *FileReader) ReadString(n int) (string, error) {
	bytes, err := fr.ReadBytes(n)
	return string(bytes), err
}

func NewFileReader() *FileReader {
	fr := &FileReader{}
	return fr

}

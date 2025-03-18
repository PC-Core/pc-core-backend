package static

import "io"

type StaticFile struct {
	File        io.Reader
	Name        string
	ContentType string
}

func NewStaticFile(file io.Reader, name string, ctype string) *StaticFile {
	return &StaticFile{
		file, name, ctype,
	}
}

package utils

import (
	"fmt"
	"hash/crc32"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func GetFirmwareInfo(relativePath string) (error, string, uint32) {
	fileInfo := struct {
		path string
		info fs.FileInfo
	}{}

	for len([]rune(relativePath)) > 0 && relativePath[0] == '/' {
		relativePath = relativePath[1:]
	}

	err := filepath.Walk(relativePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		fileInfo.path = path
		fileInfo.info = info

		return nil
	})

	if err != nil || fileInfo.path == "" {
		return err, "", 0
	}
	return nil, fileInfo.path, uint32(fileInfo.info.Size())
}

// Calculate length and crc of file
func CalculateCrcFile(path string) uint32 {
	var crc uint32 = 0

	f, err := os.Open(path)

	defer func() {
		f.Close()
	}()

	if err != nil {
		return 0
	}

	b2 := make([]byte, 256)

	for true {
		n2, err := f.Read(b2)
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0
		}

		if n2 == 0 {
			break
		}
		crc = crc32.Update(crc, crc32.IEEETable, b2[:n2])
		if n2 < 256 {
			break
		}
	}
	return crc
}

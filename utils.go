package main

import (
	"hash/crc32"
	"io"
	"os"
)

// Calculate length and crc of file
func getCrcFile(path string) uint32 {
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

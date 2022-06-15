package main

import (
	"encoding/binary"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

type vcuModel string

var updateCache map[vcuModel]struct {
	path string
	len  uint32
	crc  uint32
} = map[vcuModel]struct {
	path string
	len  uint32
	crc  uint32
}{}

func vcuModelParse(b []byte) vcuModel {
	result := vcuModel("")

	if len(b) < 4 {
		return result
	}

	typeNum := binary.BigEndian.Uint16(b[1:3])
	rev := uint8(b[3])

	result += vcuModel(b[:1])
	result += vcuModel(fmt.Sprintf("%d-%d", typeNum, rev))
	return result
}

func (model vcuModel) getUpdateInfo() (fileLen uint32, fileCrc uint32) {
	fileLen = 0
	fileCrc = 0

	if data, isOK := updateCache[model]; isOK {
		fileLen, fileCrc = data.len, data.crc
	} else {
		path, length := model.getLastUpdate()
		fileLen = length
		fileCrc = getCrcFile(path)

		// update cache
		updateCache[model] = struct {
			path string
			len  uint32
			crc  uint32
		}{path, fileLen, fileCrc}
	}
	return
}

func (model vcuModel) getLastUpdate() (path string, len uint32) {
	lastUpdate := struct {
		at   time.Time
		path string
		info fs.FileInfo
	}{}

	err := filepath.Walk("files/"+string(model), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if t, err := time.Parse("2006-01-02", strings.Split(filepath.Base(path), ".")[0]); err == nil {
			if t.Year() != 1 && lastUpdate.at.Year() == 1 || t.After(lastUpdate.at) {
				lastUpdate.at = t
				lastUpdate.path = path
				lastUpdate.info = info
			}
		}

		return nil
	})

	if err != nil || lastUpdate.path == "" {
		return "", 0
	}
	return lastUpdate.path, uint32(lastUpdate.info.Size())
}

func (model vcuModel) getUpdatePath() string {
	path := ""

	if data, isOK := updateCache[model]; isOK {
		path = data.path
	} else {
		path, _ = model.getLastUpdate()
		model.getUpdateInfo()
	}

	return path
}

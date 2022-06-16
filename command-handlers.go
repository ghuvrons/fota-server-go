package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/ghuvrons/fota-server-go/models"

	giotgo "github.com/ghuvrons/g-IoT-Go"
	giot_packet "github.com/ghuvrons/g-IoT-Go/giot_packet"
	_ "github.com/go-sql-driver/mysql"
)

var downloaderBuffer map[*giotgo.ClientHandler][]byte = map[*giotgo.ClientHandler][]byte{}

func setCmdHandlers(server *giotgo.Server) {
	server.On(CMD_GET_INFO,
		func(client *giotgo.ClientHandler, data giot_packet.Data) (giot_packet.RespStatus, *bytes.Buffer) {
			b, isOK := data.(*bytes.Buffer)
			if !isOK {
				return giot_packet.RESP_UNKNOWN_ERROR, nil
			}

			typeNum := binary.BigEndian.Uint16(b.Bytes())

			// get firmware
			vf := models.VehicleModelGetLatestFirmware(client.Context(), uint32(typeNum), models.FIRMWARE_VCU)

			if vf == nil {
				return giot_packet.RESP_UNKNOWN_ERROR, nil
			}

			buf := &bytes.Buffer{}
			binary.Write(buf, binary.BigEndian, vf.FileLength)
			binary.Write(buf, binary.BigEndian, vf.Crc)

			return giot_packet.RESP_OK, buf
		},
	)

	server.On(CMD_DOWNLOAD,
		func(client *giotgo.ClientHandler, data giot_packet.Data) (giot_packet.RespStatus, *bytes.Buffer) {
			var readFileBuffer []byte
			b, isOK := data.(*bytes.Buffer)
			if !isOK {
				return giot_packet.RESP_UNKNOWN_ERROR, nil
			}

			readFileBuffer, isOK = downloaderBuffer[client]

			if !isOK {
				readFileBuffer = make([]byte, 1024)
				downloaderBuffer[client] = readFileBuffer
			}

			typeNum := binary.BigEndian.Uint16(b.Bytes())
			b.Next(4)
			offset := binary.BigEndian.Uint32(b.Bytes()[:4])
			b.Next(4)
			readLen := binary.BigEndian.Uint32(b.Bytes()[:4])
			b.Next(4)

			// get firmware
			vf := models.VehicleModelGetLatestFirmware(client.Context(), uint32(typeNum), models.FIRMWARE_VCU)
			if vf == nil {
				return giot_packet.RESP_UNKNOWN_ERROR, nil
			}

			f, err := os.Open(vf.BinaryPath)
			defer func() {
				f.Close()
			}()
			if err != nil {
				return giot_packet.RESP_UNKNOWN_ERROR, nil
			}

			if ret, err := f.Seek(int64(offset), 0); err != nil {
				fmt.Println(err)
				fmt.Println(offset, "=>", ret)
			} else {
			}

			if readLen > uint32(cap(readFileBuffer)) {
				readLen = uint32(cap(readFileBuffer))
			}
			n2, err := f.Read(readFileBuffer)

			if err != nil {
				fmt.Println(err)
				return giot_packet.RESP_UNKNOWN_ERROR, nil
			}

			buf := bytes.NewBuffer(readFileBuffer[:n2])

			return giot_packet.RESP_OK, buf
		},
	)
}

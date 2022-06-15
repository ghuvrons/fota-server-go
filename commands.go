package main

import (
	giot_packet "github.com/ghuvrons/g-IoT-Go/giot_packet"
)

const (
	CMD_GET_INFO giot_packet.Command = iota + 0xFF11
	CMD_DOWNLOAD
)

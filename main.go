package main

import (
	"VolumePrice/helper"
	"VolumePrice/vwap"
	"VolumePrice/ws"
)

func main() {
	ws := ws.WebSocket{}
	matches := make(chan helper.Match)
	ws.SubscribeAndRead(matches)
	vwap.GetVWap(matches)
}

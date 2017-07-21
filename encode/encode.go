package encode

import (
	"fmt"
	"strings"
	"time"
)

type EncodedMsg struct {
	Timestamp   string   `json:"timestamp"`
	MessageType string   `json:"type"`
	Network     string   `json:"network"`
	Data        []string `json:"data"`
	Id          string   `json:"id"`
}

type Encoder interface {
	Encode(in string, time time.Time) (out EncodedMsg)
}

func Encode(in string, t time.Time) *EncodedMsg {
	id := in[0:2]
	data := strings.Replace(in[2:], "-", "", -1)
	return &EncodedMsg{
		Timestamp:   fmt.Sprintf("%02d %v %d %02d:%02d:%02d +0000", t.Day(), t.Month().String()[:3], t.Year(), t.Hour(), t.Minute(), t.Second()),
		MessageType: "WirelessMessage",
		Network:     "Serial",
		Data:        []string{data},
		Id:          id,
	}
}

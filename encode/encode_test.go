package encode

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	enc := Encode("AATEMP025.8", time.Now())
	if enc.MessageType != "WirelessMessage" {
		t.Error("Expecting Wirelessmessage as type")
	}
	if enc.Network != "Serial" {
		t.Error("Expecting Serial as network")
	}
	if enc.Id != "AA" {
		t.Errorf("Expecting AA id, got %v", enc.Id)
	}
	if strings.Compare(enc.Data[0], "TEMP025.8") != 0 {
		t.Errorf("Expecting message TEMP025.8, got %v", enc.Data[0])
	}
}

func TestEncodeJson(t *testing.T) {
	tm := time.Date(2017, time.July, 05, 19, 37, 07, 0, time.UTC)
	jsonMsg := "{\"timestamp\":\"05 Jul 2017 19:37:07 +0000\",\"type\":\"WirelessMessage\",\"network\":\"Serial\",\"data\":[\"TEMP025.7\"],\"id\":\"AA\"}"
	encoded := Encode("AATEMP025.7", tm)

	mJson, err := json.Marshal(&encoded)
	if err != nil {
		t.Error("Error marshalling message")
	}

	if strings.Compare(jsonMsg, string(mJson)) != 0 {
		t.Error("Messages should be identical")
	}
}

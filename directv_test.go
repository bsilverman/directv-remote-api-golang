package directv

import (
	"testing"
	"time"
)

const ip = "127.0.0.1"
const port = 8080
const serialNum = "123456789"
const accessCardID = "0000-0000-0000"
const receiverID = "0000 0000 0000"
const softwareVersion = "0x994"
const stbVersion = "1.6"
const mode = 1

func TestSTBIsConnected(t *testing.T) {
	stb := &SetTopBox{ip, port}
	connected, _ := stb.IsConnected()
	if !connected {
		t.Error("Set Top Box is not connected.")
	}
}

func TestSTBGetSerialNum(t *testing.T) {
	stb := &SetTopBox{ip, port}
	num, err := stb.GetSerialNum()
	if err != nil {
		t.Error(err)
	}
	if serialNum != num {
		t.Error("Expected", serialNum, "got", num)
	}
}

func TestSTBGetSerialNumForClient(t *testing.T) {
	stb := &SetTopBox{ip, port}
	num, err := stb.GetSerialNumForClient(0)
	if err != nil {
		t.Error(err)
	}
	if serialNum != num {
		t.Error("Expected", serialNum, "got", num)
	}

	// Should Error
	num, err = stb.GetSerialNumForClient(1)
	if err == nil {
		t.Error("Expected error for ClientAddr 1 but got num", num)
	}
}

func TestSTBGetVersion(t *testing.T) {
	stb := &SetTopBox{ip, port}
	ver, err := stb.GetVersion()
	if err != nil {
		t.Error(err)
	}

	if accessCardID != ver.AccessCardID {
		t.Error("Expected", accessCardID, "got", ver.AccessCardID)
	}
	if receiverID != ver.ReceiverID {
		t.Error("Expected", receiverID, "got", ver.ReceiverID)
	}
	if softwareVersion != ver.STBSoftwareVersion {
		t.Error("Expected", softwareVersion, "got", ver.STBSoftwareVersion)
	}
	now := time.Now()
	timeStart := now.Add(-5 * time.Minute)
	timeEnd := now.Add(5 * time.Minute)
	if ver.SystemTime < timeStart.Unix() || ver.SystemTime > timeEnd.Unix() {
		t.Error("Expected time to be between ", timeStart.Unix(), "and", timeEnd.Unix(), "but was", ver.SystemTime)
	}
	if stbVersion != ver.Version {
		t.Error("Expected", stbVersion, "got", ver.Version)
	}
}

func TestSTBGetModeForClient(t *testing.T) {
	stb := &SetTopBox{ip, port}
	num, err := stb.GetModeForClient(0)
	if err != nil {
		t.Error(err)
	}
	if mode != num {
		t.Error("Expected", mode, "got", num)
	}

	// Should Error
	num, err = stb.GetModeForClient(1)
	if err == nil {
		t.Error("Expected error for ClientAddr 1 but got num", num)
	}
}

func TestSTBGetMode(t *testing.T) {
	stb := &SetTopBox{ip, port}
	num, err := stb.GetMode()
	if err != nil {
		t.Error(err)
	}
	if mode != num {
		t.Error("Expected", mode, "got", num)
	}
}

// func TestSTBProcessKey(t *testing.T) {
// 	stb := &SetTopBox{ip, port}
// 	err := stb.ProcessKey("4", "keyPress")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestSTBProcessCommand(t *testing.T) {
// 	stb := &SetTopBox{ip, port}
// 	_, err := stb.ProcessCommand("FA9A")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestSTBGetProgInfo(t *testing.T) {
// 	stb := &SetTopBox{ip, port}
// 	res, err := stb.GetProgInfoForTime(4, 65535, time.Now().Unix())
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if res.Title != "Let's Make a Deal" {
// 		t.Error("Expected", "Let's Make a Deal", "Got", res.Title)
// 	}
// }

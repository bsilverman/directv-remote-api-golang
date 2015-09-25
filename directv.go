package directv

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// SetTopBox is the primary object to use in the directv library.
type SetTopBox struct {
	IPAddress string
	Port      int
}

// Location represents a single Set Top Box location name.
type Location struct {
	ClientAddress string `json:"clientAddr"`
	LocationName  string `json:"locationName"`
}

// Version represents the version information for the Set Top Box
type Version struct {
	AccessCardID       string
	ReceiverID         string
	STBSoftwareVersion string
	SystemTime         int64
	Version            string
}

// CommandResults respresents the data returned from the Set Top Box for a
// ProcessCommand call.
type CommandResults struct {
	Command bool            `json:"command"`
	Param   bool            `json:"param"`
	Prefix  bool            `json:"prefix"`
	Return  CommandResponse `json:"return"`
	Status  statusResponse  `json:"status"`
}

// CommandResponse represents the response portion of the return value of
// ProcessCommand
type CommandResponse struct {
	Data     string `json:"data"`
	Response int    `json:"response"`
	Value    int    `json:"value"`
}

// ProgramStatusResponse represents the data returned from
type ProgramStatusResponse struct {
	CallSign    string `json:"callsign"`
	Date        string `json:"date"`
	Duration    int    `json:"duration"`
	IsOffAir    bool   `json:"isOffAir"`
	IsPClocked  int    `json:"isPclocked"`
	IsPPV       bool   `json:"isPpv"`
	IsRecording bool   `json:"IsRecording"`
	IsVOD       bool   `json:"IsVod"`
	Major       int    `json:"major"`
	Minor       int    `json:"minor"`
	ProgramID   string `json:"programId"`
	Rating      string `json:"rating"`
	StartTime   int64  `json:"startTime"`
	StationID   int64  `json:"stationId"`
	Title       string `json:"title"`
}

type statusResponse struct {
	Code          int    `json:"code"`
	CommandResult int    `json:"commandResult"`
	Message       string `json:"msg"`
	Query         string `json:"query"`
}

// Locations represents the getLocation response returned from the API.
type getLocationsResponse struct {
	Locations []Location     `json:"locations"`
	Status    statusResponse `json:"status"`
}

type getSerialNumResponse struct {
	SerialNum string         `json:"serialNum"`
	Status    statusResponse `json:"status"`
}

type getVersionResponse struct {
	AccessCardID       string `json:"accessCardId"`
	ReceiverID         string `json:"receiverId"`
	statusResponse     `json:"status"`
	STBSoftwareVersion string `json:"stbSoftwareVersion"`
	SystemTime         int64  `json:"systemTime"`
	Version            string `json:"version"`
}

type modeResponse struct {
	Mode   int            `json:"mode"`
	Status statusResponse `json:"status"`
}

// NewSetTopBox initialized a new SetTopBox struct with the supplied ip address
// and default port.
func NewSetTopBox(ip string) *SetTopBox {
	return &SetTopBox{ip, 8080}
}

// IsConnected returns true if the current SetTopBox object can talk to the DirecTV Set Top Box.
func (stb *SetTopBox) IsConnected() (bool, error) {
	locations, err := stb.GetLocations()
	if err != nil {
		return false, err
	}

	if len(locations) > 0 {
		return true, nil
	}
	return false, nil
}

// GetLocations calls /info/getLocations and returns the returned locations.
func (stb *SetTopBox) GetLocations() ([]Location, error) {
	var locationsRes getLocationsResponse
	_, err := stb.request("/info/getLocations", nil, &locationsRes)
	if err != nil {
		return nil, err
	}

	return locationsRes.Locations, nil
}

// GetSerialNum calls /info/getSerialNum and returns the STB Serial Number
func (stb *SetTopBox) GetSerialNum() (string, error) {
	var serialNumResponse getSerialNumResponse
	_, err := stb.request("/info/getSerialNum", nil, &serialNumResponse)
	if err != nil {
		return "", err
	}

	return serialNumResponse.SerialNum, nil
}

// GetSerialNumForClient calls /info/getSerialNum for a specific STB and returns
// the STB Serial Number
func (stb *SetTopBox) GetSerialNumForClient(clientAddr int) (string, error) {
	var serialNumResponse getSerialNumResponse
	params := map[string]string{"clientAddr": strconv.FormatInt(int64(clientAddr), 10)}
	_, err := stb.request("/info/getSerialNum", params, &serialNumResponse)
	if err != nil {
		return "", err
	}

	return serialNumResponse.SerialNum, nil
}

// GetVersion returns the version information, including time, from the SetTopBox
func (stb *SetTopBox) GetVersion() (Version, error) {
	var versionResponse getVersionResponse
	var version Version
	_, err := stb.request("/info/getVersion", nil, &versionResponse)
	if err != nil {

		return version, err
	}

	version = Version{
		versionResponse.AccessCardID,
		versionResponse.ReceiverID,
		versionResponse.STBSoftwareVersion,
		versionResponse.SystemTime,
		versionResponse.Version,
	}

	return version, nil
}

// GetMode calls /info/mode and returns the mode the STB is operating in.
func (stb *SetTopBox) GetMode() (int, error) {
	var modeResponse modeResponse
	_, err := stb.request("/info/mode", nil, &modeResponse)
	if err != nil {
		return 0, err
	}

	return modeResponse.Mode, nil
}

// GetModeForClient calls /info/mode and returns the mode the STB is operating in.
func (stb *SetTopBox) GetModeForClient(clientAddr int) (int, error) {
	var modeResponse modeResponse
	params := map[string]string{"clientAddr": strconv.FormatInt(int64(clientAddr), 10)}
	_, err := stb.request("/info/mode", params, &modeResponse)
	if err != nil {
		return 0, err
	}

	return modeResponse.Mode, nil
}

// ProcessKey sends a remote key press to the STB.  The hold parameter can be
// used to specify 'keyPress' (default), 'keyDown', or 'keyUp'  Remote keys include:
// format, power, rew, pause, play, stop, ffwd, replay, advance, record, guide, active, list, exit, up, down, select, left, right, back, menu, info, red, green, yellow, blue, chanup, chandown, prev, 1, 2, 3, 4, 5, 6, 7, 8, 9, dash, 0, enter
func (stb *SetTopBox) ProcessKey(key string, hold string) error {
	var status statusResponse
	params := map[string]string{
		"key":  key,
		"hold": hold,
	}
	_, err := stb.request("/remote/processKey", params, &status)
	return err
}

// ProcessKeyForClient sends a remote key press to the STB.  The hold parameter can be
// used to specify 'keyPress' (default), 'keyDown', or 'keyUp'  Remote keys include:
// format, power, rew, pause, play, stop, ffwd, replay, advance, record, guide, active, list, exit, up, down, select, left, right, back, menu, info, red, green, yellow, blue, chanup, chandown, prev, 1, 2, 3, 4, 5, 6, 7, 8, 9, dash, 0, enter
func (stb *SetTopBox) ProcessKeyForClient(key string, hold string, clientAddr int) error {
	var status statusResponse
	params := map[string]string{
		"clientAddr": strconv.FormatInt(int64(clientAddr), 10),
		"key":        key,
		"hold":       hold,
	}
	_, err := stb.request("/remote/processKey", params, &status)
	return err
}

// ProcessCommand sends a serial command (hex value) to the Set Top Box.
// Supported commands may include:
// 'FA81' Standby
// 'FA82' Active
// 'FA83' GetPrimaryStatus
// 'FA84' GetCommandVersion
// 'FA87' GetCurrentChannel
// 'FA90' GetSignalQuality
// 'FA91' GetCurrentTime
// 'FA92' GetUserCommand
// 'FA93' EnableUserEntry
// 'FA94' DisableUserEntry
// 'FA95' GetReturnValue
// 'FA96' Reboot
// 'FAA5' SendUserCommand
// 'FAA6' OpenUserChannel
// 'FA9A' GetTuner
// 'FA8A' GetPrimaryStatusMT
// 'FA8B' GetCurrentChannelMT
// 'FA9D' GetSignalQualityMT
// 'FA9F' OpenUserChannelMT
func (stb *SetTopBox) ProcessCommand(cmd string) (interface{}, error) {
	var response interface{}
	params := map[string]string{"cmd": cmd}
	_, err := stb.request("/serial/processCommand", params, &response)
	return response, err
}

// GetProgInfo returns information about the program on the specifed channel.
func (stb *SetTopBox) GetProgInfo(channel int) (ProgramStatusResponse, error) {
	var response ProgramStatusResponse
	params := map[string]string{"major": strconv.FormatInt(int64(channel), 10)}
	_, err := stb.request("/tv/getProgInfo", params, &response)
	fmt.Println(response)
	return response, err
}

// GetProgInfoForTime returns information about the program on the specifed channel.
func (stb *SetTopBox) GetProgInfoForTime(channelMajor int, channelMinor int, time int64) (ProgramStatusResponse, error) {
	var response ProgramStatusResponse
	params := map[string]string{
		"major": strconv.FormatInt(int64(channelMajor), 10),
		"minor": strconv.FormatInt(int64(channelMinor), 10),
		"time":  strconv.FormatInt(time, 10),
	}
	_, err := stb.request("/tv/getProgInfo", params, &response)
	fmt.Println(response)
	return response, err
}

// TuneToChannel tunes the SetTopBox to a specific channel.
func (stb *SetTopBox) TuneToChannel(channel int) error {
	var response interface{}
	params := map[string]string{"major": strconv.FormatInt(int64(channel), 10)}
	_, err := stb.request("/tv/tune", params, &response)
	fmt.Println(response)
	return err
}

func (stb *SetTopBox) request(uri string, params map[string]string, targetStruct interface{}) (*http.Response, error) {
	host := fmt.Sprintf("%s:%d", stb.IPAddress, stb.Port)
	requestURL := &url.URL{
		Scheme: "http",
		Host:   host,
		Path:   uri,
	}

	values := requestURL.Query()
	for key := range params {
		values.Add(key, params[key])
	}
	requestURL.RawQuery = values.Encode()

	fmt.Println(requestURL.String())
	res, err := http.Get(requestURL.String())

	if err != nil {
		return res, err
	}

	if res.StatusCode != 200 {
		return res, errors.New("Expected Status Code 200, got " + strconv.FormatInt(int64(res.StatusCode), 10))
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, targetStruct)

	return res, err
}

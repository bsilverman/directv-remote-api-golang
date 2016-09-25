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

// CommandResponse respresents the data returned from the Set Top Box for a
// ProcessCommand call.
type CommandResponse struct {
	Command bool          `json:"command"`
	Param   bool          `json:"param"`
	Prefix  bool          `json:"prefix"`
	Return  CommandReturn `json:"return"`
}

// CommandReturn represents the return portion of the return value of
// ProcessCommand
type CommandReturn struct {
	Data     string `json:"data"`
	Response int    `json:"response"`
	Value    int    `json:"value"`
}

// ProgramStatusResponse represents the data returned from
type ProgramStatusResponse struct {
	CallSign     string         `json:"callsign"`
	Date         string         `json:"date"`
	Duration     int            `json:"duration"`
	EpisodeTitle string         `json:"episodeTitle"`
	IsOffAir     bool           `json:"isOffAir"`
	IsPClocked   int            `json:"isPclocked"`
	IsPPV        bool           `json:"isPpv"`
	IsRecording  bool           `json:"isRecording"`
	IsVOD        bool           `json:"isVod"`
	Major        int            `json:"major"`
	Minor        int            `json:"minor"`
	Offset       int            `json:"offset"`
	ProgramID    string         `json:"programId"`
	Rating       string         `json:"rating"`
	StartTime    int64          `json:"startTime"`
	StationID    int64          `json:"stationId"`
	Status       statusResponse `json:"status"`
	Title        string         `json:"title"`
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
	AccessCardID       string         `json:"accessCardId"`
	ReceiverID         string         `json:"receiverId"`
	Status             statusResponse `json:"status"`
	STBSoftwareVersion string         `json:"stbSoftwareVersion"`
	SystemTime         int64          `json:"systemTime"`
	Version            string         `json:"version"`
}

type modeResponse struct {
	Mode   int            `json:"mode"`
	Status statusResponse `json:"status"`
}

type processKeyResponse struct {
	Hold   string         `json:"hold"`
	Key    string         `json:"key"`
	Status statusResponse `json:"status"`
}

type processCommandResponse struct {
	Command bool           `json:"command"`
	Param   bool           `json:"param"`
	Prefix  bool           `json:"prefix"`
	Return  CommandReturn  `json:"return"`
	Status  statusResponse `json:"status"`
}

type tuneResponse struct {
	Status statusResponse `json:"status"`
}

// Supported keys
const (
	KeyPower    = "power"
	KeyPowerOn  = "poweron"
	KeyPowerOff = "poweroff"
	KeyFormat   = "format"
	KeyPause    = "pause"
	KeyRewind   = "rew"
	KeyReplay   = "replay"
	KeyStop     = "stop"
	KeyAdvance  = "advance"
	KeyFFwd     = "ffwd"
	KeyRecord   = "record"
	KeyPlay     = "play"
	KeyGuide    = "guide"
	KeyActive   = "active"
	KeyList     = "list"
	KeyExit     = "exit"
	KeyBack     = "back"
	KeyMenu     = "menu"
	KeyInfo     = "info"
	KeyUp       = "up"
	KeyDown     = "down"
	KeyLeft     = "left"
	KeyRight    = "right"
	KeySelect   = "select"
	KeyRed      = "red"
	KeyGreen    = "green"
	KeyYellow   = "yellow"
	KeyBlue     = "blue"
	KeyChanup   = "chanup"
	KeyChandown = "chandown"
	KeyPrev     = "prev"
	Key0        = "0"
	Key1        = "1"
	Key2        = "2"
	Key3        = "3"
	Key4        = "4"
	Key5        = "5"
	Key6        = "6"
	Key7        = "7"
	Key8        = "8"
	Key9        = "9"
	KeyDash     = "dash"
	KeyEnter    = "enter"
)

// Supported key hold types
const (
	HoldPress           = "keyDown"
	HoldRelease         = "keyUp"
	HoldPressAndRelease = "keyPress"
)

// Supported commands
const (
	CommandStandby             = "FA81"
	CommandActive              = "FA82"
	CommandGetPrimaryStatus    = "FA83"
	CommandGetCommandVersion   = "FA84"
	CommandGetCurrentChannel   = "FA87"
	CommandGetSignalQuality    = "FA90"
	CommandGetCurrentTime      = "FA91"
	CommandGetUserCommand      = "FA92"
	CommandGetUserEntry        = "FA93"
	CommandDisableUserEntry    = "FA94"
	CommandGetReturnValue      = "FA95"
	CommandReboot              = "FA96"
	CommandSendUserCommand     = "FAA5"
	CommandOpenUserChannel     = "FAA6"
	CommandGetTuner            = "FA9A"
	CommandGetPrimaryStatusMT  = "FA8A"
	CommandGetCurrentChannelMT = "FA8B"
	CommandGetSignalQualityMT  = "FA9D"
	CommandOpenUserChannelMT   = "FA9F"
)

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
		if locationsRes.Status.Message != "" {
			err = errors.New(locationsRes.Status.Message)
		}
		return nil, err
	}

	return locationsRes.Locations, nil
}

// GetSerialNum calls /info/getSerialNum and returns the STB Serial Number
func (stb *SetTopBox) GetSerialNum(clientAddr string) (string, error) {
	var serialNumResponse getSerialNumResponse
	params := map[string]string{}
	if len(clientAddr) != 0 {
		params["clientAddr"] = clientAddr
	}
	_, err := stb.request("/info/getSerialNum", params, &serialNumResponse)
	if err != nil {
		if serialNumResponse.Status.Message != "" {
			err = errors.New(serialNumResponse.Status.Message)
		}
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
		if versionResponse.Status.Message != "" {
			err = errors.New(versionResponse.Status.Message)
		}
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
func (stb *SetTopBox) GetMode(clientAddr string) (int, error) {
	var modeResponse modeResponse
	params := map[string]string{}
	if len(clientAddr) != 0 {
		params["clientAddr"] = clientAddr
	}
	_, err := stb.request("/info/mode", params, &modeResponse)
	if err != nil {
		if modeResponse.Status.Message != "" {
			err = errors.New(modeResponse.Status.Message)
		}
		return 0, err
	}

	return modeResponse.Mode, nil
}

// ProcessKey sends a remote key press to the STB.
func (stb *SetTopBox) ProcessKey(key string, hold string, clientAddr string) error {
	var processKeyResponse processKeyResponse
	params := map[string]string{
		"key": key,
	}
	if len(hold) != 0 {
		params["hold"] = hold
	}
	if len(clientAddr) != 0 {
		params["clientAddr"] = clientAddr
	}
	_, err := stb.request("/remote/processKey", params, &processKeyResponse)
	if err != nil {
		if processKeyResponse.Status.Message != "" {
			err = errors.New(processKeyResponse.Status.Message)
		}
		return err
	}

	return nil
}

// ProcessCommand sends a serial command (hex value) to the Set Top Box.
func (stb *SetTopBox) ProcessCommand(cmd string) (CommandResponse, error) {
	var response processCommandResponse
	var commandResponse CommandResponse
	params := map[string]string{"cmd": cmd}
	_, err := stb.request("/serial/processCommand", params, &response)
	if err != nil {
		if response.Status.Message != "" {
			err = errors.New(response.Status.Message)
		}
		return commandResponse, err
	}

	commandResponse = CommandResponse{
		response.Command,
		response.Param,
		response.Prefix,
		response.Return,
	}

	return commandResponse, nil
}

// GetProgInfo returns information about the program on the specifed channel.
func (stb *SetTopBox) GetProgInfo(channelMajor int, channelMinor int, time int64, clientAddr string) (ProgramStatusResponse, error) {
	var response ProgramStatusResponse
	params := map[string]string{
		"major": strconv.FormatInt(int64(channelMajor), 10),
		"minor": strconv.FormatInt(int64(channelMinor), 10),
	}
	if time != 0 {
		params["time"] = strconv.FormatInt(time, 10)
	}
	if len(clientAddr) != 0 {
		params["clientAddr"] = clientAddr
	}
	_, err := stb.request("/tv/getProgInfo", params, &response)
	if err != nil {
		if response.Status.Message != "" {
			err = errors.New(response.Status.Message)
		}
		return response, err
	}

	return response, nil
}

// GetTuned returns information about the program a STB is tuned to.
func (stb *SetTopBox) GetTuned(clientAddr string) (ProgramStatusResponse, error) {
	var response ProgramStatusResponse
	params := map[string]string{}
	if len(clientAddr) != 0 {
		params["clientAddr"] = clientAddr
	}
	_, err := stb.request("/tv/getTuned", params, &response)
	if err != nil {
		if response.Status.Message != "" {
			err = errors.New(response.Status.Message)
		}
		return response, err
	}

	return response, nil
}

// TuneToChannel tunes the SetTopBox to a specific channel.
func (stb *SetTopBox) TuneToChannel(channel string, clientAddr string) error {
	var response tuneResponse
	params := map[string]string{"major": channel}
	if len(clientAddr) != 0 {
		params["clientAddr"] = clientAddr
	}
	_, err := stb.request("/tv/tune", params, &response)
	if err != nil {
		if response.Status.Message != "" {
			err = errors.New(response.Status.Message)
		}
		return err
	}

	return nil
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
		err = errors.New("Expected Status Code 200, got " + strconv.FormatInt(int64(res.StatusCode), 10))
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, targetStruct)

	return res, err
}

package xlink

import "regexp"

type XLink struct {
	Id   Id `json:"sysid"`
	Data struct {
		Local Local `json:"local"`
	} `json:"data"`
}

func (xLink *XLink) GetId() Id {
	return xLink.Id
}

func (xLink *XLink) GetName() string {
	return xLink.Data.Local.Name
}

func (xlink *XLink) GetEncoders() []Encoder {
	return xlink.Data.Local.Enc
}

func (xlink *XLink) GetDecoders() []Decoder {
	return xlink.Data.Local.Dec
}

type Id string

type EnDecoderId string

type EnDecoderIdType int

const (
	TypeUnknown EnDecoderIdType = 0
	TypeEncoder EnDecoderIdType = 1
	TypeDecoder EnDecoderIdType = 2
)

var regTypeEncoder = regexp.MustCompile(`^X\d[a-zA-Z]\d+-E\d+$`)
var regTypeDecoder = regexp.MustCompile(`^X\d[a-zA-Z]\d+-D\d+$`)

func (edid EnDecoderId) Type() EnDecoderIdType {
	if regTypeEncoder.MatchString(string(edid)) {
		return TypeEncoder
	}
	if regTypeDecoder.MatchString(string(edid)) {
		return TypeDecoder
	}
	return TypeUnknown
}

type Local struct {
	Name    string    `json:"name"`
	Enc     []Encoder `json:"enc"`
	Dec     []Decoder `json:"dec"`
	Network struct {
		Nets []Ethernet `json:"nets"`
	} `json:"network"`
}

type EnDecoder interface {
	GetName() (string, bool)
	PhyicalNumber() (int, bool)
	HasVideoSignal() (bool, bool)
	HasAudioSignal() (bool, bool)
	IsRunning() (bool, bool)
	IsConnected() (bool, bool)
}

type EnDecoderStateUpdateValues interface {
	IsVideoEnabled() (bool, bool)
	IsAudioEnabled() (bool, bool)
}

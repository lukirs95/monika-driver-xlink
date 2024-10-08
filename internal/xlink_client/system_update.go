package xlinkclient

import (
	"fmt"

	"github.com/lukirs95/monika-driver-xlink/internal/xlink"
	"github.com/lukirs95/monika-gosdk/pkg/types"
)

func ParseSystemUpdate(stateMap StateMap, xlink *xlink.XLink) {
	for _, encoder := range xlink.GetEncoders() {
		ParseEncoderUpdate(stateMap, &encoder)
	}

	for _, decoder := range xlink.GetDecoders() {
		ParseDecoderUpdate(stateMap, &decoder)
	}
}

func ParseEncoderUpdate(stateMap StateMap, encoder *xlink.Encoder) {
	state, found := stateMap.GetStateOf(encoder.GetId())
	if !found {
		return
	}

	module := state.GetModule()
	videoIn := state.video
	audioIn := state.audio
	videoOut := module.GetIOlet(types.IOletId(fmt.Sprintf("%s:vo", encoder.GetId())))
	audioOut := module.GetIOlet(types.IOletId(fmt.Sprintf("%s:ao", encoder.GetId())))

	if encoder.GetReceiver() != nil {
		if _, found := stateMap.GetStateOf(encoder.GetReceiver().GetId()); !found {
			stateMap.Add(encoder.GetReceiver().GetId(), NewStateEntry("", module, videoOut, audioOut))
		}
	}

	if value, ok := encoder.GetName(); ok {
		module.SetName(value)
	}

	videoInStatus := videoIn.GetStatus()
	audioInStatus := audioIn.GetStatus()

	if value, ok := encoder.HasVideoSignal(); ok {
		videoInStatus.SetReceiving(value)
	}

	if value, ok := encoder.HasAudioSignal(); ok {
		audioInStatus.SetReceiving(value)
	}

	if value, ok := encoder.IsRunning(); ok {
		videoInStatus.SetSending(value)
		audioInStatus.SetSending(value)
	}

	if value, ok := encoder.IsVideoEnabled(); ok {
		videoInStatus.SetEnabled(value)
	}

	if value, ok := encoder.IsAudioEnabled(); ok {
		audioInStatus.SetEnabled(value)
	}

	videoIn.SetStatus(videoInStatus)
	audioIn.SetStatus(audioInStatus)

	videoOutStatus := videoOut.GetStatus()
	audioOutStatus := audioOut.GetStatus()

	receiver := encoder.GetReceiver()

	if receiver != nil {
		connected := false
		if value, ok := receiver.IsConnected(); ok {
			connected = value
			videoOutStatus.SetOK(value)
			audioOutStatus.SetOK(value)
			if !value {
				videoOutStatus.SetReceiving(false)
				videoOutStatus.SetSending(false)
				videoOutStatus.SetEnabled(false)
				audioOutStatus.SetReceiving(false)
				audioOutStatus.SetSending(false)
				audioOutStatus.SetEnabled(false)
			}
		}

		if value, ok := receiver.HasVideoSignal(); ok && connected {
			videoOutStatus.SetSending(value)
		}

		if value, ok := receiver.HasAudioSignal(); ok && connected {
			audioOutStatus.SetSending(value)
		}

		if value, ok := receiver.IsRunning(); ok && connected {
			videoOutStatus.SetReceiving(value)
			audioOutStatus.SetReceiving(value)
		}

		if value, ok := receiver.IsVideoEnabled(); ok && connected {
			videoOutStatus.SetEnabled(value)
		}

		if value, ok := receiver.IsAudioEnabled(); ok && connected {
			audioOutStatus.SetEnabled(value)
		}
	}

	videoOut.SetStatus(videoOutStatus)
	audioOut.SetStatus(audioOutStatus)

	moduleStatus := module.GetStatus()

	moduleStatus.SetOK(videoInStatus.OK() && audioInStatus.OK() && videoOutStatus.OK() && audioOutStatus.OK())

	module.SetStatus(moduleStatus)
}

func ParseDecoderUpdate(stateMap StateMap, decoder *xlink.Decoder) {
	state, found := stateMap.GetStateOf(decoder.GetId())
	if !found {
		return
	}

	module := state.GetModule()
	videoOut := state.video
	audioOut := state.audio
	videoIn := module.GetIOlet(types.IOletId(fmt.Sprintf("%s:vi", decoder.GetId())))
	audioIn := module.GetIOlet(types.IOletId(fmt.Sprintf("%s:ai", decoder.GetId())))

	if decoder.GetSender() != nil {
		if _, found := stateMap.GetStateOf(decoder.GetSender().GetId()); !found {
			stateMap.Add(decoder.GetSender().GetId(), NewStateEntry("", module, videoIn, audioIn))
		}
	}

	if value, ok := decoder.GetName(); ok {
		module.SetName(value)
	}

	videoOutStatus := videoOut.GetStatus()
	audioOutStatus := audioOut.GetStatus()

	if value, ok := decoder.HasVideoSignal(); ok {
		videoOutStatus.SetReceiving(value)
	}

	if value, ok := decoder.HasAudioSignal(); ok {
		audioOutStatus.SetReceiving(value)
	}

	if value, ok := decoder.IsRunning(); ok {
		videoOutStatus.SetSending(value)
		audioOutStatus.SetSending(value)
	}

	if value, ok := decoder.IsVideoEnabled(); ok {
		videoOutStatus.SetEnabled(value)
	}

	if value, ok := decoder.IsAudioEnabled(); ok {
		audioOutStatus.SetEnabled(value)
	}

	videoOut.SetStatus(videoOutStatus)
	audioOut.SetStatus(audioOutStatus)

	videoInStatus := videoIn.GetStatus()
	audioInStatus := audioIn.GetStatus()

	sender := decoder.GetSender()

	if sender == nil {
		videoInStatus.SetOK(false)
		videoInStatus.SetReceiving(false)
		videoInStatus.SetSending(false)
		videoInStatus.SetEnabled(false)
		audioInStatus.SetOK(false)
		audioInStatus.SetReceiving(false)
		audioInStatus.SetSending(false)
		audioInStatus.SetEnabled(false)
	} else {
		connected := false
		if value, ok := sender.IsConnected(); ok {
			connected = value
			videoInStatus.SetOK(value)
			audioInStatus.SetOK(value)
			if !value {
				videoInStatus.SetReceiving(false)
				videoInStatus.SetSending(false)
				videoInStatus.SetEnabled(false)
				audioInStatus.SetReceiving(false)
				audioInStatus.SetSending(false)
				audioInStatus.SetEnabled(false)
			}
		}

		if value, ok := sender.HasVideoSignal(); ok && connected {
			videoInStatus.SetSending(value)
		}

		if value, ok := sender.HasAudioSignal(); ok && connected {
			audioInStatus.SetSending(value)
		}

		if value, ok := sender.IsRunning(); ok && connected {
			videoInStatus.SetReceiving(value)
			audioInStatus.SetReceiving(value)
		}

		if value, ok := sender.IsVideoEnabled(); ok && connected {
			videoInStatus.SetEnabled(value)
		}

		if value, ok := sender.IsAudioEnabled(); ok && connected {
			audioInStatus.SetEnabled(value)
		}
	}

	videoIn.SetStatus(videoInStatus)
	audioIn.SetStatus(audioInStatus)

	moduleStatus := module.GetStatus()

	moduleStatus.SetOK(videoInStatus.OK() && audioInStatus.OK() && videoOutStatus.OK() && audioOutStatus.OK())

	module.SetStatus(moduleStatus)
}

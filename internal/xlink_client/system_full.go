package xlinkclient

import (
	"fmt"

	"github.com/lukirs95/monika-driver-xlink/internal/xlink"
	"github.com/lukirs95/monika-gosdk/pkg/types"
)

func (client *Client) ParseSystem(state StateMap, xlink *xlink.XLink) {
	client.device.SetId(types.DeviceId(xlink.GetId()))
	for _, encoder := range xlink.GetEncoders() {
		if _, found := state.GetStateOf(encoder.Id); !found {
			newModule := types.NewModule(types.ModuleId(encoder.GetId()), types.ModuleType_AV, encoder.Name)
			videoInput := types.NewIOlet(types.IOletId(fmt.Sprintf("%s:vi", encoder.GetId())), getVideoTypeIN(&encoder), fmt.Sprintf("%s:vi", encoder.GetId()))
			audioInput := types.NewIOlet(types.IOletId(fmt.Sprintf("%s:ai", encoder.GetId())), getAudioTypeIN(&encoder), fmt.Sprintf("%s:ai", encoder.GetId()))
			videoOutput := types.NewIOlet(types.IOletId(fmt.Sprintf("%s:vo", encoder.GetId())), getVideoTypeOUT(encoder.GetReceiver()), fmt.Sprintf("%s:vo", encoder.GetReceiver().GetId()))
			audioOutput := types.NewIOlet(types.IOletId(fmt.Sprintf("%s:ao", encoder.GetId())), getAudioTypeOUT(encoder.GetReceiver()), fmt.Sprintf("%s:ao", encoder.GetReceiver().GetId()))
			state.Add(encoder.GetId(), NewStateEntry(xlink.Id, newModule, videoInput, audioInput))
			state.Add(encoder.GetReceiver().GetId(), NewStateEntry(xlink.Id, newModule, videoOutput, audioOutput))
			newModule.AddIOlet(videoInput)
			newModule.AddIOlet(audioInput)
			newModule.AddIOlet(videoOutput)
			newModule.AddIOlet(audioOutput)
			newModule.AddAction(types.ModuleControl_START, client.startModule)
			newModule.AddAction(types.ModuleControl_STOP, client.stopModule)
			newModule.AddAction(types.ModuleControl_RESTART, client.restartModule)
			client.device.AddModule(newModule)
		}

	}
	for _, decoder := range xlink.GetDecoders() {
		if _, found := state.GetStateOf(decoder.Id); !found {
			newModule := types.NewModule(types.ModuleId(decoder.GetId()), types.ModuleType_AV, decoder.Name)
			videoInput := types.NewIOlet(types.IOletId(fmt.Sprintf("%s:vi", decoder.GetId())), getVideoTypeIN(decoder.GetSender()), fmt.Sprintf("%s:vi", decoder.GetSender().GetId()))
			audioInput := types.NewIOlet(types.IOletId(fmt.Sprintf("%s:ai", decoder.GetId())), getAudioTypeIN(decoder.GetSender()), fmt.Sprintf("%s:ai", decoder.GetSender().GetId()))
			videoOutput := types.NewIOlet(types.IOletId(fmt.Sprintf("%s:vo", decoder.GetId())), getVideoTypeOUT(&decoder), fmt.Sprintf("%s:vo", decoder.GetId()))
			audioOutput := types.NewIOlet(types.IOletId(fmt.Sprintf("%s:ao", decoder.GetId())), getAudioTypeOUT(&decoder), fmt.Sprintf("%s:ao", decoder.GetId()))
			state.Add(decoder.GetId(), NewStateEntry(xlink.Id, newModule, videoOutput, audioOutput))
			state.Add(decoder.GetSender().GetId(), NewStateEntry(xlink.Id, newModule, videoInput, audioInput))
			newModule.AddIOlet(videoInput)
			newModule.AddIOlet(videoOutput)
			newModule.AddIOlet(audioInput)
			newModule.AddIOlet(audioOutput)
			newModule.AddAction(types.ModuleControl_START, client.startModule)
			newModule.AddAction(types.ModuleControl_STOP, client.stopModule)
			newModule.AddAction(types.ModuleControl_RESTART, client.restartModule)
			client.device.AddModule(newModule)
		}
	}
}

func getVideoTypeIN(encoder xlink.EnDecoder) types.IOletType {
	if card, ok := encoder.PhyicalNumber(); ok && card == 12 {
		return types.IOletType_IPVIDEOIN
	} else {
		return types.IOletType_BBVIDEOIN
	}
}

func getAudioTypeIN(encoder xlink.EnDecoder) types.IOletType {
	if card, ok := encoder.PhyicalNumber(); ok && card == 12 {
		return types.IOletType_IPAUDIOIN
	} else {
		return types.IOletType_BBAUDIOIN
	}
}

func getVideoTypeOUT(encoder xlink.EnDecoder) types.IOletType {
	if card, ok := encoder.PhyicalNumber(); ok && card == 12 {
		return types.IOletType_IPVIDEOOUT
	} else {
		return types.IOletType_BBVIDEOOUT
	}
}

func getAudioTypeOUT(encoder xlink.EnDecoder) types.IOletType {
	if card, ok := encoder.PhyicalNumber(); ok && card == 12 {
		return types.IOletType_IPAUDIOOUT
	} else {
		return types.IOletType_BBAUDIOOUT
	}
}

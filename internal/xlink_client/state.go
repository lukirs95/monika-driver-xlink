package xlinkclient

import (
	"github.com/lukirs95/monika-driver-xlink/internal/xlink"
	"github.com/lukirs95/monika-gosdk/pkg/types"
)

type StateMap map[xlink.EnDecoderId]*StateEntry

func NewStateMap() StateMap {
	return make(StateMap)
}

func (stateMap StateMap) Add(endecoder xlink.EnDecoderId, newEntry *StateEntry) (created bool) {
	if oldState, ok := stateMap[endecoder]; !ok {
		stateMap[endecoder] = newEntry
		created = true
	} else {
		if string(oldState.video.GetId()) != oldState.video.GetName() {
			oldState.video.SetName("-")
		}
		if string(oldState.audio.GetId()) != oldState.audio.GetName() {
			oldState.audio.SetName("-")
		}
		stateMap[endecoder] = newEntry
		created = false
	}
	return
}

func (stateMap StateMap) GetStateOf(endecoder xlink.EnDecoderId) (*StateEntry, bool) {
	stateEntry, ok := stateMap[endecoder]
	return stateEntry, ok
}

type StateEntry struct {
	xlink  xlink.Id
	module types.Module
	video  types.IOlet
	audio  types.IOlet
}

func NewStateEntry(xlink xlink.Id, module types.Module, video types.IOlet, audio types.IOlet) *StateEntry {
	return &StateEntry{
		xlink:  xlink,
		module: module,
		video:  video,
		audio:  audio,
	}
}

func (stateEntry *StateEntry) GetXLink() xlink.Id {
	return stateEntry.xlink
}

func (stateEntry *StateEntry) GetModule() types.Module {
	return stateEntry.module
}

func (stateEntry *StateEntry) GetVideoIOlet() types.IOlet {
	return stateEntry.video
}

func (stateEntry *StateEntry) GetAudioIOlet() types.IOlet {
	return stateEntry.audio
}

package schedule

import (
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/testutil"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

type fakeMessageSender struct {
	messages      []proto.Message
	messagesMutex sync.RWMutex
}

func (fms *fakeMessageSender) SendMessage(message proto.Message) error {
	fms.messagesMutex.Lock()
	defer fms.messagesMutex.Unlock()
	if fms.messages == nil {
		fms.messages = []proto.Message{message}
	} else {
		fms.messages = append(fms.messages, message)
	}
	return nil
}

func (fms *fakeMessageSender) Messages() []proto.Message {
	fms.messagesMutex.RLock()
	defer fms.messagesMutex.RUnlock()
	return testutil.CloneMessages(fms.messages)
}

func assertFakeMessageSenderMessages(t *testing.T, fms *fakeMessageSender, expected []proto.Message, name string) {
	actual := fms.Messages()
	assert.Equal(t, expected, actual, "serverState %s send the wrong messages", name)
}

func TestServerState_sendVolume(t *testing.T) {
	for i := 0; i < 16; i++ {
		fms := new(fakeMessageSender)
		ss := &serverState{volume: float64(i)}
		ss.sendVolume(fms)

		assertFakeMessageSenderMessages(t, fms, []proto.Message{&comm.SetVolumeRequest{Volume: float64(i)}}, "sendVolume")
	}
}

func TestServerState_sendNewestSong(t *testing.T) {
	for _, sm := range testSongMessages {
		ss := &serverState{newestSong: sm}
		fms := new(fakeMessageSender)
		ss.sendNewestSong(fms)

		assertFakeMessageSenderMessages(t, fms, []proto.Message{sm}, "sendNewestSong")
	}

	ss := &serverState{}
	fms := new(fakeMessageSender)
	ss.sendNewestSong(fms)
	messages := fms.Messages()
	assert.Zero(t, len(messages), "serverState sendNewest song send a message without newest song being set: %v", messages)
}

func TestServerState_sendPauses(t *testing.T) {
	for _, pms := range testPauseMessages {
		ss := &serverState{pauses: pms}
		fms := new(fakeMessageSender)
		ss.sendPauses(fms)

		expected := make([]proto.Message, len(pms))
		for i, pm := range pms {
			expected[i] = pm
		}

		assertFakeMessageSenderMessages(t, fms, expected, "sendPauses")
	}
}

func TestServerState_removablePauses(t *testing.T) {
	for _, c := range removePauseCases {
		ss := &serverState{pauses: c.pauses, newestSong: &comm.NewSongInfo{FirstSampleOfSongIndex: c.songStartSample}}
		assert.Equal(t, c.removable, ss.removablePauses(), "serverState removablePauses returned incorrect number of removable pauses for case %v", c)
	}
}

func TestServerState_removeOldPauses(t *testing.T) {
	for _, c := range removePauseCases {
		ss := &serverState{pauses: c.pauses, newestSong: &comm.NewSongInfo{FirstSampleOfSongIndex: c.songStartSample}}
		ss.removeOldPauses()
		assert.Equal(t, c.result, ss.pauses, "serverState removeOldPauses resulted with the wrong pauses for case %v", c)
	}
}

func TestToWireLyrics(t *testing.T) {
	for _, c := range toWireLyricsCases {
		actual := toWireLyrics(c.lyrics)
		assert.Equal(t, c.result, actual, "toWireLyrics returned the wrong wire lyrics")
	}
}

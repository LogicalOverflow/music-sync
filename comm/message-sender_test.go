package comm

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestSingleMessageSender_SendMessage(t *testing.T) {
	for _, p := range testPackages {
		expectedBytes, err := toWire(p)
		if assert.Nil(t, err, "toWire returned an error for package %v: %v", p, err) {
			conn := newBufferConn()
			sms := singleMessageSender{connection: conn}
			sms.SendMessage(p)
			actualBytes := conn.Bytes()
			assert.True(t, bytes.Equal(expectedBytes, actualBytes), "singleMessageSender sendMessage did not write toWire to the connection for package %v", p)
		}
	}
}

func TestMultiMessageSender_SendMessage(t *testing.T) {
	for i, p := range testPackages {
		expectedBytes, err := toWire(p)
		if assert.Nil(t, err, "toWire returned an error for package %v: %v", p, err) {
			channels := testPackageChannels[i]

			connNoChan := newBufferConn()
			connNotInDict := newBufferConn()
			connAudio := newBufferConn()
			connMeta := newBufferConn()
			connBoth := newBufferConn()

			mms := &multiMessageSender{
				connections: []net.Conn{connNoChan, connNotInDict, connAudio, connMeta, connBoth},
				channels: map[net.Conn][]Channel{
					connNoChan: {},
					connAudio:  {Channel_AUDIO},
					connMeta:   {Channel_META},
					connBoth:   {Channel_AUDIO, Channel_META},
				}}
			mms.SendMessage(p)

			hasAudio := false
			hasMeta := false
			if len(channels) == 0 {
				hasAudio = true
				hasMeta = true
				assert.True(t, bytes.Equal(expectedBytes, connNoChan.Bytes()), "multiMessageSender sendMessage did not write toWire to the connection without channels for package %v", p)
				assert.True(t, bytes.Equal(expectedBytes, connNotInDict.Bytes()), "multiMessageSender sendMessage did not write toWire to the connection not in the channels dictionary for package %v", p)
			} else {
				assert.Zero(t, len(connNoChan.Bytes()), "multiMessageSender sendMessage did write to the connection without channels for package %v", p)
				assert.Zero(t, len(connNotInDict.Bytes()), "multiMessageSender sendMessage did write to the connection not in the channels dictionary for package %v", p)
			}
			for _, c := range channels {
				if c == Channel_AUDIO {
					hasAudio = true
				}
				if c == Channel_META {
					hasMeta = true
				}
			}
			if hasAudio {
				assert.True(t, bytes.Equal(expectedBytes, connAudio.Bytes()), "multiMessageSender sendMessage did not write toWire to the connection with the AUDIO channel for package %v", p)
			} else {
				assert.Zero(t, len(connAudio.Bytes()), "multiMessageSender sendMessage did write to the connection with the AUDIO channel for package %v", p)
			}
			if hasMeta {
				assert.True(t, bytes.Equal(expectedBytes, connMeta.Bytes()), "multiMessageSender sendMessage did not write toWire to the connection with the META channel for package %v", p)
			} else {

				assert.Zero(t, len(connMeta.Bytes()), "multiMessageSender sendMessage did write to the connection with the META channel for package %v", p)
			}

			assert.True(t, bytes.Equal(expectedBytes, connBoth.Bytes()), "multiMessageSender sendMessage did not write toWire to the connection with both channels for package %v", p)
		}
	}
}

func TestMultiMessageSender_AddConn(t *testing.T) {
	mms := &multiMessageSender{connections: make([]net.Conn, 0), channels: make(map[net.Conn][]Channel, 0)}
	conns := []net.Conn{newBufferConn(), newBufferConn(), newBufferConn()}

	assert.Zero(t, len(mms.connections), "multiMessageSender contains elements before adding any connections")

	for i, c := range conns {
		mms.AddConn(c)
		assert.Equal(t, i+1, len(mms.connections), "multiMessageSender has the wrong number of connections after adding %d connections", i+i)
		exists := false
		for _, e := range mms.connections {
			if e == c {
				exists = true
				break
			}
		}
		assert.True(t, exists, "multiMessageSender connections slice does not contain the added connection after adding %d connections", i+1)
	}
}

func TestMultiMessageSender_DelConn(t *testing.T) {
	conns := []net.Conn{newBufferConn(), newBufferConn(), newBufferConn()}
	mms := &multiMessageSender{connections: make([]net.Conn, len(conns)), channels: make(map[net.Conn][]Channel, 0)}
	copy(mms.connections, conns)

	for i, c := range conns {
		mms.DelConn(c)
		assert.Equal(t, len(conns)-i-1, len(mms.connections), "multiMessageSender has the wrong number of connections after deleting %d connections", i+i)
		exists := false
		for _, e := range mms.connections {
			if e == c {
				exists = true
				break
			}
		}
		assert.False(t, exists, "multiMessageSender connections slice contains the connection after deleting %d connections", i+1)
	}

	assert.Zero(t, len(mms.connections), "multiMessageSender contains elements after deleting all connections")
}

func TestMultiMessageSender_Subscribe(t *testing.T) {
	conn := newBufferConn()
	mms := &multiMessageSender{connections: []net.Conn{conn}, channels: make(map[net.Conn][]Channel, 0)}

	assert.False(t, mms.isSubscribed(conn, []Channel{Channel_AUDIO}), "multiMessageSender claims the connection is subscribed to the AUDIO channel")
	assert.False(t, mms.isSubscribed(conn, []Channel{Channel_META}), "multiMessageSender claims the connection is subscribed to the META channel")

	mms.Subscribe(conn, Channel_AUDIO)
	assert.True(t, mms.isSubscribed(conn, []Channel{Channel_AUDIO}), "multiMessageSender claims the connection is not subscribed to the AUDIO channel")
	assert.False(t, mms.isSubscribed(conn, []Channel{Channel_META}), "multiMessageSender claims the connection is subscribed to the META channel")

	mms.Subscribe(conn, Channel_META)
	assert.True(t, mms.isSubscribed(conn, []Channel{Channel_AUDIO}), "multiMessageSender claims the connection is not subscribed to the AUDIO channel")
	assert.True(t, mms.isSubscribed(conn, []Channel{Channel_META}), "multiMessageSender claims the connection is not subscribed to the META channel")
}

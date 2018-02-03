package comm

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"io/ioutil"
	"net"
	"reflect"
)

const zlibLevel = 5

func channelOf(m proto.Message) ([]Channel, bool) {
	switch m.(type) {
	case *QueueChunkRequest:
		return []Channel{Channel_AUDIO}, true
	case *SetVolumeRequest:
		return []Channel{Channel_AUDIO, Channel_META}, true
	case *ChunkInfo, *NewSongInfo, *PauseInfo:
		return []Channel{Channel_META}, true
	default:
		return []Channel{}, false
	}
}

func toWire(m proto.Message) ([]byte, error) {
	inner, err := proto.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal inner protobuf (%s): %v", proto.MessageName(m), err)
	}

	compressed := bytes.NewBuffer([]byte{})
	w, _ := zlib.NewWriterLevel(compressed, zlibLevel)
	io.Copy(w, bytes.NewReader(inner))
	w.Close()

	envelope := &Envelope{
		Type: proto.MessageName(m),
		Data: compressed.Bytes(),
	}

	env, err := proto.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal envelope protobuf: %v", err)
	}

	wire := make([]byte, len(env)+8)
	binary.LittleEndian.PutUint64(wire[:8], uint64(len(env)))
	copy(wire[8:], env)

	return wire, nil
}

func readWire(conn net.Conn) (proto.Message, error) {
	sizeBuf, err := readBytes(conn, 8)
	if err != nil {
		return nil, err
	}
	size := int64(binary.LittleEndian.Uint64(sizeBuf))
	env, err := readBytes(conn, size)

	envelope := &Envelope{}
	if err := proto.Unmarshal(env, envelope); err != nil {
		return nil, fmt.Errorf("failed to unmarshal envelope protobuf: %v", err)
	}

	mType := proto.MessageType(envelope.Type)
	if mType == nil {
		return nil, fmt.Errorf("unknown inner type: %s", envelope.Type)
	}
	m := reflect.New(mType.Elem()).Interface().(proto.Message)

	r, _ := zlib.NewReader(bytes.NewBuffer(envelope.Data))
	data, _ := ioutil.ReadAll(r)
	r.Close()

	if err := proto.Unmarshal(data, m); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inner protobuf (%s): %v", envelope.Type, err)
	}

	return m, nil
}

func sendWire(m proto.Message, conn net.Conn) error {
	data, err := toWire(m)
	if err != nil {
		return err
	}
	send := 0
	for send < len(data) {
		n, err := conn.Write(data[send:])
		if err != nil {
			return err
		}
		send += n
	}
	return nil
}

func readBytes(conn net.Conn, number int64) ([]byte, error) {
	buf := make([]byte, number)
	read := 0
	for read < len(buf) {
		n, err := conn.Read(buf[read:])
		if err != nil {
			return nil, err
		}
		read += n
	}
	return buf, nil
}

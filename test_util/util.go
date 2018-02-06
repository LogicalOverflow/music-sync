package test_util

import "github.com/golang/protobuf/proto"

// CloneMessages clones the message slice
func CloneMessages(m []proto.Message) []proto.Message {
	if m == nil || len(m) == 0 {
		return []proto.Message{}
	}
	r := make([]proto.Message, len(m))
	copy(r, m)
	return r
}

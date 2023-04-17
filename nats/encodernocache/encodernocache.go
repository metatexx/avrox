package encodernocache

import (
	"errors"

	"github.com/metatexx/avrox"
	"github.com/nats-io/nats.go"
)

const (
	AVROX_ENCODER = "avrox"
)

func init() {
	// Register AvroX encoder
	ae := &AvroXEncoder{}
	nats.RegisterEncoder(AVROX_ENCODER, ae)
}

// AvroXEncoder is a AvroX implementation for EncodedConn
// We use an alternative cache to the default one though
type AvroXEncoder struct {
	Compression avrox.CompressionID
}

var (
	ErrInvalidAvroXMsgEncode = errors.New("nats: Invalid avrox object passed to encode")
	ErrInvalidAvroXMsgDecode = errors.New("nats: Invalid avrox object passed to decode")
)

// Encode
func (pb *AvroXEncoder) Encode(_ string, v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	i, found := v.(avrox.Schemer)
	if !found {
		return nil, ErrInvalidAvroXMsgEncode
	}
	data, err := avrox.Marshal(i, pb.Compression, nil)
	if err != nil {
		return nil, errors.Join(ErrInvalidAvroXMsgEncode, err)
	}
	return data, nil
}

// Decode
func (pb *AvroXEncoder) Decode(_ string, data []byte, vPtr interface{}) error {
	if _, ok := vPtr.(*interface{}); ok {
		return nil
	}
	i, found := vPtr.(avrox.Schemer)
	if !found {
		return ErrInvalidAvroXMsgDecode
	}
	return avrox.Unmarshal(data, i, nil)
}

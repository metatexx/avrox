package encoder

import (
	"errors"
	"github.com/hamba/avro/v2"
	"github.com/metatexx/avrox"
	"github.com/nats-io/nats.go"
	"sync"
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
	Compression  avrox.CompressionID
	SchmemaCache sync.Map
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
	cacheIndex := uint64(i.NamespaceID())*uint64(avrox.SchemaMax+1) + uint64(i.SchemaID())
	schema, found := pb.SchmemaCache.Load(cacheIndex)
	if !found {
		var errParse error
		schema, errParse = avro.Parse(i.AVSC())
		if errParse != nil {
			return nil, avrox.ErrSchemaInvalid
		}
		pb.SchmemaCache.Store(cacheIndex, schema)
	}

	data, err := avrox.Marshal(i, pb.Compression, schema.(avro.Schema))
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
	cacheIndex := uint64(i.NamespaceID())*uint64(avrox.SchemaMax+1) + uint64(i.SchemaID())
	schema, found := pb.SchmemaCache.Load(cacheIndex)
	if !found {
		var errParse error
		schema, errParse = avro.Parse(i.AVSC())
		if errParse != nil {
			return avrox.ErrSchemaInvalid
		}
		pb.SchmemaCache.Store(cacheIndex, schema)
	}
	return avrox.Unmarshal(data, i, schema.(avro.Schema))
}

package avrox

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/golang/snappy"
	"github.com/hamba/avro/v2"
	"github.com/metatexx/mxx/wfl"
)

// 0xAC 0bCCCNNNNN 0bSSSSSSSS 0bSSSSSPPP

var Marker byte = 0x93 // one of the best when analysing our data

type Magic [4]byte

type CompressionID int
type NamespaceID int
type SchemaID int

const (
	CompNone   CompressionID = 0
	CompSnappy CompressionID = 1
	CompMax    CompressionID = 7

	// NamespacePrivate means that it is not registered and we use private schemas
	NamespacePrivate NamespaceID = 0
	// NamespaceBasic is reserved for the basic types and structs that are implemented through avrox
	NamespaceBasic NamespaceID = 1
	// NamespaceReserved1 is reserved for later
	NamespaceReserved1 NamespaceID = 2
	// NamespaceReserved2 is reserved for later
	NamespaceReserved2 NamespaceID = 3
	// NamespaceReserved3 is reserved for later
	NamespaceReserved3 NamespaceID = 4
	NamespaceMax       NamespaceID = 31

	// Schema 0 means that it is not defined (but may belong to a namespace)
	SchemaUndefined SchemaID = 0
	SchemaMax       SchemaID = 8191
)

var (
	ErrLengthInvalid           = errors.New("data length should be exactly 4 bytes")
	ErrNamespaceIDOutOfRange   = errors.New("namespace must be between 0 and 31")
	ErrCompressionIDOutOfRange = errors.New("compression must be between 0 and 7")
	ErrCompressionUnsupported  = errors.New("compression type is unsupported")
	ErrSchemaIDOutOfRange      = errors.New("schema must be between 0 and 8191")
	ErrMarkerInvalid           = fmt.Errorf("data should start with magic marker (0x%02x)", Marker)
	ErrParityCheckFailed       = errors.New("parity check failed")
	ErrMarshallingFailed       = errors.New("marshalling failed")
	ErrSchemaNil               = errors.New("schema is nil")
	ErrSchemaInvalid           = errors.New("schema is invalid")
	ErrDecompress              = errors.New("can not decompress")
	ErrDataFormatNotDetected   = errors.New("message format was not detected")
	ErrNoBasicNamespace        = errors.New("no basic namespace")
	ErrNoBasicSchema           = errors.New("no basic schema")
	ErrNoBasicString           = errors.New("no basic string")
	ErrNoBasicInt              = errors.New("no basic int")
	ErrNoBasicByteSlice        = errors.New("no basic byte slice")
	ErrNoBasicMapStringAny     = errors.New("no basic map string any")
	ErrWrongNamespace          = errors.New("namespace from schemer does not fit the magic entry")
	ErrWrongSchema             = errors.New("schema from schemer does not fit the magic entry")
	//ErrBasicTypeNotSupported    = errors.New("basic type not supported")
)

type Schemer interface {
	AVSC() string
	NamespaceID() NamespaceID
	SchemaID() SchemaID
}

func calculateParity(data []byte) byte {
	if len(data) != 4 {
		panic(errors.New("wrong len for calculate parity"))
	}
	// this was the best checksum when testing different implementations
	return (data[1] ^ data[2] + (data[3] & 0xF8)) & 0x07
}

func MustEncodePrivateMagic(compression CompressionID) Magic {
	m, err := EncodeMagic(0, 0, compression)
	if err != nil {
		panic(err)
	}
	return m
}

func MustEncodeBasicMagic(schemaID SchemaID, compression CompressionID) Magic {
	m, err := EncodeMagic(NamespaceBasic, schemaID, compression)
	if err != nil {
		panic(err)
	}
	return m
}

func EncodePrivateMagic(compression CompressionID) (Magic, error) {
	return EncodeMagic(0, 0, compression)
}

func EncodeMagic(namespace NamespaceID, schema SchemaID, compression CompressionID) (Magic, error) {
	if namespace < 0 || namespace > NamespaceMax {
		return Magic{}, ErrNamespaceIDOutOfRange
	}
	if compression < 0 || compression > CompMax {
		return Magic{}, ErrCompressionIDOutOfRange
	}
	if schema < 0 || schema > SchemaMax {
		return Magic{}, ErrSchemaIDOutOfRange
	}

	data := Magic{}
	data[0] = Marker
	data[1] = byte(int(compression<<5) | int(namespace))
	data[2] = byte(schema >> 5)
	data[3] = byte(schema << 3)

	calculatedParity := calculateParity(data[:])

	data[3] |= calculatedParity

	return data, nil
}

func DecodeMagic(data []byte) (NamespaceID, SchemaID, CompressionID, error) {
	if len(data) != 4 {
		return 0, 0, 0, ErrLengthInvalid
	}

	if data[0] != Marker {
		return 0, 0, 0, ErrMarkerInvalid
	}

	compression := CompressionID(int(data[1] >> 5))
	namespace := NamespaceID(int(data[1] & 0x1F))
	schema := SchemaID(int(data[2])<<5 | int(data[3]>>3))

	parityBits := data[3] & 0x07
	calculatedParity := calculateParity(data)

	if parityBits != calculatedParity {
		return 0, 0, 0, ErrParityCheckFailed
	}

	return namespace, schema, compression, nil
}

func IsMagic(data []byte) bool {
	if len(data) != 4 {
		return false
	}

	if data[0] != Marker {
		return false
	}

	parityBits := data[3] & 0x07
	calculatedParity := calculateParity(data)

	return parityBits == calculatedParity
}

func Marshal(src Schemer, cID CompressionID, schema avro.Schema) ([]byte, error) {
	if schema == nil {
		var parseErr error
		schema, parseErr = avro.Parse(src.AVSC())
		if parseErr != nil {
			return nil, ErrSchemaInvalid
		}
	}
	return MarshalAny(src, schema, src.NamespaceID(), src.SchemaID(), cID)
}

func MarshalAny(src any, schema avro.Schema, nID NamespaceID, sID SchemaID, cID CompressionID) ([]byte, error) {
	if schema == nil {
		return nil, wfl.ErrorWithSkip(ErrSchemaNil, 2)
	}
	someValue := reflect.ValueOf(src)
	if someValue.Kind() == reflect.Interface {
		someValue = someValue.Elem()
	}
	if someValue.Kind() == reflect.Ptr {
		someValue = someValue.Elem()
	}
	magicField := someValue.FieldByName("Magic")
	magic, errMagic := EncodeMagic(nID, sID, cID)
	if errMagic != nil {
		return nil, wfl.ErrorWithSkip(errMagic, 2)
	}
	magicField.Set(reflect.ValueOf(magic))
	data, errMarshal := avro.Marshal(schema, src)
	if errMarshal != nil {
		return nil, errors.Join(ErrMarshallingFailed, wfl.ErrorWithSkip(errMarshal, 2))
	}
	//nolint:exhaustive // can't be exhaustive
	switch cID {
	case CompNone:
	case CompSnappy:
		edata := snappy.Encode(nil, data[4:])
		data = append(data[0:4], edata...)
	default:
		return nil, wfl.ErrorWithSkip(ErrCompressionUnsupported, 2)
	}
	return data, nil
}

func MarshalBasic(src any) ([]byte, error) {
	switch v := src.(type) {
	case string:
		kind := &BasicString{
			Magic: MustEncodeBasicMagic(BasicStringID, CompNone),
			Value: v,
		}
		return avro.Marshal(avro.MustParse(BasicStringAVSC), kind)
	case int:
		kind := &BasicInt{
			Magic: MustEncodeBasicMagic(BasicIntID, CompNone),
			Value: v,
		}
		return avro.Marshal(avro.MustParse(BasicIntAVSC), kind)
	case []byte:
		kind := &BasicByteSlice{
			Magic: MustEncodeBasicMagic(BasicByteSliceID, CompNone),
			Value: v,
		}
		return avro.Marshal(avro.MustParse(BasicByteSliceAVSC), kind)
	case map[string]any:
		kind := &BasicMapStringAny{
			Magic: MustEncodeBasicMagic(BasicMapStringAnyID, CompNone),
			Value: v,
		}
		return avro.Marshal(avro.MustParse(BasicMapStringAnyAVSC), kind)
	default:
		return nil, errors.New("unsupported type")
	}
}

// Unmarshal uses the give schema for unmarshalling and checks if
// it fits to the decode data. This function is faster if the schema is given
// When the schema is not given it will parse the Schemer info.
// If the schema is given, it will check that this matches to the Schemer info
func Unmarshal(data []byte, dst Schemer, schema avro.Schema) error {
	if len(data) > 1 && data[0] == '{' {
		// TODO: JSON We need to get the namespace and schema from the json data
		return json.Unmarshal(data, dst)
	} else if len(data) > 4 && IsMagic(data[0:4]) {
		var cID CompressionID
		if schema == nil {
			var errSchema error
			schema, errSchema = avro.Parse(dst.AVSC())
			if errSchema != nil {
				return errSchema
			}
			var errDecodeMagic error
			_, _, cID, errDecodeMagic = DecodeMagic(data[0:4])
			if errDecodeMagic != nil {
				return errDecodeMagic
			}
		} else {
			var nID NamespaceID
			var sID SchemaID
			var errDecodeMagic error
			nID, sID, cID, errDecodeMagic = DecodeMagic(data[0:4])
			if errDecodeMagic != nil {
				return errDecodeMagic
			}
			if nID != dst.NamespaceID() {
				return ErrWrongNamespace
			}
			if sID != dst.SchemaID() {
				return ErrWrongSchema
			}
		}
		//nolint:exhaustive // can't be exhaustive
		switch cID {
		case CompNone:
			break
		case CompSnappy:
			// we encode the part after the magic (this can most likely be optimized)
			edata, errDecode := snappy.Decode(nil, data[4:])
			if errDecode != nil {
				return wfl.Error(ErrDecompress)
			}
			data = append(data[0:4], edata...)
		default:
			return ErrCompressionUnsupported
		}
		return avro.Unmarshal(schema, data, dst)
	}
	return ErrDataFormatNotDetected
}

func UnmarshalAny[T any](data []byte, schema avro.Schema, dst *T) (NamespaceID, SchemaID, error) {
	if len(data) > 1 && data[0] == '{' {
		// TODO: JSON We need to get the namespace and schema from the json data
		return 0, 0, json.Unmarshal(data, dst)
	}

	if len(data) <= 4 || !IsMagic(data[0:4]) {
		return 0, 0, ErrDataFormatNotDetected
	}

	// AVRO (with MagicV1)
	nID, sID, cID, errMagic := DecodeMagic(data[0:4])
	if errMagic != nil {
		return 0, 0, errMagic
	}
	//nolint:exhaustive // can't be exhaustive
	switch cID {
	case CompNone:
		break
	case CompSnappy:
		var errDecode error
		data, errDecode = snappy.Decode(nil, data[:])
		if errDecode != nil {
			return 0, 0, ErrDecompress
		}
	default:
		return 0, 0, ErrCompressionUnsupported
	}
	return nID, sID, avro.Unmarshal(schema, data, dst)
}

func UnmarshalBasic(src []byte) (any, error) {
	nID, sID, cID, err := DecodeMagic(src[:4])
	if err != nil {
		return nil, err
	}
	if nID != NamespaceBasic {
		return nil, ErrNoBasicNamespace
	}
	//nolint:exhaustive // can't be exhaustive
	switch cID {
	case CompNone:
	case CompSnappy:
		src, err = snappy.Decode(nil, src[:])
		if err != nil {
			return nil, ErrDecompress
		}
	}
	var errUnmarshalAny error

	//nolint:exhaustive // only runs for NamespaceBasic
	switch sID {
	case BasicStringID:
		kind := &BasicString{}
		nID, sID, errUnmarshalAny = UnmarshalAny(src, avro.MustParse(BasicStringAVSC), kind)
		if errUnmarshalAny != nil {
			return nil, errUnmarshalAny
		}
		if nID != NamespaceBasic && sID != BasicStringID {
			return nil, ErrNoBasicString
		}
		return kind.Value, nil
	case BasicIntID:
		kind := &BasicInt{}
		nID, sID, errUnmarshalAny = UnmarshalAny(src, avro.MustParse(BasicIntAVSC), kind)
		if errUnmarshalAny != nil {
			return nil, errUnmarshalAny
		}
		if nID != NamespaceBasic && sID != BasicIntID {
			return nil, ErrNoBasicInt
		}
		return kind.Value, nil
	case BasicByteSliceID:
		kind := &BasicByteSlice{}
		nID, sID, errUnmarshalAny = UnmarshalAny(src, avro.MustParse(BasicByteSliceAVSC), kind)
		if errUnmarshalAny != nil {
			return nil, errUnmarshalAny
		}
		if nID != NamespaceBasic && sID != BasicByteSliceID {
			return nil, ErrNoBasicByteSlice
		}
		return kind.Value, nil
	case BasicMapStringAnyID:
		kind := &BasicMapStringAny{}
		nID, sID, errUnmarshalAny = UnmarshalAny(src, avro.MustParse(BasicMapStringAnyAVSC), kind)
		if errUnmarshalAny != nil {
			return nil, errUnmarshalAny
		}
		if nID != NamespaceBasic && sID != BasicByteSliceID {
			return nil, ErrNoBasicMapStringAny
		}
		return kind.Value, nil
	default:
		return nil, ErrNoBasicSchema
	}
}

func UnmarshalString(data []byte) (string, error) {
	x, err := UnmarshalBasic(data)
	if err != nil {
		return "", err
	}
	if str, ok := x.(string); ok {
		return str, nil
	}
	return "", ErrNoBasicString
}

func UnmarshalInt(data []byte) (int, error) {
	x, err := UnmarshalBasic(data)
	if err != nil {
		return 0, err
	}
	if n, ok := x.(int); ok {
		return n, nil
	}
	return 0, ErrNoBasicString
}

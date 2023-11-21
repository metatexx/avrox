package avrox

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hamba/avro/v2"
	"github.com/metatexx/mxx/wfl"
	"reflect"
)

type CompressionID int
type NamespaceID int
type SchemaID int // Schema<<8 | 8 bit version

func PackSchemVer(schemaID SchemaID, version int) SchemaID {
	return SchemaID(int(schemaID)<<8 | (version & 0xff))
}

func UnpackSchemVer(schemaVerID SchemaID) (int, int) {
	return int(schemaVerID >> 8), int(schemaVerID & 0xff)
}

const (
	CompNone   CompressionID = 0
	CompSnappy CompressionID = 1
	CompFlate  CompressionID = 2 // Uses -1 as compression parameter
	CompGZip   CompressionID = 3 // Uses -1 as compression parameter
	CompMax    CompressionID = 255

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
	NamespaceMax       NamespaceID = 65535

	// Schema 0 means that it is not defined (but may belong to a namespace)
	SchemaUndefined SchemaID = 0
	SchemaMax       SchemaID = 16777215 // Schema<<8 | Version

	MagicFieldName = "Magic"
)

var (
	ErrLengthInvalid           = errors.New("data length should be exactly 8 bytes")
	ErrNamespaceIDOutOfRange   = errors.New("namespace must be between 0 and 31")
	ErrCompressionIDOutOfRange = errors.New("compression must be between 0 and 7")
	ErrCompressionUnsupported  = errors.New("compression type is unsupported")
	ErrSchemaIDOutOfRange      = errors.New("schema must be between 0 and 8191")
	ErrMarkerInvalid           = fmt.Errorf("data should start with magic marker (0x%02x)", Marker)
	ErrParityCheckFailed       = errors.New("parity check failed")
	ErrMarshallingFailed       = errors.New("marshalling failed")
	ErrMissingMagicField       = errors.New("missing magic field in struct")
	ErrMarshallAnyWithoutPtr   = errors.New("no ptr src for MarshalAny")
	ErrSchemaNil               = errors.New("schema is nil")
	ErrSchemaInvalid           = errors.New("schema is invalid")
	ErrDecompress              = errors.New("can not decompress")
	ErrDataFormatNotDetected   = errors.New("message format was not detected")
	ErrNoData                  = errors.New("no data")
	ErrNoBasicNamespace        = errors.New("no basic namespace")
	ErrNoBasicSchema           = errors.New("no basic schema")
	ErrNoBasicString           = errors.New("no basic string")
	ErrNoBasicInt              = errors.New("no basic int")
	ErrNoBasicByteSlice        = errors.New("no basic byte slice")
	ErrNoBasicMapStringAny     = errors.New("no basic map string any")
	ErrNoBasicTime             = errors.New("no basic time")
	ErrWrongNamespace          = errors.New("namespace from schemer does not fit the magic entry")
	ErrWrongSchema             = errors.New("schema from schemer does not fit the magic entry")
	ErrNotAvroX                = errors.New("data is not avrox")
	ErrNoPointerDestination    = errors.New("not a pointer destination")
	ErrSchemerNotFound         = errors.New("schema from schemer is not in the given slice")
	//ErrBasicTypeNotSupported    = errors.New("basic type not supported")
)

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
	} else {
		return nil, ErrMarshallAnyWithoutPtr
	}
	magicField := someValue.FieldByName(MagicFieldName)
	if !magicField.IsValid() {
		return nil, ErrMissingMagicField
	}
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
	return CompressData(data, cID)

}

func unmarshalHelper(data []byte) ([]byte, NamespaceID, SchemaID, error) {
	if len(data) <= MagicLen || !IsMagic(data[0:MagicLen]) {
		return nil, 0, 0, ErrDataFormatNotDetected
	}

	nID, sID, cID, errMagic := DecodeMagic(data[0:MagicLen])
	if errMagic != nil {
		return nil, 0, 0, errMagic
	}

	uncompressed, err := DecompressData(data, cID)
	return uncompressed, nID, sID, err
}

func UnmarshalAny[T any](data []byte, schema avro.Schema, dst *T) (NamespaceID, SchemaID, error) {
	if len(data) == 0 {
		return 0, 0, nil
	}

	if len(data) > 1 && data[0] == '{' {
		// TODO: JSON We need to get the namespace and schema from the json data
		return 0, 0, json.Unmarshal(data, dst)
	}

	data, nID, sID, errHelper := unmarshalHelper(data)
	if errHelper != nil {
		return 0, 0, errHelper
	}

	return nID, sID, avro.Unmarshal(schema, data, dst)
}

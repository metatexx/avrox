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
type SchemaID int

const (
	CompNone   CompressionID = 0
	CompSnappy CompressionID = 1
	CompFlate  CompressionID = 2 // Uses -1 as compression parameter
	CompGZip   CompressionID = 3 // Uses -1 as compression parameter
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

	MagicFieldName = "Magic"
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
	ErrMissingMagicField       = errors.New("missing magic field in struct")
	ErrMarshallAnyWithoutPtr   = errors.New("no ptr src for MarshalAny")
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
	ErrNoBasicTime             = errors.New("no basic time")
	ErrWrongNamespace          = errors.New("namespace from schemer does not fit the magic entry")
	ErrWrongSchema             = errors.New("schema from schemer does not fit the magic entry")
	ErrNotAvroX                = errors.New("data is not avrox")
	ErrNoPointerDestination    = errors.New("not a pointer destination")
	ErrSchemerNotFound         = errors.New("schema from schemer is not in the given slice")
	//ErrBasicTypeNotSupported    = errors.New("basic type not supported")
)

type Schemer interface {
	AVSC() string
	NamespaceID() NamespaceID
	SchemaID() SchemaID
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
	return compressData(data, cID)

}

func unmarshalHelper(data []byte) ([]byte, NamespaceID, SchemaID, error) {
	if len(data) <= 4 || !IsMagic(data[0:4]) {
		return nil, 0, 0, ErrDataFormatNotDetected
	}

	nID, sID, cID, errMagic := DecodeMagic(data[0:4])
	if errMagic != nil {
		return nil, 0, 0, errMagic
	}

	var err error
	data, err = decompressData(data, cID)
	return data, nID, sID, err
}

// Unmarshal uses the give schema for unmarshalling and checks if
// it fits to the decode data. This function is faster if the schema is given
// When the schema is not given it will parse the Schemer info.
// If the schema is given, it will check that this matches to the Schemer info
func Unmarshal(data []byte, dst Schemer, schema avro.Schema) error {
	if len(data) == 0 {
		dst = nil
		return nil
	}

	if len(data) > 1 && data[0] == '{' {
		// TODO: JSON We need to get the namespace and schema from the json data
		return json.Unmarshal(data, dst)
	}

	data, nID, sID, errHelper := unmarshalHelper(data)
	if errHelper != nil {
		return errHelper
	}

	if schema == nil {
		var errSchema error
		schema, errSchema = avro.Parse(dst.AVSC())
		if errSchema != nil {
			return errSchema
		}
	}
	if nID != dst.NamespaceID() {
		return ErrWrongNamespace
	}
	if sID != dst.SchemaID() {
		return ErrWrongSchema
	}

	return avro.Unmarshal(schema, data, dst)
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

// UnmarshalUnion expects a slice with preallocated schemers and uses the magic
// in the data to unmarshal the correct one. It will return the used schema as any
// If no schema fits it will return an error
func UnmarshalSchemer(src []byte, schemers ...Schemer) (any, error) {
	if len(src) == 0 {
		return nil, nil
	}
	if len(src) < 4 {
		return nil, ErrNotAvroX
	}
	nID, sID, _, err := DecodeMagic(src[:4])
	if err != nil {
		return nil, err
	}

	for _, schemer := range schemers {
		if schemer.NamespaceID() == nID && schemer.SchemaID() == sID {
			// When we only get a nil pointer of the Schemer type, we allocate one
			v := reflect.ValueOf(schemer)
			if v.Kind() == reflect.Ptr {
				if v.IsNil() {
					schemer = reflect.New(v.Type().Elem()).Interface().(Schemer)
				}
			} else {
				return nil, ErrNoPointerDestination
			}
			return schemer, Unmarshal(src, schemer, nil)
		}
	}
	return nil, ErrSchemerNotFound
}

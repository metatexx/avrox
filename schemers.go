package avrox

import (
	"encoding/json"
	"github.com/hamba/avro/v2"
	"reflect"
	"strings"
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

// UnmarshalUnion expects a slice with preallocated schemers and uses the magic
// in the data to unmarshal the correct one. It will return the used schema as any
// If no schema fits it will return an error
func UnmarshalSchemer(src []byte, schemers ...Schemer) (any, error) {
	if len(src) == 0 {
		return nil, nil
	}
	if len(src) < MagicLen {
		return nil, ErrNotAvroX
	}
	nID, sID, _, err := DecodeMagic(src[:MagicLen])
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

// JoinedSchemas returns a json array of all schemers in the arguments
func JoinedSchemas(schemers ...Schemer) string {
	var sb strings.Builder
	if len(schemers) > 1 {
		sb.WriteString("[")
		for idx, schema := range schemers {
			if idx > 0 {
				sb.WriteRune(',')
			}
			sb.WriteString(schema.AVSC())
		}
		sb.WriteString("]")
	} else if len(schemers) == 1 {
		sb.WriteString(schemers[0].AVSC())
	}
	return sb.String()
}

package avrox

import _ "embed"

//go:generate avscgen -n "basics" -o avsc/ . BasicByteSlice
//go:embed avsc/basic_byte_slice.avsc
var BasicByteSliceAVSC string

type BasicByteSlice struct {
	Magic [MagicLen]byte // 1.3.1
	Value []byte
}

// AVSC returns the AVRO schema for the BasicString struct type
func (_ *BasicByteSlice) Schema() string {
	return BasicByteSliceAVSC
}

// NamespaceID returns the namespace id for the BasicInt struct type
func (_ *BasicByteSlice) NamespaceID() NamespaceID {
	return NamespaceBasic
}

// SchemaID returns the schema id for the BasicInt struct type
func (_ *BasicByteSlice) SchemaID() SchemaID {
	return BasicByteSliceSchemaID
}

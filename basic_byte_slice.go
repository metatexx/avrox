package avrox

import _ "embed"

//go:generate avscgen -ns "testing" . BasicByteSlice
//go:embed basic_byte_slice.avsc
var BasicByteSliceAVSC string

type BasicByteSlice struct {
	Magic [4]byte
	Value []byte
}

// AVSC returns the AVRO schema for the BasicString struct type
func (s *BasicByteSlice) AVSC() string {
	return BasicByteSliceAVSC
}

// NamespaceID returns the namespace id for the BasicInt struct type
func (s *BasicByteSlice) NamespaceID() NamespaceID {
	return NamespaceBasic
}

// SchemaID returns the schema id for the BasicInt struct type
func (s *BasicByteSlice) SchemaID() SchemaID {
	return BasicByteSliceID
}

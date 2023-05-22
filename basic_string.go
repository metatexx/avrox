package avrox

import _ "embed"

// Implementation of BasicString

// BasicString is the container type to store a string in a single avro schema
type BasicString struct {
	Magic [MagicLen]byte // 1.1.1
	Value string
}

//go:generate avscgen -n "basics" -o avsc/ . BasicString
//go:embed avsc/basic_string.avsc
var BasicStringAVSC string

// AVSC returns the AVRO schema for the BasicString struct type
func (_ *BasicString) Schema() string {
	return BasicStringAVSC
}

// NamespaceID returns the namespace id for the BasicString struct type
func (_ *BasicString) NamespaceID() NamespaceID {
	return NamespaceBasic
}

// SchemaID returns the schema id for the BasicString struct type
func (_ *BasicString) SchemaID() SchemaID {
	return BasicStringSchemaID
}

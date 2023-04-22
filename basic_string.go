package avrox

import _ "embed"

// Implementation of BasicString

// BasicString is the container type to store a string in a single avro schema
type BasicString struct {
	Magic [4]byte
	Value string
}

//go:generate avscgen -ns "testing" . BasicString
//go:embed basic_string.avsc
var BasicStringAVSC string

// AVSC returns the AVRO schema for the BasicString struct type
func (s *BasicString) AVSC() string {
	return BasicStringAVSC
}

// NamespaceID returns the namespace id for the BasicString struct type
func (s *BasicString) NamespaceID() NamespaceID {
	return NamespaceBasic
}

// SchemaID returns the schema id for the BasicString struct type
func (s *BasicString) SchemaID() SchemaID {
	return BasicStringID
}

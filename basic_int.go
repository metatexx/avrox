package avrox

import _ "embed"

//go:generate avscgen -ns "testing" . BasicInt
//go:embed basic_int.avsc
var BasicIntAVSC string

type BasicInt struct {
	Magic [4]byte
	Value int
}

// AVSC returns the AVRO schema for the BasicString struct type
func (s *BasicInt) AVSC() string {
	return BasicIntAVSC
}

// NamespaceID returns the namespace id for the BasicInt struct type
func (s *BasicInt) NamespaceID() NamespaceID {
	return NamespaceBasic
}

// SchemaID returns the schema id for the BasicInt struct type
func (s *BasicInt) SchemaID() SchemaID {
	return BasicIntID
}

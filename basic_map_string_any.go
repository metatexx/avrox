package avrox

import _ "embed"

//go:generate avscgen -ns "basics" -o avsc/ . BasicMapStringAny
//go:embed avsc/basic_map_string_any.avsc
var BasicMapStringAnyAVSC string

type BasicMapStringAny struct {
	Magic [4]byte
	Value map[string]any
}

// AVSC returns the AVRO schema for the BasicString struct type
func (s *BasicMapStringAny) AVSC() string {
	return BasicMapStringAnyAVSC
}

// NamespaceID returns the namespace id for the BasicInt struct type
func (s *BasicMapStringAny) NamespaceID() NamespaceID {
	return NamespaceBasic
}

// SchemaID returns the schema id for the BasicInt struct type
func (s *BasicMapStringAny) SchemaID() SchemaID {
	return BasicMapStringAnyID
}

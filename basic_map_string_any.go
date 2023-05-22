package avrox

import _ "embed"

//go:generate avscgen -n "basics" -o avsc/ . BasicMapStringAny
//go:embed avsc/basic_map_string_any.avsc
var BasicMapStringAnyAVSC string

type BasicMapStringAny struct {
	Magic [MagicLen]byte // 1.4.1
	Value map[string]any
}

// AVSC returns the AVRO schema for the BasicString struct type
func (_ *BasicMapStringAny) Schema() string {
	return BasicMapStringAnyAVSC
}

// NamespaceID returns the namespace id for the BasicInt struct type
func (_ *BasicMapStringAny) NamespaceID() NamespaceID {
	return NamespaceBasic
}

// SchemaID returns the schema id for the BasicInt struct type
func (_ *BasicMapStringAny) SchemaID() SchemaID {
	return BasicMapStringAnySchemaID
}

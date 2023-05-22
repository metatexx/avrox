package avrox

import _ "embed"

//go:generate avscgen -n "basics" -o avsc/ . BasicInt
//go:embed avsc/basic_int.avsc
var BasicIntAVSC string

type BasicInt struct {
	Magic [MagicLen]byte // 1.2.1
	Value int
}

// AVSC returns the AVRO schema for the BasicString struct type
func (_ *BasicInt) Schema() string {
	return BasicIntAVSC
}

// NamespaceID returns the namespace id for the BasicInt struct type
func (_ *BasicInt) NamespaceID() NamespaceID {
	return NamespaceBasic
}

// SchemaID returns the schema id for the BasicInt struct type
func (_ *BasicInt) SchemaID() SchemaID {
	return BasicIntSchemaID
}

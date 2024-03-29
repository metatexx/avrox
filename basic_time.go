package avrox

import (
	_ "embed"
	"time"
)

// Implementation of BasicTime
var _ Schemer = (*BasicTime)(nil)

// BasicTime is the container type to store a timestamp in a single avro schema
type BasicTime struct {
	Magic [MagicLen]byte // 1.5.1
	Value time.Time
}

//go:generate avscgen -n "basics" -o avsc/ . BasicTime
//go:embed avsc/basic_time.avsc
var BasicTimeAVSC string

// Schema returns the AVRO schema for the BasicTime struct type
func (BasicTime) Schema() string {
	return BasicTimeAVSC
}

// NamespaceID returns the namespace id for the BasicTime struct type
func (BasicTime) NamespaceID() NamespaceID {
	return NamespaceBasic
}

// SchemaID returns the schema id for the BasicTime struct type
func (BasicTime) SchemaID() SchemaID {
	return BasicTimeSchemaID
}

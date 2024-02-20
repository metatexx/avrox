package avrox

import (
	_ "embed"
	"github.com/metatexx/avrox/rawdate"
)

// Implementation of BasicTime
var _ Schemer = (*BasicRawDate)(nil)

/*
type RawDate struct {
	Year  int
	Month int8
	Day   int8
}
*/

// BasicRawDate is the container type to store a timestamp in a single avro schema
type BasicRawDate struct {
	Magic [MagicLen]byte // 1.5.1
	Value rawdate.RawDate
}

//go:generate avscgen -n "basics" -o avsc/ . BasicRawDate
//go:embed avsc/basic_raw_date.avsc
var BasicRawDateAVSC string

// Schema returns the AVRO schema for the BasicRawDate struct type
func (BasicRawDate) Schema() string {
	return BasicRawDateAVSC
}

// NamespaceID returns the namespace id for the BasicRawDate struct type
func (BasicRawDate) NamespaceID() NamespaceID {
	return NamespaceBasic
}

// SchemaID returns the schema id for the BasicRawDate struct type
func (BasicRawDate) SchemaID() SchemaID {
	return BasicRawDateSchemaID
}

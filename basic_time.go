package avrox

import (
	_ "embed"
	"time"
)

// Implementation of BasicString

// BasicString is the container type to store a string in a single avro schema
type BasicTime struct {
	Magic [4]byte
	Value time.Time
}

//go:generate avscgen -ns "testing" . BasicTime
//go:embed basic_time.avsc
var BasicTimeAVSC string

// AVSC returns the AVRO schema for the BasicTime struct type
func (s *BasicTime) AVSC() string {
	return BasicTimeAVSC
}

// NamespaceID returns the namespace id for the BasicTime struct type
func (s *BasicTime) NamespaceID() NamespaceID {
	return NamespaceBasic
}

// SchemaID returns the schema id for the BasicTime struct type
func (s *BasicTime) SchemaID() SchemaID {
	return BasicTimeID
}

package avrox

import (
	_ "embed"
	"math/big"
)

// Implementation of BasicDecimal

// BasicDecima is the container type to store a *bigRat value into a single avro schema
type BasicDecimal struct {
	Magic [MagicLen]byte
	Value *big.Rat
}

//go:generate avscgen -ns "basics" -o avsc/ . BasicDecimal
//go:embed avsc/basic_decimal.avsc
var BasicDecimalAVSC string

// AVSC returns the AVRO schema for the BasicDecimal struct type
func (s *BasicDecimal) AVSC() string {
	return BasicDecimalAVSC
}

// NamespaceID returns the namespace id for the BasicDecimal struct type
func (s *BasicDecimal) NamespaceID() NamespaceID {
	return NamespaceBasic
}

// SchemVerID returns the schema id for the BasicDecimal struct type
func (s *BasicDecimal) SchemaID() SchemVerID {
	return BasicTimeID
}
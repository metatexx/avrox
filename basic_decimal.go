package avrox

import (
	_ "embed"
	"math/big"
)

// Implementation of BasicDecimal

// BasicDecima is the container type to store a *bigRat value into a single avro schema
type BasicDecimal struct {
	Magic [4]byte
	Value *big.Rat `avsc:"precision:12,scale:2"`
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

// SchemaID returns the schema id for the BasicDecimal struct type
func (s *BasicDecimal) SchemaID() SchemaID {
	return BasicTimeID
}

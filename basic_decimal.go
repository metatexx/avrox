package avrox

import (
	_ "embed"
	"math/big"
)

// Implementation of BasicDecimal
var _ Schemer = (*BasicDecimal)(nil)

// BasicDecimal is the container type to store a *bigRat value into a single avro schema
type BasicDecimal struct {
	Magic [MagicLen]byte // 1.6.1
	Value *big.Rat
}

//go:generate avscgen -n "basics" -o avsc/ . BasicDecimal
//go:embed avsc/basic_decimal.avsc
var BasicDecimalAVSC string

// Schema returns the AVRO schema for the BasicDecimal struct type
func (BasicDecimal) Schema() string {
	return BasicDecimalAVSC
}

// NamespaceID returns the namespace id for the BasicDecimal struct type
func (BasicDecimal) NamespaceID() NamespaceID {
	return NamespaceBasic
}

// SchemaID returns the schema id for the BasicDecimal struct type
func (BasicDecimal) SchemaID() SchemaID {
	return BasicTimeSchemaID
}

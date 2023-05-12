package avrox

// This contains some of the basic types in a predefined namespace

// SchemaIDs for the supported basic types

const (
	// BasicStringSchemaID is the id for the avro schema of struct BasicString
	BasicStringSchemaID SchemaID = 1

	// BasicIntSchemaID is the id for the avro schema of struct BasicInt
	BasicIntSchemaID SchemaID = 2

	// BasicByteSliceSchemaID is the id for the avro schema of struct BasicInt
	BasicByteSliceSchemaID SchemaID = 3

	// BasicMapStringAnySchemaID is the id for the avro schema of struct BasicMapStringAny
	BasicMapStringAnySchemaID SchemaID = 4

	// BasicTimeSchemaID is the id for the avro schema of struct BasicTime
	BasicTimeSchemaID SchemaID = 5

	// BasicDecimalSchemaID is the id for the avro schema of struct BasicDecimal (*big.Rat / decimal.ficed)
	BasicDecimalSchemaID SchemaID = 6
)

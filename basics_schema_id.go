package avrox

// This contains some of the basic types in a predefined namespace

// SchemaIDs for the supported basic types

const (
	// BasicStringID is the id for the avro schema of struct BasicString
	BasicStringID SchemaID = 1

	// BasicIntID is the id for the avro schema of struct BasicInt
	BasicIntID SchemaID = 2

	// BasicByteSliceID is the id for the avro schema of struct BasicInt
	BasicByteSliceID SchemaID = 3

	// BasicMapStringAnyID is the id for the avro schema of struct BasicMapStringAny
	BasicMapStringAnyID SchemaID = 4

	// BasicTimeID is the id for the avro schema of struct BasicTime
	BasicTimeID SchemaID = 5

	// BasicDecimalID is the id for the avro schema of struct BasicDecimal (*big.Rat / decimal.ficed)
	BasicDecimalID SchemaID = 6
)

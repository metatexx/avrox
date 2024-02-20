package avrox

// This contains some of the basic types in a predefined namespace

// SchemaIDs for the supported basic types

const (
	// BasicStringSchemaID is the id for the avro schema of struct BasicString
	BasicStringSchemaID SchemaID = 1<<8 + 1

	// BasicIntSchemaID is the id for the avro schema of struct BasicInt
	BasicIntSchemaID SchemaID = 2<<8 + 1

	// BasicByteSliceSchemaID is the id for the avro schema of struct BasicInt
	BasicByteSliceSchemaID SchemaID = 3<<8 + 1

	// BasicMapStringAnySchemaID is the id for the avro schema of struct BasicMapStringAny
	BasicMapStringAnySchemaID SchemaID = 4<<8 + 1

	// BasicTimeSchemaID is the id for the avro schema of struct BasicTime
	BasicTimeSchemaID SchemaID = 5<<8 + 1

	// BasicDecimalSchemaID is the id for the avro schema of struct BasicDecimal (*big.Rat / decimal.fixed)
	BasicDecimalSchemaID SchemaID = 6<<8 + 1

	// BasicRawDateSchemaID is the id for the avro schema of struct BasicRawDate (rawdate.Rawdate)
	BasicRawDateSchemaID SchemaID = 7<<8 + 1
)

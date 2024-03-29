package avrox_test

import (
	"errors"
	"github.com/hamba/avro/v2"
	"github.com/metatexx/avrox/rawdate"
	"github.com/metatexx/avrox/testdata"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/metatexx/avrox"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	testCases := []struct {
		namespace   avrox.NamespaceID
		compression avrox.CompressionID
		schema      avrox.SchemaID
		shouldError bool
	}{
		{0, 0, 0, false},
		{1, 3, 1023, false},
		{31, 7, 8191, false},
		{17, 2, 4567, false},
		{-1, 0, 0, true},
		{avrox.NamespaceMax + 1, 0, 0, true},
		{avrox.NamespaceMax, 0, 0, false},
		{0, -1, 0, true},
		{0, avrox.CompMax, 0, false},
		{0, avrox.CompMax + 1, 0, true},
		{0, 0, -1, true},
		{0, 0, avrox.SchemaMax, false},
		{0, 0, avrox.SchemaMax + 1, true},
	}

	// Test encoding with valid and invalid input values
	for _, testCase := range testCases {
		encodedData, err := avrox.EncodeMagic(testCase.namespace, testCase.schema, testCase.compression)
		if testCase.shouldError && err == nil {
			t.Errorf("Expected an error, but got no error (Namespace: %d, Compression: %d, Schema: %d)",
				testCase.namespace, testCase.compression, testCase.schema)
			continue
		}
		if !testCase.shouldError && err != nil {
			t.Errorf("Unexpected error during encoding: %v", err)
			continue
		}
		if testCase.shouldError {
			continue
		}

		// Test decoding with correct data
		dNS, dSchema, dCompression, err := avrox.DecodeMagic(encodedData[:])
		if err != nil {
			t.Errorf("Unexpected error during decoding: %v", err)
			continue
		}
		if dNS != testCase.namespace || dCompression != testCase.compression || dSchema !=
			testCase.schema {
			t.Errorf("Mismatch in decoded values. Expected: (Namespace: %d, Compression: %d, Schema: %d) "+
				"Got: (Namespace: %d, Compression: %d, Schema: %d)",
				testCase.namespace, testCase.compression, testCase.schema, dNS, dCompression, dSchema)
		}
	}

	// Test decoding with incorrect data length
	_, _, _, err := avrox.DecodeMagic([]byte{avrox.Marker})
	if err == nil {
		t.Error("Expected an error for incorrect data length, but got no error")
	}

	// Test decoding with incorrect first byte
	_, _, _, err = avrox.DecodeMagic([]byte{0xAD, 0x53, 0x24, 0x1A})
	if err == nil {
		t.Error("Expected an error for incorrect first byte, but got no error")
	}

	// Test decoding with incorrect parity
	_, _, _, err = avrox.DecodeMagic([]byte{avrox.Marker, 0x53, 0x24, 0x1B})
	if err == nil {
		t.Error("Expected an error for incorrect parity, but got no error")
	}
}

func TestMarshal(t *testing.T) {
	// we can use BasicString struct here (just use another magic)
	mt := avrox.BasicString{
		Value: "foo",
	}
	_, err := avrox.Marshal(&mt, avrox.CompMax+1, nil)
	//t.Logf("%s", err)
	assert.True(t, errors.Is(err, avrox.ErrCompressionIDOutOfRange))

	data, err2 := avrox.Marshal(&mt, avrox.CompNone, nil)
	assert.NoError(t, err2)

	n, s, c, err3 := avrox.DecodeMagic(data[:avrox.MagicLen])
	assert.NoError(t, err3)
	assert.Equal(t, avrox.NamespaceBasic, n)
	assert.Equal(t, avrox.BasicStringSchemaID, s)
	assert.Equal(t, avrox.CompNone, c)

	basicString := &avrox.BasicString{}
	err4 := avrox.Unmarshal(data, basicString, nil)
	assert.NoError(t, err4)
	assert.Equal(t, "foo", basicString.Value)
}

func TestMarshalSnappy(t *testing.T) {
	// we can use BasicString struct here (just use another magic)
	mt := avrox.BasicString{
		Value: "foobarfoo",
	}
	data, errMarshal := avrox.Marshal(&mt, avrox.CompSnappy, nil)
	assert.NoError(t, errMarshal)

	n, s, c, errDecodeMagic := avrox.DecodeMagic(data[:avrox.MagicLen])
	assert.NoError(t, errDecodeMagic)
	assert.Equal(t, avrox.NamespaceBasic, n)
	assert.Equal(t, avrox.BasicStringSchemaID, s)
	assert.Equal(t, avrox.CompSnappy, c)

	basicString := &avrox.BasicString{}
	errUnmarshal := avrox.Unmarshal(data, basicString, nil)
	assert.NoError(t, errUnmarshal)
	assert.Equal(t, "foobarfoo", basicString.Value)
}

func TestMarshalBasicString(t *testing.T) {
	data, err := avrox.MarshalBasic("bar", avrox.CompNone)
	assert.NoError(t, err)

	n, s, c, err3 := avrox.DecodeMagic(data[:avrox.MagicLen])
	assert.NoError(t, err3)
	assert.Equal(t, avrox.NamespaceBasic, n)
	assert.Equal(t, avrox.BasicStringSchemaID, s)
	assert.Equal(t, avrox.CompNone, c)

	text, err4 := avrox.UnmarshalBasic(data)
	assert.NoError(t, err4)
	assert.Equal(t, "bar", text)
}

func TestMarshalBasicInt(t *testing.T) {
	data, err := avrox.MarshalBasic(42, avrox.CompNone)
	assert.NoError(t, err)

	n, s, c, err3 := avrox.DecodeMagic(data[:avrox.MagicLen])
	assert.NoError(t, err3)
	assert.Equal(t, avrox.NamespaceBasic, n)
	assert.Equal(t, avrox.BasicIntSchemaID, s)
	assert.Equal(t, avrox.CompNone, c)

	text, err4 := avrox.UnmarshalBasic(data)
	assert.NoError(t, err4)
	assert.Equal(t, 42, text)
}

func TestUnmarshalBasicTypes(t *testing.T) {
	inString := "bar"
	data, err := avrox.MarshalBasic(inString, avrox.CompNone)
	assert.NoError(t, err)
	var outString string
	outString, err = avrox.UnmarshalString(data)
	assert.NoError(t, err)
	assert.Equal(t, inString, outString)

	inInt := -42
	data, err = avrox.MarshalBasic(inInt, avrox.CompNone)
	assert.NoError(t, err)
	var outInt int
	outInt, err = avrox.UnmarshalInt(data)
	assert.NoError(t, err)
	assert.Equal(t, inInt, outInt)
}

func TestMarshalBasicByteSlice(t *testing.T) {
	value := []byte{0, 1, 2, 3}
	data, err := avrox.MarshalBasic(value, avrox.CompNone)
	assert.NoError(t, err)

	n, s, c, err3 := avrox.DecodeMagic(data[:avrox.MagicLen])
	assert.NoError(t, err3)
	assert.Equal(t, avrox.NamespaceBasic, n)
	assert.Equal(t, avrox.BasicByteSliceSchemaID, s)
	assert.Equal(t, avrox.CompNone, c)

	result, err4 := avrox.UnmarshalBasic(data)
	assert.NoError(t, err4)
	assert.Equal(t, value, result)
}

func TestMarshalBasicTime(t *testing.T) {
	value := avrox.AvroTime(time.Now())
	data, err := avrox.MarshalBasic(value, avrox.CompNone)
	assert.NoError(t, err)

	n, s, c, err3 := avrox.DecodeMagic(data[:avrox.MagicLen])
	assert.NoError(t, err3)
	assert.Equal(t, avrox.NamespaceBasic, n)
	assert.Equal(t, avrox.BasicTimeSchemaID, s)
	assert.Equal(t, avrox.CompNone, c)

	result, err4 := avrox.UnmarshalBasic(data)
	assert.NoError(t, err4)
	assert.Equal(t, value, result)
}

func TestMarshalBasicDecimal(t *testing.T) {
	value, _ := (&big.Rat{}).SetString("1/10")
	data, err := avrox.MarshalBasic(value, avrox.CompNone)
	assert.NoError(t, err)

	n, s, c, err3 := avrox.DecodeMagic(data[:avrox.MagicLen])
	assert.NoError(t, err3)
	assert.Equal(t, avrox.NamespaceBasic, n)
	assert.Equal(t, avrox.BasicDecimalSchemaID, s)
	assert.Equal(t, avrox.CompNone, c)

	result, err4 := avrox.UnmarshalBasic(data)
	assert.NoError(t, err4)
	assert.Equal(t, value, result)
}

func TestMarshalBasicRawDate(t *testing.T) {
	value, _ := rawdate.New(2023, 2, 20)
	data, err := avrox.MarshalBasic(value, avrox.CompNone)
	assert.NoError(t, err)

	n, s, c, err3 := avrox.DecodeMagic(data[:avrox.MagicLen])
	assert.NoError(t, err3)
	assert.Equal(t, avrox.NamespaceBasic, n)
	assert.Equal(t, avrox.BasicRawDateSchemaID, s)
	assert.Equal(t, avrox.CompNone, c)

	result, err4 := avrox.UnmarshalBasic(data)
	assert.NoError(t, err4)
	assert.Equal(t, value, result)
	assert.True(t, result.(rawdate.RawDate).String() == "2023-02-20")
}

func TestMarshalBasicNulRawDate(t *testing.T) {
	value := rawdate.RawDate{}
	data, err := avrox.MarshalBasic(value, avrox.CompNone)
	assert.NoError(t, err)

	n, s, c, err3 := avrox.DecodeMagic(data[:avrox.MagicLen])
	assert.NoError(t, err3)
	assert.Equal(t, avrox.NamespaceBasic, n)
	assert.Equal(t, avrox.BasicRawDateSchemaID, s)
	assert.Equal(t, avrox.CompNone, c)

	result, err4 := avrox.UnmarshalBasic(data)
	assert.NoError(t, err4)
	assert.True(t, result.(rawdate.RawDate).IsZero())
	assert.True(t, result.(rawdate.RawDate).String() == "0001-01-01")
}

func TestMarshalMapStringAnySnappy(t *testing.T) {
	value := map[string]any{"foo": 1, "bar": "baz", "error": false}
	data, err := avrox.MarshalBasic(value, avrox.CompSnappy)
	assert.NoError(t, err)

	n, s, c, err3 := avrox.DecodeMagic(data[:avrox.MagicLen])
	assert.NoError(t, err3)
	assert.Equal(t, avrox.NamespaceBasic, n)
	assert.Equal(t, avrox.BasicMapStringAnySchemaID, s)
	assert.Equal(t, avrox.CompSnappy, c)

	text, err4 := avrox.UnmarshalBasic(data)
	assert.NoError(t, err4)
	assert.Equal(t, value, text)
}

func TestMarshalMapStringAnySnappyPtr(t *testing.T) {
	value := map[string]any{"foo": 1, "bar": "baz", "error": false}
	data, err := avrox.MarshalBasic(&value, avrox.CompSnappy)
	assert.NoError(t, err)

	n, s, c, err3 := avrox.DecodeMagic(data[:avrox.MagicLen])
	assert.NoError(t, err3)
	assert.Equal(t, avrox.NamespaceBasic, n)
	assert.Equal(t, avrox.BasicMapStringAnySchemaID, s)
	assert.Equal(t, avrox.CompSnappy, c)

	text, err4 := avrox.UnmarshalBasic(data)
	assert.NoError(t, err4)
	assert.Equal(t, value, text)
}

func TestMarshalAny(t *testing.T) {
	// we can use BasicString struct here (just use another magic)
	var int69 int8 = 69
	maxFloat64 := math.MaxFloat64
	maxFloat32 := float32(math.MaxFloat32)
	mt := &testdata.TestStruct{
		FieldInt: 42,
		FieldSubStruct: testdata.TestSubStruct{
			FieldSubString: "Bar",
		},
		FieldPtrSubStruct: &testdata.TestPtrSubStruct{
			FieldSubString:  "Baz",
			FieldPtrInt8:    &int69,
			FieldPtrFloat64: &maxFloat64,
			FieldPtrFloat32: &maxFloat32,
		},
		FieldSlicePtrSubStruct: []*testdata.TestPtrSubStruct{{
			FieldSubString: "Bla",
			FieldPtrInt8:   &int69,
		}, nil, {
			FieldSubString: "Blub",
			FieldPtrInt8:   nil,
		}},
		FieldString:  "Foo",
		FieldTime:    avrox.AvroTime(time.Now()),
		FieldDate:    avrox.AvroDate(time.Now()),
		FieldRawDate: rawdate.MustNew(2024, 02, 20),
	}
	data, err := avrox.MarshalAny(mt, nil, avrox.NamespacePrivate, avrox.SchemaUndefined, avrox.CompNone)
	assert.True(t, errors.Is(err, avrox.ErrSchemaNil))

	schema, errParse := avro.Parse(testdata.TestStructAVSC)
	assert.NoError(t, errParse)
	if errParse != nil {
		assert.FailNow(t, "future tests skipped")
	}
	assert.NotNil(t, schema)

	data, err = avrox.MarshalAny(mt, schema, avrox.NamespacePrivate, avrox.SchemaUndefined, avrox.CompNone)
	assert.NoError(t, err)
	if err != nil {
		assert.FailNow(t, "future tests skipped")
	}

	n, s, c, err3 := avrox.DecodeMagic(data[:avrox.MagicLen])
	assert.NoError(t, err3)
	assert.Equal(t, avrox.NamespacePrivate, n)
	assert.Equal(t, avrox.SchemaUndefined, s)
	assert.Equal(t, avrox.CompNone, c)

	basicString := &avrox.BasicString{}
	err4 := avrox.Unmarshal(data, basicString, nil)
	assert.True(t, errors.Is(err4, avrox.ErrWrongNamespace))

	testStruct := &testdata.TestStruct{}
	nID, sID, err5 := avrox.UnmarshalAny(data, schema, testStruct)
	assert.Equal(t, avrox.NamespacePrivate, nID)
	assert.Equal(t, avrox.SchemaUndefined, sID)
	assert.NoError(t, err5)
	assert.Equal(t, mt, testStruct)
}

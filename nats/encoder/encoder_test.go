package encoder_test

import (
	"github.com/metatexx/avrox"
	"github.com/metatexx/avrox/nats/encoder"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncoder(t *testing.T) {
	assert.Equal(t, 1, 1)
	ae := &encoder.AvroXEncoder{}
	value := &avrox.BasicMapStringAny{
		Magic: avrox.MustEncodeBasicMagic(avrox.BasicStringID, avrox.CompNone),
		Value: map[string]any{"Foo": 1, "bar": "baz"},
	}
	encoded, err := ae.Encode("", value)
	assert.NoError(t, err)
	decoded := &avrox.BasicMapStringAny{}
	err = ae.Decode("", encoded, decoded)
	assert.NoError(t, err)
	assert.Equal(t, value, decoded)
	assert.Equal(t, avrox.NamespaceBasic, decoded.NamespaceID())
	assert.Equal(t, avrox.BasicMapStringAnyID, decoded.SchemaID())
}

func encodeDecode(n int) int {
	ae := &encoder.AvroXEncoder{}
	value := &avrox.BasicMapStringAny{
		Magic: avrox.MustEncodeBasicMagic(avrox.BasicStringID, avrox.CompNone),
		Value: map[string]any{"Foo": n, "bar": "baz"},
	}
	encoded, err := ae.Encode("", value)
	if err != nil {
		panic(err)
	}
	decoded := &avrox.BasicMapStringAny{}
	err = ae.Decode("", encoded, decoded)
	return decoded.Value["Foo"].(int)
}

var resultSum int

func BenchmarkEncoder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		sum := 0
		for x := 0; x < 10; x++ {
			sum += encodeDecode(x)
		}
		resultSum = sum
	}
	//fmt.Printf("%#v\n", avro.DefaultSchemaCache)
}

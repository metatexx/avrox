package avrox

import (
	"errors"
	"github.com/hamba/avro/v2"
	"math/big"
	"time"
)

func MustEncodeBasicMagic(schemaID SchemaID, compression CompressionID) Magic {
	m, err := EncodeMagic(NamespaceBasic, schemaID, compression)
	if err != nil {
		panic(err)
	}
	return m
}

func MarshalBasic(src any, cID CompressionID) ([]byte, error) {
	var data []byte
	var errMarshall error
	switch v := src.(type) {
	case string:
		kind := &BasicString{
			Magic: MustEncodeBasicMagic(BasicStringSchemaID, cID),
			Value: v,
		}
		data, errMarshall = avro.Marshal(avro.MustParse(BasicStringAVSC), kind)
	case *string:
		kind := &BasicString{
			Magic: MustEncodeBasicMagic(BasicStringSchemaID, cID),
			Value: *v,
		}
		data, errMarshall = avro.Marshal(avro.MustParse(BasicStringAVSC), kind)
	case int:
		kind := &BasicInt{
			Magic: MustEncodeBasicMagic(BasicIntSchemaID, cID),
			Value: v,
		}
		data, errMarshall = avro.Marshal(avro.MustParse(BasicIntAVSC), kind)
	case *int:
		kind := &BasicInt{
			Magic: MustEncodeBasicMagic(BasicIntSchemaID, cID),
			Value: *v,
		}
		data, errMarshall = avro.Marshal(avro.MustParse(BasicIntAVSC), kind)
	case []byte:
		kind := &BasicByteSlice{
			Magic: MustEncodeBasicMagic(BasicByteSliceSchemaID, cID),
			Value: v,
		}
		data, errMarshall = avro.Marshal(avro.MustParse(BasicByteSliceAVSC), kind)
	case *[]byte:
		kind := &BasicByteSlice{
			Magic: MustEncodeBasicMagic(BasicByteSliceSchemaID, cID),
			Value: *v,
		}
		data, errMarshall = avro.Marshal(avro.MustParse(BasicByteSliceAVSC), kind)
	case map[string]any:
		kind := &BasicMapStringAny{
			Magic: MustEncodeBasicMagic(BasicMapStringAnySchemaID, cID),
			Value: v,
		}
		data, errMarshall = avro.Marshal(avro.MustParse(BasicMapStringAnyAVSC), kind)
	case *map[string]any:
		kind := &BasicMapStringAny{
			Magic: MustEncodeBasicMagic(BasicMapStringAnySchemaID, cID),
			Value: *v,
		}
		data, errMarshall = avro.Marshal(avro.MustParse(BasicMapStringAnyAVSC), kind)
	case time.Time:
		kind := &BasicTime{
			Magic: MustEncodeBasicMagic(BasicTimeSchemaID, cID),
			Value: v,
		}
		data, errMarshall = avro.Marshal(avro.MustParse(BasicTimeAVSC), kind)
	case *time.Time:
		kind := &BasicTime{
			Magic: MustEncodeBasicMagic(BasicTimeSchemaID, cID),
			Value: *v,
		}
		data, errMarshall = avro.Marshal(avro.MustParse(BasicTimeAVSC), kind)
	case *big.Rat:
		kind := &BasicDecimal{
			Magic: MustEncodeBasicMagic(BasicDecimalSchemaID, cID),
			Value: v,
		}
		data, errMarshall = avro.Marshal(avro.MustParse(BasicDecimalAVSC), kind)
	default:
		return nil, errors.New("unsupported type")
	}
	if errMarshall != nil {
		return nil, errMarshall
	}
	return CompressData(data, cID)
}

func UnmarshalBasic(src []byte) (any, error) {
	if len(src) == 0 {
		return nil, nil
	}
	if len(src) < MagicLen {
		return nil, ErrNotAvroX
	}
	nID, sID, _, err := DecodeMagic(src[:MagicLen])
	if err != nil {
		return nil, err
	}
	if nID != NamespaceBasic {
		return nil, ErrNoBasicNamespace
	}

	var errUnmarshalAny error

	//nolint:exhaustive // only runs for NamespaceBasic
	switch sID {
	case BasicStringSchemaID:
		kind := &BasicString{}
		nID, sID, errUnmarshalAny = UnmarshalAny(src, avro.MustParse(BasicStringAVSC), kind)
		if errUnmarshalAny != nil {
			return nil, errUnmarshalAny
		}
		if nID != NamespaceBasic && sID != BasicStringSchemaID {
			return nil, ErrNoBasicString
		}
		return kind.Value, nil
	case BasicIntSchemaID:
		kind := &BasicInt{}
		nID, sID, errUnmarshalAny = UnmarshalAny(src, avro.MustParse(BasicIntAVSC), kind)
		if errUnmarshalAny != nil {
			return nil, errUnmarshalAny
		}
		if nID != NamespaceBasic && sID != BasicIntSchemaID {
			return nil, ErrNoBasicInt
		}
		return kind.Value, nil
	case BasicByteSliceSchemaID:
		kind := &BasicByteSlice{}
		nID, sID, errUnmarshalAny = UnmarshalAny(src, avro.MustParse(BasicByteSliceAVSC), kind)
		if errUnmarshalAny != nil {
			return nil, errUnmarshalAny
		}
		if nID != NamespaceBasic && sID != BasicByteSliceSchemaID {
			return nil, ErrNoBasicByteSlice
		}
		return kind.Value, nil
	case BasicMapStringAnySchemaID:
		kind := &BasicMapStringAny{}
		nID, sID, errUnmarshalAny = UnmarshalAny(src, avro.MustParse(BasicMapStringAnyAVSC), kind)
		if errUnmarshalAny != nil {
			return nil, errUnmarshalAny
		}
		if nID != NamespaceBasic && sID != BasicByteSliceSchemaID {
			return nil, ErrNoBasicMapStringAny
		}
		return kind.Value, nil
	case BasicTimeSchemaID:
		kind := &BasicTime{}
		nID, sID, errUnmarshalAny = UnmarshalAny(src, avro.MustParse(BasicTimeAVSC), kind)
		if errUnmarshalAny != nil {
			return nil, errUnmarshalAny
		}
		if nID != NamespaceBasic && sID != BasicTimeSchemaID {
			return nil, ErrNoBasicTime
		}
		return kind.Value, nil
	case BasicDecimalSchemaID:
		kind := &BasicDecimal{}
		nID, sID, errUnmarshalAny = UnmarshalAny(src, avro.MustParse(BasicDecimalAVSC), kind)
		if errUnmarshalAny != nil {
			return nil, errUnmarshalAny
		}
		if nID != NamespaceBasic && sID != BasicTimeSchemaID {
			return nil, ErrNoBasicTime
		}
		return kind.Value, nil
	default:
		return nil, ErrNoBasicSchema
	}
}

func UnmarshalString(data []byte) (string, error) {
	x, err := UnmarshalBasic(data)
	if err != nil {
		return "", err
	}
	if str, ok := x.(string); ok {
		return str, nil
	}
	return "", ErrNoBasicString
}

func UnmarshalInt(data []byte) (int, error) {
	x, err := UnmarshalBasic(data)
	if err != nil {
		return 0, err
	}
	if n, ok := x.(int); ok {
		return n, nil
	}
	return 0, ErrNoBasicString
}

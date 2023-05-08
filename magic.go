package avrox

import "errors"

// 0x93 0bCCCNNNNN 0bSSSSSSSS 0bSSSSSPPP

var Marker byte = 0x93 // one of the best when analysing our data

type Magic [4]byte

func calculateParity(data []byte) byte {
	if len(data) != 4 {
		panic(errors.New("wrong len for calculate parity"))
	}
	// this was the best checksum when testing different implementations
	return (data[1] ^ data[2] + (data[3] & 0xF8)) & 0x07
}

func MustEncodePrivateMagic(compression CompressionID) Magic {
	m, err := EncodeMagic(0, 0, compression)
	if err != nil {
		panic(err)
	}
	return m
}

func EncodePrivateMagic(compression CompressionID) (Magic, error) {
	return EncodeMagic(0, 0, compression)
}

func EncodeMagic(namespace NamespaceID, schema SchemaID, compression CompressionID) (Magic, error) {
	if namespace < 0 || namespace > NamespaceMax {
		return Magic{}, ErrNamespaceIDOutOfRange
	}
	if compression < 0 || compression > CompMax {
		return Magic{}, ErrCompressionIDOutOfRange
	}
	if schema < 0 || schema > SchemaMax {
		return Magic{}, ErrSchemaIDOutOfRange
	}

	data := Magic{}
	data[0] = Marker
	data[1] = byte(int(compression<<5) | int(namespace))
	data[2] = byte(schema >> 5)
	data[3] = byte(schema << 3)

	calculatedParity := calculateParity(data[:])

	data[3] |= calculatedParity

	return data, nil
}

func DecodeMagic(data []byte) (NamespaceID, SchemaID, CompressionID, error) {
	if len(data) != 4 {
		return 0, 0, 0, ErrLengthInvalid
	}

	if data[0] != Marker {
		return 0, 0, 0, ErrMarkerInvalid
	}

	compression := CompressionID(int(data[1] >> 5))
	namespace := NamespaceID(int(data[1] & 0x1F))
	schema := SchemaID(int(data[2])<<5 | int(data[3]>>3))

	parityBits := data[3] & 0x07
	calculatedParity := calculateParity(data)

	if parityBits != calculatedParity {
		return 0, 0, 0, ErrParityCheckFailed
	}

	return namespace, schema, compression, nil
}

func IsMagic(data []byte) bool {
	if len(data) != 4 {
		return false
	}

	if data[0] != Marker {
		return false
	}

	parityBits := data[3] & 0x07
	calculatedParity := calculateParity(data)

	return parityBits == calculatedParity
}

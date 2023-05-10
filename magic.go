package avrox

import (
	"errors"
)

// //0x93 0bCCCNNNNN 0bSSSSSSSS 0bSSSSSPPP
// 0x93 0bCCCCCCCC 0bNNNNNNNN 0bNNSSSSSS 0bSSSSSSSS 0bVVVVVVVV
// 0x93 0bCCCCCCCC 0bNNNNNNNN 0bNNNNNNNN 0bSSSSSSSS 0bSSSSSSSS 0bSSSSSSSS 0bPPPPPPPP

var Marker byte = 0x93 // one of the best when analysing our data

const MagicLen = 8 // The struct will be aligned to 8 bytes anyway

type Magic [MagicLen]byte

const (
	offsetBasis uint = 2166136261
	prime       uint = 16777619
)

func calculateParity(data []byte) byte {
	if len(data) != MagicLen {
		panic(errors.New("wrong len for calculate parity"))
	}
	/*
		hash := offsetBasis
		hash ^= uint(data[1])
		hash *= prime
		hash ^= uint(data[2])
		hash *= prime
		hash ^= uint(data[3])
		hash *= prime
		hash ^= uint(data[4])
		hash *= prime
		hash ^= uint(data[5])
		hash *= prime
		hash ^= uint(data[6])
		hash *= prime
	*/
	return data[1] ^ data[2] - data[3] ^ data[4] + data[5] ^ data[6]
	//return data[1] - data[2] + data[3] - data[4] + data[5] - data[6]
	//return data[1] ^ data[2] ^ data[3] ^ data[4] ^ data[5] ^ data[6]
	//return byte((int(data[1]) + int(data[2]) + int(data[3]) + int(data[4]) + int(data[5]) + int(data[6])) % 255)

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

func EncodeMagic(namespace NamespaceID, schema SchemVerID, compression CompressionID) (Magic, error) {
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
	data[1] = byte(compression)
	data[2] = byte(namespace >> 8)
	data[3] = byte(namespace)
	data[4] = byte(schema >> 16)
	data[5] = byte(schema >> 8)
	data[6] = byte(schema) // used as version!
	//fmt.Printf("%x\n", data)
	/*
		data[1] = byte(compression & 0xff) // 32 (+3 bit Parity)
		data[2] = byte(namespace >> 2)     // 1024
		data[5] = byte(schema & 0xff)      // Should be used as version of the schema
		data[6] = byte(schema & 0xff)      // Should be used as version of the schema
		data[3] = byte(namespace<<6) | byte(schema>>16)
		data[4] = byte((schema >> 8) & 0xff)
		data[7] = byte(schema & 0xff) // Should be used as version of the schema
	*/
	calculatedParity := calculateParity(data[:])

	data[7] |= calculatedParity

	return data, nil
}

func DecodeMagic(data []byte) (NamespaceID, SchemVerID, CompressionID, error) {
	if len(data) != MagicLen {
		return 0, 0, 0, ErrLengthInvalid
	}

	if data[0] != Marker {
		return 0, 0, 0, ErrMarkerInvalid
	}

	compression := CompressionID(int(data[1]))
	namespace := NamespaceID((int(data[2]) << 8) | int(data[3]))
	schema := SchemVerID((int(data[4]) << 16) | (int(data[5]) << 8) | int(data[6]))
	parityBits := data[7]
	calculatedParity := calculateParity(data)

	if parityBits != calculatedParity {
		return 0, 0, 0, ErrParityCheckFailed
	}

	return namespace, schema, compression, nil
}

func IsMagic(data []byte) bool {
	if len(data) != MagicLen {
		return false
	}

	if data[0] != Marker {
		return false
	}

	parityBits := data[7]
	calculatedParity := calculateParity(data)

	return parityBits == calculatedParity
}

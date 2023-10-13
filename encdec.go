package avrox

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/golang/snappy"
	"github.com/metatexx/mxx/wfl"
)

func CompressData(data []byte, cID CompressionID) ([]byte, error) {
	switch cID {
	case CompNone:
		return data, nil
	case CompSnappy:
		edata := snappy.Encode(nil, data[MagicLen:])
		return append(data[0:MagicLen], edata...), nil
	case CompFlate:
		var eData bytes.Buffer
		w, err := flate.NewWriter(&eData, flate.DefaultCompression)
		if err != nil {
			return nil, err
		}
		_, err = w.Write(data[MagicLen:])
		if err != nil {
			return nil, err
		}
		err = w.Close()
		if err != nil {
			return nil, err
		}
		return append(data[0:MagicLen], eData.Bytes()...), nil
	case CompGZip:
		var eData bytes.Buffer
		w := gzip.NewWriter(&eData)
		_, err := w.Write(data[MagicLen:])
		if err != nil {
			return nil, err
		}
		err = w.Close()
		if err != nil {
			return nil, err
		}
		return append(data[0:MagicLen], eData.Bytes()...), nil
	default:
		return nil, wfl.ErrorWithSkip(ErrCompressionUnsupported, 3)
	}
}

func DecompressData(data []byte, cID CompressionID) ([]byte, error) {
	//nolint:exhaustive // can't be exhaustive
	var uncompressed []byte
	switch cID {
	case CompNone:
		return data, nil
	case CompSnappy:
		dData, errDecode := snappy.Decode(nil, data[MagicLen:])
		if errDecode != nil {
			return nil, errors.Join(ErrDecompress, errDecode)
		}
		uncompressed = dData
	case CompFlate:
		b := flate.NewReader(bytes.NewReader(data[MagicLen:]))
		var dData bytes.Buffer
		_, err := dData.ReadFrom(b)
		if err != nil {
			return nil, fmt.Errorf("read flate error: %w", err)
		}
		err = b.Close()
		if err != nil {
			return nil, fmt.Errorf("close flate error: %w", err)
		}
		uncompressed = dData.Bytes()
	case CompGZip:
		b, err := gzip.NewReader(bytes.NewReader(data[MagicLen:]))
		if err != nil {
			return nil, fmt.Errorf("new reader gzip error: %w", err)
		}
		var dData bytes.Buffer
		_, err = dData.ReadFrom(b)
		if err != nil {
			return nil, fmt.Errorf("read gzip error: %w", err)
		}
		err = b.Close()
		if err != nil {
			return nil, fmt.Errorf("close gzip error: %w", err)
		}
		uncompressed = dData.Bytes()
	default:
		return nil, ErrCompressionUnsupported
	}
	//return append(data[:MagicLen], uncompressed...), nil

	// rewrite the magic header, so it does not say it is compressed anymore
	ns, sc, _, errDecode := DecodeMagic(data[0:MagicLen])
	if errDecode != nil {
		return nil, errDecode
	}
	magic, errMagic := EncodeMagic(ns, sc, CompNone)
	if errMagic != nil {
		return nil, errMagic
	}
	return append(magic[:], uncompressed...), nil
}

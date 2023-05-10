package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/metatexx/avrox"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

// Mostly generated together with GPT-4 and just an experiment

const maxLenForFilenameDisplay = 70

var found = 0

var counts = make(map[byte]int, 256)

func Scanner(path string) {
	scanned := 0
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			//nolint:nilerr // skip all errors
			return nil
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !info.IsDir() && info.Mode().IsRegular() && info.Size() >= avrox.MagicLen {
			scanned++
			p := filepath.Dir(path)
			if len(p) > maxLenForFilenameDisplay {
				p = p[len(p)-maxLenForFilenameDisplay:]
			}
			p = filepath.Join(p, filepath.Base(path))
			if len(p) > maxLenForFilenameDisplay {
				p = p[:maxLenForFilenameDisplay-10] + "..." + p[len(p)-7:]
			}

			fmt.Printf("\u001B[KScanned %d / found %d: %s\r",
				scanned, found, p)
			err = findMagicBytes(ctx, path, true, true)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", path, err)
			}
		}
		return nil
	})

	fmt.Printf("\u001B[KScanned %d / found %d\n",
		scanned, found)

	if err != nil {
		fmt.Printf("Error walking the path: %v\n", err)
	}
}

func findMagicBytes(ctx context.Context, path string, onlyAtStart bool, verbose bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, file.Close())
	}()

	offset := 0
	reader := bufio.NewReader(file)
	buffer := make([]byte, avrox.MagicLen)

	var n int
	n, err = io.ReadFull(reader, buffer)
	if err != nil {
		return err
	}
	offset += n

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if buffer[0] == avrox.Marker && avrox.IsMagic(buffer) {
			nID, sID, cID, _ := avrox.DecodeMagic(buffer)
			reader.Size()
			found++
			if verbose {
				fmt.Printf("Found magic bytes at offset %d in file: %s (N: %d / S: %d / C: %d)\n",
					offset-avrox.MagicLen, path, nID, sID, cID)
			}
		}

		if onlyAtStart {
			return nil
		}

		var byteRead byte
		byteRead, err = reader.ReadByte()
		counts[byteRead]++
		offset++
		if errors.Is(err, io.EOF) {
			err = nil
			break
		}
		if err != nil {
			return err
		}

		buffer = append(buffer[1:], byteRead)
	}
	return err
}

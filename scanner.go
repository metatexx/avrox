package avrox

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"syscall"
)

// Mostly generated together with GPT-4 and just an experiment

const maxLen = 70

var found = 0

var counts = make(map[byte]int, 256)

func Scanner(path string) {
	scanned := 0
	for x := 0; x < 256; x++ {
		counts[byte(x)] = 0
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			//nolint:nilerr // skip all errors
			return nil
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !info.IsDir() && info.Mode().IsRegular() && info.Size() >= 4 {
			scanned++
			p := filepath.Dir(path)
			if len(p) > maxLen {
				p = p[len(p)-maxLen:]
			}
			p = filepath.Join(p, filepath.Base(path))
			if len(p) > maxLen {
				p = p[:maxLen-10] + "..." + p[len(p)-7:]
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

	byteIntPairs := make([]struct {
		key   byte
		value int
	}, 0, len(counts))

	for k, v := range counts {
		byteIntPairs = append(byteIntPairs, struct {
			key   byte
			value int
		}{k, v})
	}

	sort.Slice(byteIntPairs, func(i, j int) bool {
		return byteIntPairs[i].value < byteIntPairs[j].value
	})

	topBytes := byteIntPairs[:10]

	fmt.Println("Top 10 markers with the least false positives:")
	for i, pair := range topBytes {
		fmt.Printf("%d. 0x%02x Char: %c: count: %d\n", i+1, pair.key, pair.key, pair.value)
	}
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
	buffer := make([]byte, 4)

	var n int
	n, err = io.ReadFull(reader, buffer)
	if err != nil {
		return err
	}
	offset += n

	// checking parities for all of the bytes
	for x := 0; x < 256; x++ {
		if buffer[0] == byte(x) {
			parityBits := buffer[3] & 0x07
			calculatedParity := calculateParity(buffer)
			if parityBits == calculatedParity {
				counts[byte(x)]++
			}
		}
	}
	/*
		counts[buffer[0]]++
		if !onlyAtStart {
			counts[buffer[1]]++
			counts[buffer[2]]++
			counts[buffer[3]]++
		}*/

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if buffer[0] == Marker && IsMagic(buffer) {
			nID, sID, cID, _ := DecodeMagic(buffer)
			reader.Size()
			found++
			if verbose {
				fmt.Printf("Found magic bytes at offset %d in file: %s (N: %d / S: %d / C: %d)\n", offset-4, path, nID, sID, cID)
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

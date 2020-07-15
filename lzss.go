package main

import (
	"fmt"
	"io/ioutil"

	"github.com/alexflint/go-arg"
)

func createBitMask(inputs []bool) byte {
	if len(inputs) != 8 {
		panic("Bit mask input must be 8 logical values.")
	}

	var mask byte = 0
	for _, input := range inputs {
		var bit byte = 0
		if input {
			bit = 1
		}

		mask = mask << 1
		mask += bit
	}

	return mask
}

func readBitMask(input byte) []bool {
	readMask := byte(0b10000000)
	output := make([]bool, 0, 8)

	for i := 0; i < 8; i++ {
		tmp := input&readMask >= 1
		output = append(output, tmp)
		readMask = readMask >> 1
	}

	return output
}

func createPointer(offset int, length int) []byte {
	output := make([]byte, 2)
	output[0] = byte(offset >> 4)
	output[1] = byte(offset<<4) + byte(length)

	return output
}

func readPointer(input []byte) (int, int) {
	length := int(input[1] & 0b00001111)
	offset := (int(input[0]) << 4) | (int(input[1]) >> 4)

	return offset, length
}

func findMatch(input []byte, startPos int) (int, int) {
	lookback := 4095
	if startPos < lookback {
		lookback = startPos
	}

	maxOffset := 0
	maxLength := 0

	for i := startPos - lookback; i < startPos; i++ {
		if input[i] != input[startPos] {
			continue
		}

		offset := startPos - i
		length := 1
		for startPos+length < len(input) && input[i+length] == input[startPos+length] && length < 15 {
			length++
		}

		if length > maxLength {
			maxLength = length
			maxOffset = offset
		}

	}

	return maxOffset, maxLength
}

func encode(input []byte) []byte {
	inLen := len(input)
	inPos := 0
	var masks []bool
	var encoded []byte

	// Write the compressed data to output and leave a byte free every 9 bytes
	// for the bit masks
	for inPos < inLen {
		offset, length := findMatch(input, inPos)
		if length <= 2 {
			encoded = append(encoded, input[inPos])
			masks = append(masks, false)
			inPos++
			continue
		}

		encoded = append(encoded, createPointer(offset, length)...)
		masks = append(masks, true, true)
		inPos += length
	}

	// Add masks so the length of masks is a multiple of 8
	masks = append(masks, make([]bool, 8)...)

	// Add the bit masks
	// At the moment this needs to copy the hole array
	// Todo: dont copy the encoded byte by byte
	output := make([]byte, 0, len(encoded)+len(encoded)/9)
	for i := 0; i < len(encoded); i++ {
		if i%8 == 0 {
			mask := masks[i : i+8] // Maybe 7
			output = append(output, createBitMask(mask))
		}

		output = append(output, encoded[i])
	}

	return output
}

func decode(input []byte) []byte {

	// Read all mask bytes
	masks := make([]bool, 0, len(input)/9*8)
	encoded := make([]byte, 0, len(input)-len(input)/9)
	for i := 0; i < len(input); i++ {
		if i%9 == 0 {
			masks = append(masks, readBitMask(input[i])...)
			continue
		}
		encoded = append(encoded, input[i])
	}

	// Decompress
	var output []byte
	for i := 0; i < len(encoded); i++ {
		if !masks[i] {
			output = append(output, encoded[i])
			continue
		}

		offset, length := readPointer(encoded[i : i+2])
		index := len(output)
		for j := 0; j < length; j++ {
			output = append(output, output[index-offset+j])
		}
		i++
	}

	return output
}

func main() {
	// Setup the command line arguments
	var args struct {
		Input      string `arg:"-i,required"`
		Output     string `arg:"-o,required"`
		Compress   bool   `arg:"-c"`
		Decompress bool   `arg:"-d"`
		Verbose    bool   `arg:"-v"`
	}

	p := arg.MustParse(&args)
	if args.Compress == false && args.Decompress == false {
		p.Fail("you must provide either --compress or --decompress")
	}

	// Load the file and compress or decompress
	data, err := ioutil.ReadFile(args.Input)
	if err != nil {
		fmt.Printf("Cannot read file %s: %s", args.Input, err.Error())
		return
	}

	var output []byte
	if args.Compress {
		output = encode(data)
	}

	if args.Decompress {
		output = decode(data)
	}

	err = ioutil.WriteFile(args.Output, output, 0644)
	if err != nil {
		fmt.Printf("Cannot write file %s: %s", args.Output, err.Error())
		return
	}

	if args.Verbose && args.Compress {
		fmt.Printf("Original size: %d bytes\n", len(data))
		fmt.Printf("Compressed size: %d bytes\n", len(output))
		fmt.Printf("Reduced to %0.1f%% of the original\n", float64(len(output))/float64(len(data))*100)
	}
	if args.Verbose && args.Decompress {
		fmt.Printf("Original size: %d\n", len(output))
		fmt.Printf("Compressed size: %d\n", len(data))
	}
}

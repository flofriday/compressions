package main

import "fmt"

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

func findMatch(input []byte, startPos int, lookback int) (int, int) {
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
		for startPos+length < len(input) && input[i+length] == input[startPos+length] && length < 4094 {
			length++
		}

		if length > maxLength {
			maxLength = length
			maxOffset = offset
		}

	}

	fmt.Printf("%v - %v\n", maxOffset, maxLength)
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
		offset, length := findMatch(input, inPos, 2<<12)
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

func main() {
	s := []byte("Hose Dose Rose")
	for i, b := range s {
		fmt.Printf("%d: %08b\n", i, b)
	}
	fmt.Println()

	tmp := encode(s)

	for i, b := range tmp {
		fmt.Printf("%d: %08b\n", i, b)
	}
}

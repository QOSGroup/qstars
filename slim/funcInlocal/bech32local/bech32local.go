package bech32local

import (
	"fmt"
	"github.com/pkg/errors"
)

//ConvertAndEncode converts from a base64 encoded byte string to base32 encoded byte string and then to bech32
func ConvertAndEncode(hrp string, data []byte) (string, error) {
	converted, err := ConvertBits(data, 8, 5, true)
	if err != nil {
		return "", errors.Wrap(err, "encoding bech32 failed")
	}
	return Encode(hrp, converted)

}

// ConvertBits converts a byte slice where each byte is encoding fromBits bits,
// to a byte slice where each byte is encoding toBits bits.
func ConvertBits(data []byte, fromBits, toBits uint8, pad bool) ([]byte, error) {
	if fromBits < 1 || fromBits > 8 || toBits < 1 || toBits > 8 {
		return nil, fmt.Errorf("only bit groups between 1 and 8 allowed")
	}

	// The final bytes, each byte encoding toBits bits.
	var regrouped []byte

	// Keep track of the next byte we create and how many bits we have
	// added to it out of the toBits goal.
	nextByte := byte(0)
	filledBits := uint8(0)

	for _, b := range data {

		// Discard unused bits.
		b = b << (8 - fromBits)

		// How many bits remaining to extract from the input data.
		remFromBits := fromBits
		for remFromBits > 0 {
			// How many bits remaining to be added to the next byte.
			remToBits := toBits - filledBits

			// The number of bytes to next extract is the minimum of
			// remFromBits and remToBits.
			toExtract := remFromBits
			if remToBits < toExtract {
				toExtract = remToBits
			}

			// Add the next bits to nextByte, shifting the already
			// added bits to the left.
			nextByte = (nextByte << toExtract) | (b >> (8 - toExtract))

			// Discard the bits we just extracted and get ready for
			// next iteration.
			b = b << toExtract
			remFromBits -= toExtract
			filledBits += toExtract

			// If the nextByte is completely filled, we add it to
			// our regrouped bytes and start on the next byte.
			if filledBits == toBits {
				regrouped = append(regrouped, nextByte)
				filledBits = 0
				nextByte = 0
			}
		}
	}

	// We pad any unfinished group if specified.
	if pad && filledBits > 0 {
		nextByte = nextByte << (toBits - filledBits)
		regrouped = append(regrouped, nextByte)
		filledBits = 0
		nextByte = 0
	}

	// Any incomplete group must be <= 4 bits, and all zeroes.
	if filledBits > 0 && (filledBits > 4 || nextByte != 0) {
		return nil, fmt.Errorf("invalid incomplete group")
	}

	return regrouped, nil
}

func Encode(hrp string, data []byte) (string, error) {
	// Calculate the checksum of the data and append it at the end.
	checksum := bech32Checksum(hrp, data)
	combined := append(data, checksum...)

	// The resulting bech32 string is the concatenation of the hrp, the
	// separator 1, data and checksum. Everything after the separator is
	// represented using the specified charset.
	dataChars, err := toChars(combined)
	if err != nil {
		return "", fmt.Errorf("unable to convert data bytes to chars: "+
			"%v", err)
	}
	return hrp + "1" + dataChars, nil
}

// For more details on the checksum calculation, please refer to BIP 173.
func bech32Checksum(hrp string, data []byte) []byte {
	// Convert the bytes to list of integers, as this is needed for the
	// checksum calculation.
	integers := make([]int, len(data))
	for i, b := range data {
		integers[i] = int(b)
	}
	values := append(bech32HrpExpand(hrp), integers...)
	values = append(values, []int{0, 0, 0, 0, 0, 0}...)
	polymod := bech32Polymod(values) ^ 1
	var res []byte
	for i := 0; i < 6; i++ {
		res = append(res, byte((polymod>>uint(5*(5-i)))&31))
	}
	return res
}

// For more details on HRP expansion, please refer to BIP 173.
func bech32HrpExpand(hrp string) []int {
	v := make([]int, 0, len(hrp)*2+1)
	for i := 0; i < len(hrp); i++ {
		v = append(v, int(hrp[i]>>5))
	}
	v = append(v, 0)
	for i := 0; i < len(hrp); i++ {
		v = append(v, int(hrp[i]&31))
	}
	return v
}

var gen = []int{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}

func bech32Polymod(values []int) int {
	chk := 1
	for _, v := range values {
		b := chk >> 25
		chk = (chk&0x1ffffff)<<5 ^ v
		for i := 0; i < 5; i++ {
			if (b>>uint(i))&1 == 1 {
				chk ^= gen[i]
			}
		}
	}
	return chk
}

const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

// toChars converts the byte slice 'data' to a string where each byte in 'data'
// encodes the index of a character in 'charset'.
func toChars(data []byte) (string, error) {
	result := make([]byte, 0, len(data))
	for _, b := range data {
		if int(b) >= len(charset) {
			return "", fmt.Errorf("invalid data byte: %v", b)
		}
		result = append(result, charset[b])
	}
	return string(result), nil
}

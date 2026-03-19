package hash

import (
	"math/big"
	"strconv"
)

// Encode converts a source string to a base 36 representation
func Encode(source uint64) string {
	converted := convert(source)
	reversedBit := reverseBit(converted)

	decimalNumber := new(big.Int)
	decimalNumber.SetString(reversedBit, 10)

	return decimalNumber.Text(36)
}

func Decode(target string) uint64 {
	decimalNumber := convertBase(target, 36, 10)
	reversedBit := reverseBit(decimalNumber)
	return invert(reversedBit)
}

func convert(number uint64) string {
	number = number*11 + 1
	return strconv.FormatUint(number, 10)
}

func invert(source string) uint64 {
	number, _ := strconv.ParseUint(source, 10, 64)
	number = (number - 1) / 11
	return number
}

func convertBase(number string, from, to int) string {
	sourceNumber, _ := strconv.ParseUint(number, from, 64)
	targerNumber := strconv.FormatUint(sourceNumber, to)
	return targerNumber
}

func reverseBit(number string) string {
	n, _ := strconv.ParseUint(number, 10, 64)
	reversedBit := n ^ ((1 << 62) - 1)
	return strconv.FormatUint(reversedBit, 10)
}

package utils

import (
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"magic.pathao.com/parcel/prism/internal/cnst"
)

func SanitizePhoneNumber(number string) string {
	// Remove hyphens
	number = strings.ReplaceAll(number, "-", "")
	// Remove white spaces, tabs, line ends
	re := regexp.MustCompile(`\s+`)
	number = re.ReplaceAllString(number, "")
	return number
}

func GenerateUniqueValue(id int, capitalize bool) string {
	const p = 10006890001
	var returnable string

	if id > 0 {
		returnable = strconv.FormatInt(int64(p+id), 36)
	} else {
		returnable = strconv.FormatInt(time.Now().Unix(), 36)
	}

	if capitalize {
		returnable = strings.ToUpper(returnable)
	}

	return returnable
}

func GetCurrentTime(countryId uint) time.Time {
	timeZone := cnst.BD_TIMEZONE
	if countryId == cnst.NEPAL {
		timeZone = cnst.NP_TIMEZONE
	}
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		fmt.Println("Error loading location: ", err)
		return time.Now()
	} else {
		return time.Now().In(location)
	}
}

// Encode converts a source string to a base 36 representation
func Encode(source string) string {
	converted := convert(source)
	reversedBit := reverseBit(converted)

	decimalNumber := new(big.Int)
	decimalNumber.SetString(reversedBit, 10)

	return decimalNumber.Text(36)
}

func Decode(target string) string {
	decimalNumber := new(big.Int)
	decimalNumber.SetString(target, 36)

	reversedBit := reverseBit(decimalNumber.String())
	return invert(reversedBit)
}

func convert(source string) string {
	number, _ := strconv.Atoi(source)
	number = number*11 + 1

	return strconv.Itoa(number)
}

func invert(source string) string {
	number, _ := strconv.Atoi(source)
	number = (number - 1) / 11

	return strconv.Itoa(number)
}

func reverseBit(number string) string {
	n, _ := strconv.ParseInt(number, 10, 64)
	reversedBit := n ^ (int64(math.Pow(2, 62)) - 1)

	return strconv.FormatInt(reversedBit, 10)
}

func ToSnakeCase(str string) string {
	parts := strings.Fields(str)
	return strings.ToLower(strings.Join(parts, "_"))
}

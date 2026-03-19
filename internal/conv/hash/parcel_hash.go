package hash

import (
	"strconv"
	"strings"
	"time"
)

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

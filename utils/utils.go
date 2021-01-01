package utils

import (
	"fmt"
	"strings"
)

func Prettify(buffer []byte) string {
	s := make([]string, 0, len(buffer))

	for i, b := range buffer {
		if i >= 16 {
			s = append(s, "<truncated>")
			break
		}

		s = append(s, fmt.Sprintf("%02x", b))
	}

	return fmt.Sprintf("[ %s ] (%d)", strings.Join(s, " "), len(buffer))
}

func Dump(buffer []byte) string {
	s := make([]string, 0, len(buffer))
	for _, b := range buffer {
		s = append(s, fmt.Sprintf("0x%02x", b))
	}

	return fmt.Sprintf("{%s}", strings.Join(s, ", "))
}

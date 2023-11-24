package utils

import (
	"crypto/rand"
	"io"
	"strings"
)

var otpList = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GenerateOtp(max int) string {
	b := make([]byte, max)
	n, _ := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		return strings.Repeat("9", max)
	}
	for i := 0; i < len(b); i++ {
		b[i] = otpList[int(b[i])%len(otpList)]
	}
	return string(b)
}

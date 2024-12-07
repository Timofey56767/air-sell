package tools

import (
	"crypto/sha256"
	"fmt"
)



func Sum256(s []byte) string{
	sum := sha256.Sum256(s)
	return fmt.Sprintf("%x", sum)
}
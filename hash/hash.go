package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func HashStringMD5(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

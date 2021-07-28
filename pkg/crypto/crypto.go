package crypto

import (
	"crypto/md5"
	"fmt"
)

func GetMD5(load string) string {
	res := md5.Sum([]byte(load))
	return fmt.Sprintf("%x", res)
}
